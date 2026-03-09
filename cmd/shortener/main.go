package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/mabishka/lupanova/internal/auth"
	"github.com/mabishka/lupanova/internal/compress"
	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/handler"
	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/internal/proto"
	"github.com/mabishka/lupanova/internal/repository/audit"
	"github.com/mabishka/lupanova/internal/repository/connloader"
	"github.com/mabishka/lupanova/internal/repository/fileloader"
)

const stopTimeout = 5 * time.Second
const defaultGrpcAddress = ":5300"

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	fmt.Println("Build version: ", buildVersion)
	fmt.Println("Build date: ", buildDate)
	fmt.Println("Build commit: ", buildCommit)

	if err := create(context.WithCancelCause(context.Background())); err != nil {
		log.Fatalf("exist with error: %v", err)
	}
}

func create(ctx context.Context, fnCancel context.CancelCauseFunc) error {

	config := config.New()
	if err := logger.InitLogger(config.GetLogLevel()); err != nil {
		panic(err)
	}

	logger.Log().Info("config",
		zap.String("serverAddress", config.GetServerAddress()),
		zap.String("baseAddress", config.GetBaseAddress()),
		zap.String("logLevel", config.GetLogLevel()),
		zap.String("fileName", config.GetFileName()),
		zap.String("connAddress", config.GetConnAddress()),
	)

	server := handler.New(config.GetBaseAddress())

	var loader model.StorageLoader
	var conn model.ConnLoader
	if config.GetConnAddress() != "" {
		loader = connloader.New(config.GetConnAddress())
		conn, _ = loader.(model.ConnLoader)
		if err := server.Load(context.Background(), loader); err != nil {
			logger.Log().Info("conn not loaded", zap.Error(err))
			loader = nil
		} else {
			logger.Log().Info("conn storage usage")
		}

	}

	if loader == nil && config.GetFileName() != "" {
		loader = fileloader.New(config.GetFileName())

		if err := server.Load(context.Background(), loader); err != nil {
			logger.Log().Error("file not loaded", zap.Error(err))
			loader = nil
		} else {
			logger.Log().Info("file storage usage")
		}
	}

	if loader == nil {
		logger.Log().Info("memory storage usage")
	}

	connServer := handler.NewConn(conn)

	auditEvent := audit.NewAuditEvent()

	if config.GetAuditFile() != "" {
		auditEvent.Register(audit.NewFileObserver(config.GetAuditFile()))
	}

	if config.GetAuditAddress() != "" {
		auditEvent.Register(audit.NewAddressObserver(config.GetAuditAddress()))
	}

	server.SetAudit(auditEvent)
	server.SetTrustedSubnet(config.GetTrustedSubnet())

	router := chi.NewRouter()

	router.Use(logger.WithLogging)
	router.Use(compress.WithCompress)
	router.Use(auth.WithAuth)

	router.Mount("/debug", middleware.Profiler())

	router.Post("/", server.HandlerPostFull)
	router.Post("/api/shorten", server.HandlerPostFullJSON)
	router.Post("/api/shorten/batch", server.HandlerPostBatch)
	router.Get("/{id}", server.HandlerGetFull)
	router.Get("/ping", connServer.HandlerGetPing)
	router.Delete("/api/user/urls", server.HandlerDelete)
	router.Get("/api/user/urls", server.HandlerGetUser)
	router.Get("/api/internal/stats", server.HandlerGetStat)

	tlsConfig := &tls.Config{}

	if config.IsEnableHTTPS() {
		if cert, err := makeCertificate(); err == nil {
			tlsConfig.Certificates = cert
		}
	}

	grpcServer := grpc.NewServer()
	proto.RegisterShortenerServiceServer(grpcServer, server)

	httpServer := &http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    tlsConfig,
	}

	go configureStop(ctx, grpcServer, httpServer)

	var wg errgroup.Group

	wg.Go(func() error {
		return runHttp(defaultGrpcAddress, httpServer, config.IsEnableHTTPS())
	})

	wg.Go(func() error {
		return runGrpc(config.GetServerAddress(), grpcServer)
	})

	if err := wg.Wait(); err != nil {
		fnCancel(err)
		return err
	}
	fnCancel(nil)

	logger.Log().Info("exit")

	return nil
}

func configureStop(ctx context.Context, gsrv *grpc.Server, srv *http.Server) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-sigint:
		logger.Log().Info("stop with signal", zap.String("signal", s.String()))
	case <-ctx.Done():
		logger.Log().Info("stop with context", zap.Error(context.Cause(ctx)))
	}

	gsrv.GracefulStop()

	stopCtx, cancel := context.WithTimeoutCause(context.Background(), stopTimeout, fmt.Errorf("server Shutdown with timeout %v", stopTimeout))
	defer cancel()
	if err := srv.Shutdown(stopCtx); err != nil {
		logger.Log().Info("HTTP server shutdown", zap.Error(err))
	}

}

func runHttp(address string, srv *http.Server, isEnableHTTPS bool) error {

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	if isEnableHTTPS {
		if err := srv.ServeTLS(l, "", ""); err != http.ErrServerClosed {
			logger.Log().Info("HTTP server Serve", zap.Error(err))
			return err
		}
	} else {
		if err := srv.Serve(l); err != http.ErrServerClosed {
			logger.Log().Info("HTTP server Serve", zap.Error(err))
			return err
		}

	}
	return nil
}

func runGrpc(address string, gsrv *grpc.Server) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	if err := gsrv.Serve(l); err != http.ErrServerClosed {
		logger.Log().Info("GRPC server Serve", zap.Error(err))
		return err
	}
	return nil
}

/*
func run(_ context.Context, l net.Listener, gsrv *grpc.Server, srv *http.Server, isEnableHTTPS bool) error {

		var wg errgroup.Group

		wg.Go(func() error {
			if isEnableHTTPS {
				if err := srv.ServeTLS(l, "", ""); err != http.ErrServerClosed {
					logger.Log().Info("HTTP server Serve", zap.Error(err))
					return err
				}
			} else {
				if err := srv.Serve(l); err != http.ErrServerClosed {
					logger.Log().Info("HTTP server Serve", zap.Error(err))
					return err
				}

			}
			return nil
		})

		wg.Go(func() error {
			if err := gsrv.Serve(l); err != http.ErrServerClosed {
				logger.Log().Info("GRPC server Serve", zap.Error(err))
				return err
			}
			return nil
		})

		logger.Log().Info("exit")
		return wg.Wait()
	}
*/
func makeCertificate() ([]tls.Certificate, error) {
	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	var certPEM bytes.Buffer
	if err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return nil, err
	}

	var privateKeyPEM bytes.Buffer
	if err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return nil, err
	}

	certPair, err := tls.X509KeyPair(certPEM.Bytes(), privateKeyPEM.Bytes())
	if err != nil {
		return nil, err
	}

	return []tls.Certificate{certPair}, nil
}
