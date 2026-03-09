package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/handler"
	"github.com/mabishka/lupanova/internal/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_create(t *testing.T) {

	ctx, fnTimeoutCancel := context.WithTimeoutCause(context.Background(), time.Second*10, errors.New("stop timeout test"))
	defer fnTimeoutCancel()
	tests := []struct {
		name    string // description of this test case
		wantErr bool
	}{
		{
			name:    "positive",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := create(context.WithCancelCause(ctx))
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_runHttp(t *testing.T) {

	cfg := config.New()
	srv := handler.New(cfg.GetBaseAddress())

	ctx, fnCancel := context.WithTimeoutCause(context.Background(), time.Second*2, errors.New("test stop with timeout 2 seconds"))
	defer fnCancel()
	tests := []struct {
		name        string // description of this test case
		enableHTTPS bool
		wantErr     bool
	}{
		{
			name:        "positive_http",
			enableHTTPS: false,
			wantErr:     true,
		},
		{
			name:        "positive_https",
			enableHTTPS: true,
			wantErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			grpcServer := grpc.NewServer()
			proto.RegisterShortenerServiceServer(grpcServer, srv)

			tlsConfig := &tls.Config{}

			if cfg.IsEnableHTTPS() {
				if cert, err := makeCertificate(); err == nil {
					tlsConfig.Certificates = cert
				}
			}

			httpServer := &http.Server{
				Addr:         cfg.GetServerAddress(),
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
				IdleTimeout:  120 * time.Second,
				TLSConfig:    tlsConfig,
			}

			go configureStop(ctx, grpcServer, httpServer)

			runHttp(cfg.GetServerAddress(), httpServer, test.enableHTTPS)
		})
	}
}

func Test_runGrpc(t *testing.T) {

	cfg := config.New()
	srv := handler.New(cfg.GetBaseAddress())

	ctx, fnCancel := context.WithTimeoutCause(context.Background(), time.Second*2, errors.New("test stop with timeout 2 seconds"))
	defer fnCancel()
	tests := []struct {
		name        string // description of this test case
		enableHTTPS bool
		wantErr     bool
	}{
		{
			name:        "positive_grpc",
			enableHTTPS: false,
			wantErr:     true,
		},
		{
			name:        "negative_grpc",
			enableHTTPS: true,
			wantErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			grpcServer := grpc.NewServer()
			proto.RegisterShortenerServiceServer(grpcServer, srv)

			tlsConfig := &tls.Config{}

			if cfg.IsEnableHTTPS() {
				if cert, err := makeCertificate(); err == nil {
					tlsConfig.Certificates = cert
				}
			}

			httpServer := &http.Server{
				Addr:         cfg.GetServerAddress(),
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
				IdleTimeout:  120 * time.Second,
				TLSConfig:    tlsConfig,
			}

			go configureStop(ctx, grpcServer, httpServer)

			runHttp(cfg.GetServerAddress(), httpServer, test.enableHTTPS)
		})
	}
}
