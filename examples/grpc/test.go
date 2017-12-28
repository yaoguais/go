package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Fatal("connect")
	}

	defer cc.Close()

	health(cc)
	echo(cc)
}

func health(cc *grpc.ClientConn) {
	c := healthpb.NewHealthClient(cc)
	req := &healthpb.HealthCheckRequest{}
	resp, err := c.Check(context.Background(), req)
	if err != nil {
		logrus.WithError(err).Error("health check")
	} else {
		logrus.WithField("response", resp).Info("health check")
	}
}

func echo(cc *grpc.ClientConn) {
	c := NewEchoClient(cc)
	req := &Request{
		Message: "hello",
	}
	resp, err := c.Echo(context.Background(), req)
	if err != nil {
		logrus.WithError(err).Error("echo")
	} else {
		logrus.WithField("response", resp).Info("echo")
	}
}
