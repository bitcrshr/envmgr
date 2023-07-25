package grpc

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/bitcrshr/envmgr/api/ent"
	"github.com/bitcrshr/envmgr/api/shared"
	pb "github.com/bitcrshr/envmgr/proto/compiled/go"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	port     = flag.Int("port", 50052, "environment_manager server port")
	jwkCache *jwk.Cache
)

const (
	JWKS_URI string = "https://envmgr-dev.us.auth0.com/.well-known/jwks"
)

type server struct {
	pb.UnimplementedEnvironmentManagerServer
}

func Serve(entMigrationDoneChannel chan bool) {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		shared.Logger().Fatal("failed to listen: %v", zap.Error(err))
	}

	ctx := context.Background()
	jwkCache = jwk.NewCache(ctx)

	jwkCache.Register(
		JWKS_URI,
		jwk.WithMinRefreshInterval(15*time.Minute),
	)
	_, err = jwkCache.Refresh(ctx, JWKS_URI)
	if err != nil {
		shared.Logger().Fatal("failed to do initial jwk refresh: %v", zap.Error(err))
	}

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

	entMigrationDoneChannel <- true

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

func authCheck(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	jwks, err := jwkCache.Get(ctx, JWKS_URI)
	if err != nil {
		return nil, err
	}

	tok, err := jwt.Parse([]byte(token), jwt.WithKeySet(jwks), jwt.WithValidate(true))
	if err != nil {
		err := status.Error(codes.Unauthenticated, "invalid token")
		return nil, err
	}

	ctx = context.WithValue(ctx, "user_id", tok.Subject())

	return ctx, nil
}
