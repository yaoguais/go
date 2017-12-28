package main

import (
	"context"
	"crypto/tls"
	"flag"
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
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
	go func() {
		time.Sleep(1 * time.Second)
		test()
	}()
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
	gprcServer := GrpcServer(tls, gopts...)
	l, err := net.Listen("tcp", address)

	if err != nil {
		logrus.WithError(err).Fatal("listen")
	}

	logrus.WithField("address", address).Info("listen")

	if err := gprcServer.Serve(l); err != nil {
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

	RegisterUserServer(grpcServer, &user{})

	hsrv := health.NewServer()
	hsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, hsrv)

	return grpcServer
}

type user struct {
}

func (*user) Create(c context.Context, req *CreateRequest) (*CreateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &CreateResponse{
		Id:       1,
		Username: req.Username,
		Password: "",
		Age:      25,
	}, nil
}

func test() {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Fatal("connect")
	}

	defer cc.Close()

	requests := []*CreateRequest{
		&CreateRequest{
			Username: "hello",
			Password: "123456",
			Age:      25,
		},
		&CreateRequest{
			Username: "",
			Password: "123456",
			Age:      25,
		},
		&CreateRequest{
			Username: "hello",
			Password: "",
			Age:      25,
		},
		&CreateRequest{
			Username: "hello",
			Password: "123456",
			Age:      500,
		},
	}

	c := NewUserClient(cc)

	for _, req := range requests {
		resp, err := c.Create(context.Background(), req)
		if err != nil {
			logrus.WithError(err).Error("create user")
		} else {
			logrus.WithField("response", resp).Info("create user")
		}
	}
}
