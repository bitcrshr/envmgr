package grpc

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bitcrshr/envmgr/api/ent"
	"github.com/bitcrshr/envmgr/api/shared"
	pb "github.com/bitcrshr/envmgr/proto/compiled/go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50052, "environment_manager server port")
)

type server struct {
	pb.UnimplementedEnvironmentManagerServer
}

func Serve() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		shared.Logger().Fatal("failed to listen: %v", zap.Error(err))
	}

	ctx := context.Background()

	s := grpc.NewServer()

	reflection.Register(s)
	pb.RegisterEnvironmentManagerServer(s, &server{})

	errChan := make(chan error)
	stopChan := make(chan os.Signal)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	connString := os.Getenv("POSTGRES_CONN_STRING")
	shared.EntClient, err = ent.Open("postgres", connString)
	if err != nil {
		shared.Logger().Fatal("failed to serve", zap.Error(err))
		errChan <- err
	}

	defer shared.EntClient.Close()

	if err := shared.EntClient.Schema.Create(ctx); err != nil {
		shared.Logger().Fatal("failed to create db schema", zap.Error(err))
		errChan <- err
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			shared.Logger().Fatal("failed to serve", zap.Error(err))
			errChan <- err
		}
	}()

	defer func() {
		s.GracefulStop()
	}()

	select {
	case err := <-errChan:
		shared.Logger().Info("Fatal error: %v\n", zap.Error(err))
	case <-stopChan:
		shared.Logger().Info("Shutting down!")
	}
}
