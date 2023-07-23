package grpc

import (
	"context"

	"github.com/bitcrshr/envmgr/api/ent"
	"github.com/bitcrshr/envmgr/api/ent/environment"
	"github.com/bitcrshr/envmgr/api/ent/project"
	"github.com/bitcrshr/envmgr/api/grpc/pb_grpc_environment_manager"
	pb "github.com/bitcrshr/envmgr/api/proto/compiled/go"
	"github.com/bitcrshr/envmgr/api/shared"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) CreateEnvironment(ctx context.Context, req *pb.CreateEnvironmentRequest) (*pb.CreateEnvironmentResponse, error) {
	environment := req.Environment
	if environment == nil {
		err := status.Errorf(codes.InvalidArgument, "environment is required")
		shared.Logger().Error("CreateEnvironment failed", zap.Error(err))
		return nil, err
	}

	entEnv, err := shared.EntClient.Environment.Create().
		SetName(environment.Name).
		SetKind(pb_grpc_environment_manager.EntEnvironmentKind(environment.Kind)).
		Save(ctx)
	if err != nil {
		shared.Logger().Error("CreateEnvironment failed: couldn't save Environment to db", zap.Error(err))
		return nil, err
	}

	pbEnv := pb_grpc_environment_manager.PbEnvironment(entEnv)

	return &pb.CreateEnvironmentResponse{
		Environment: pbEnv,
	}, nil
}

func (s *server) GetOneEnvironment(ctx context.Context, req *pb.GetOneEnvironmentRequest) (*pb.GetOneEnvironmentResponse, error) {
	envUuid, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Errorf(codes.InvalidArgument, "provided id [%s] was not a valid uuid", req.Id)
		shared.Logger().Error("GetOneEnvironment failed", zap.Error(err))
		return nil, err
	}

	entEnv, err := shared.EntClient.Environment.Get(ctx, envUuid)
	if err != nil {
		if ent.IsNotFound(err) {
			err = status.Errorf(codes.NotFound, "Environment with id [%s] was not found", req.Id)
		}

		shared.Logger().Error("GetOneEnvironment failed", zap.Error(err))
		return nil, err
	}

	pbEnv := pb_grpc_environment_manager.PbEnvironment(entEnv)

	return &pb.GetOneEnvironmentResponse{
		Environment: pbEnv,
	}, nil
}

func (s *server) GetAllEnvironments(ctx context.Context, req *pb.GetAllEnvironmentsRequest) (*pb.GetAllEnvironmentsResponse, error) {
	projUuid, err := uuid.Parse(req.ProjectId)
	if err != nil {
		err := status.Errorf(codes.InvalidArgument, "provided id [%s] was not a valid uuid", req.ProjectId)
		shared.Logger().Error("GetAllEnvironments failed", zap.Error(err))
		return nil, err
	}

	entEnvs, err := shared.EntClient.Environment.Query().
		Where(environment.HasProjectWith(project.ID(projUuid))).
		All(ctx)
	if err != nil {
		shared.Logger().Error("GetAllEnvironments failed", zap.Error(err))
		return nil, err
	}

	pbEnvs := pb_grpc_environment_manager.PbEnvironments(entEnvs)

	return &pb.GetAllEnvironmentsResponse{
		Environments: pbEnvs,
	}, nil
}
