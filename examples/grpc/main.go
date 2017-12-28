package main

import (
	"crypto/tls"
	"flag"
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:5000", "server address e.g. 127.0.0.1:5000")
	flag.Parse()
}

func main() {
	go startGrpcServer(addr, nil)
	watch()
}

func watch() {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGSEGV,
		syscall.SIGTERM,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGSTOP,
	)

	for {
		s := <-c
		logrus.WithField("signal", s).Info("receive signal")

		if len(c) == 0 {
			break
		}
	}
}

func startGrpcServer(address string, tls *tls.Config, gopts ...grpc.ServerOption) {
	grpcServer := GrpcServer(tls, gopts...)
	ln, err := net.Listen("tcp", address)

	if err != nil {
		logrus.WithError(err).Fatal("listen")
	}

	logrus.WithField("address", address).Info("listen")

	if err := grpcServer.Serve(ln); err != nil {
		logrus.WithError(err).Fatal("serve")
	}
}

func GrpcServer(tls *tls.Config, gopts ...grpc.ServerOption) *grpc.Server {
	var opts []grpc.ServerOption
	if tls != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(tls)))
	}

	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_logrus.StreamServerInterceptor(logrusEntry),
		grpc_recovery.StreamServerInterceptor(),
	)))
	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_logrus.UnaryServerInterceptor(logrusEntry),
		grpc_recovery.UnaryServerInterceptor(),
	)))
	opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))
	opts = append(opts, grpc.MaxConcurrentStreams(math.MaxUint32))

	grpcServer := grpc.NewServer(append(opts, gopts...)...)

	RegisterEchoServer(grpcServer, &echo{})

	hsrv := health.NewServer()
	hsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, hsrv)

	return grpcServer
}

type echo struct {
}

func (*echo) Echo(context.Context, *Request) (*Response, error) {
	return &Response{
		Message: "world",
	}, nil
}
