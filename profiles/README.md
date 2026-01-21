File: shortener.exe
Build ID: D:\Practicum\projects\iter17\cmd\shortener\shortener.exe2026-01-21 22:05:17.6074281 +0300 MSK
Type: inuse_space
Time: 2026-01-21 22:05:21 MSK
Showing nodes accounting for -6382.84kB, 53.24% of 11989.89kB total
      flat  flat%   sum%        cum   cum%
-4512.93kB 37.64% 37.64% -8137.14kB 67.87%  compress/flate.NewWriter (inline)
-1951.87kB 16.28% 53.92% -3624.21kB 30.23%  compress/flate.(*compressor).init
 1762.94kB 14.70% 39.22%  1762.94kB 14.70%  runtime/pprof.StartCPUProfile
-1672.34kB 13.95% 53.16% -1672.34kB 13.95%  compress/flate.newDeflateFast (inline)
   -1026kB  8.56% 61.72%    -1026kB  8.56%  runtime.allocm
 -521.37kB  4.35% 66.07% -8658.50kB 72.22%  runtime/pprof.(*profileBuilder).emitLocation
  514.63kB  4.29% 61.78%   514.63kB  4.29%  vendor/golang.org/x/net/http2/hpack.init
  512.12kB  4.27% 57.51%   512.12kB  4.27%  net/http.init
 -512.05kB  4.27% 61.78%  -512.05kB  4.27%  github.com/jackc/pgx/v5/pgconn/ctxwatch.(*ContextWatcher).Watch.func1
  512.01kB  4.27% 57.51%   512.01kB  4.27%  github.com/jackc/pgx/v5/pgconn.(*PgConn).makeCommandTag (inline)
  512.01kB  4.27% 53.24%   512.01kB  4.27%  internal/syscall/windows.errnoErr (inline)
         0     0% 53.24% -8137.14kB 67.87%  compress/gzip.(*Writer).Write
         0     0% 53.24%   512.01kB  4.27%  database/sql.(*DB).BeginTx
         0     0% 53.24%   512.01kB  4.27%  database/sql.(*DB).BeginTx.func1
         0     0% 53.24%   512.01kB  4.27%  database/sql.(*DB).begin
         0     0% 53.24%   512.01kB  4.27%  database/sql.(*DB).beginDC
         0     0% 53.24%   512.01kB  4.27%  database/sql.(*DB).beginDC.func1
         0     0% 53.24%   512.01kB  4.27%  database/sql.(*DB).retry
         0     0% 53.24%   512.01kB  4.27%  database/sql.ctxDriverBegin
         0     0% 53.24%   512.01kB  4.27%  database/sql.withLock
         0     0% 53.24%   512.01kB  4.27%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 53.24%   512.01kB  4.27%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 53.24%   512.01kB  4.27%  github.com/jackc/pgx/v5.(*Conn).BeginTx
         0     0% 53.24%   512.01kB  4.27%  github.com/jackc/pgx/v5.(*Conn).Exec
         0     0% 53.24%   512.01kB  4.27%  github.com/jackc/pgx/v5.(*Conn).exec
         0     0% 53.24%   512.01kB  4.27%  github.com/jackc/pgx/v5.(*Conn).execSimpleProtocol
         0     0% 53.24%   512.01kB  4.27%  github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).NextResult
         0     0% 53.24%   512.01kB  4.27%  github.com/jackc/pgx/v5/stdlib.(*Conn).BeginTx
         0     0% 53.24%   512.01kB  4.27%  github.com/mabishka/lupanova/internal/auth.WithAuth.func1
         0     0% 53.24%   512.01kB  4.27%  github.com/mabishka/lupanova/internal/compress.WithCompress.func1
         0     0% 53.24% -8658.50kB 72.22%  github.com/mabishka/lupanova/internal/handler.(*StorageServer).HandlerDelete.func1
         0     0% 53.24%   512.01kB  4.27%  github.com/mabishka/lupanova/internal/handler.(*StorageServer).HandlerPostFull
         0     0% 53.24%   512.01kB  4.27%  github.com/mabishka/lupanova/internal/logger.WithLogging.func1
         0     0% 53.24%   512.01kB  4.27%  github.com/mabishka/lupanova/internal/repository/connloader.(*ConnLoader).GetShort
         0     0% 53.24%   512.01kB  4.27%  github.com/mabishka/lupanova/internal/service.(*Server).GetShort
         0     0% 53.24%   512.01kB  4.27%  internal/poll.(*FD).Read
         0     0% 53.24%   512.01kB  4.27%  internal/poll.execIO
         0     0% 53.24%   512.01kB  4.27%  internal/syscall/windows.WSAGetOverlappedResult
         0     0% 53.24%  1762.94kB 14.70%  main.main
         0     0% 53.24%   512.01kB  4.27%  net.(*conn).Read
         0     0% 53.24%   512.01kB  4.27%  net.(*netFD).Read
         0     0% 53.24%   512.01kB  4.27%  net/http.(*conn).serve
         0     0% 53.24%   512.01kB  4.27%  net/http.(*connReader).backgroundRead
         0     0% 53.24%   512.01kB  4.27%  net/http.HandlerFunc.ServeHTTP
         0     0% 53.24%   512.01kB  4.27%  net/http.serverHandler.ServeHTTP
         0     0% 53.24%  1026.75kB  8.56%  runtime.doInit (inline)
         0     0% 53.24%  1026.75kB  8.56%  runtime.doInit1
         0     0% 53.24%     -513kB  4.28%  runtime.goexit0
         0     0% 53.24%  2789.69kB 23.27%  runtime.main
         0     0% 53.24%    -1026kB  8.56%  runtime.mcall
         0     0% 53.24%    -1026kB  8.56%  runtime.newm
         0     0% 53.24%     -513kB  4.28%  runtime.park_m
         0     0% 53.24%    -1026kB  8.56%  runtime.resetspinning
         0     0% 53.24%    -1026kB  8.56%  runtime.schedule
         0     0% 53.24%    -1026kB  8.56%  runtime.startm
         0     0% 53.24%    -1026kB  8.56%  runtime.wakep
         0     0% 53.24% -8658.50kB 72.22%  runtime/pprof.(*Profile).WriteTo
         0     0% 53.24% -8658.50kB 72.22%  runtime/pprof.(*profileBuilder).appendLocsForStack
         0     0% 53.24% -8137.14kB 67.87%  runtime/pprof.(*profileBuilder).flush
         0     0% 53.24% -8658.50kB 72.22%  runtime/pprof.writeHeap
         0     0% 53.24% -8658.50kB 72.22%  runtime/pprof.writeHeapInternal
         0     0% 53.24% -8658.50kB 72.22%  runtime/pprof.writeHeapProto