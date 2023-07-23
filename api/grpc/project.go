package grpc

import (
	"context"

	"github.com/bitcrshr/envmgr/api/ent"
	"github.com/bitcrshr/envmgr/api/grpc/pb_grpc_environment_manager"
	pb "github.com/bitcrshr/envmgr/api/proto/compiled/go"
	"github.com/bitcrshr/envmgr/api/shared"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	project := req.Project
	if project == nil {
		err := status.Errorf(codes.InvalidArgument, "project is required")
		shared.Logger().Error("CreateProject failed", zap.Error(err))
		return nil, err
	}

	entProj, err := shared.EntClient.Project.Create().
		SetDisplayName(project.DisplayName).
		Save(ctx)
	if err != nil {
		shared.Logger().Error("CreateProject failed: couldn't save project to db", zap.Error(err))
		return nil, err
	}

	pbProj := pb_grpc_environment_manager.PbProject(entProj)

	return &pb.CreateProjectResponse{
		Project: pbProj,
	}, nil
}

func (s *server) GetOneProject(ctx context.Context, req *pb.GetOneProjectRequest) (*pb.GetOneProjectResponse, error) {
	projUuid, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Errorf(codes.InvalidArgument, "provided id [%s] was not a valid uuid", req.Id)
		shared.Logger().Error("GetOneProject failed", zap.Error(err))
		return nil, err
	}

	entProj, err := shared.EntClient.Project.Get(ctx, projUuid)
	if err != nil {
		if ent.IsNotFound(err) {
			err = status.Errorf(codes.NotFound, "project with id [%s] was not found", req.Id)
		}

		shared.Logger().Error("GetOneProject failed", zap.Error(err))
		return nil, err
	}

	pbProj := pb_grpc_environment_manager.PbProject(entProj)

	return &pb.GetOneProjectResponse{
		Project: pbProj,
	}, nil
}

func (s *server) GetAllProjects(ctx context.Context, req *pb.GetAllProjectsRequest) (*pb.GetAllProjectsResponse, error) {
	entProjs, err := shared.EntClient.Project.Query().All(ctx)
	if err != nil {
		shared.Logger().Error("GetAllProjects failed", zap.Error(err))
		return nil, err
	}

	pbProjs := pb_grpc_environment_manager.PbProjects(entProjs)

	return &pb.GetAllProjectsResponse{
		Projects: pbProjs,
	}, nil
}
