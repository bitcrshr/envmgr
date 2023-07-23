package grpc

import (
	"context"

	"github.com/bitcrshr/envmgr/api/ent"
	"github.com/bitcrshr/envmgr/api/ent/environment"
	"github.com/bitcrshr/envmgr/api/ent/predicate"
	"github.com/bitcrshr/envmgr/api/ent/project"
	"github.com/bitcrshr/envmgr/api/ent/variable"
	"github.com/bitcrshr/envmgr/api/grpc/pb_grpc_environment_manager"
	"github.com/bitcrshr/envmgr/api/shared"
	pb "github.com/bitcrshr/envmgr/proto/compiled/go"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) CreateVariable(ctx context.Context, req *pb.CreateVariableRequest) (*pb.CreateVariableResponse, error) {
	variable := req.Variable
	if variable == nil {
		err := status.Errorf(codes.InvalidArgument, "variable is required")
		shared.Logger().Error("CreateVariable failed", zap.Error(err))
		return nil, err
	}

	entVar, err := shared.EntClient.Variable.Create().
		SetKey(variable.Key).
		SetValue(variable.Value).
		Save(ctx)
	if err != nil {
		shared.Logger().Error("CreateVariable failed: couldn't save Variable to db", zap.Error(err))
		return nil, err
	}

	pbVar := pb_grpc_environment_manager.PbVariable(entVar)

	return &pb.CreateVariableResponse{
		Variable: pbVar,
	}, nil
}

func (s *server) GetOneVariable(ctx context.Context, req *pb.GetOneVariableRequest) (*pb.GetOneVariableResponse, error) {
	varUuid, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Errorf(codes.InvalidArgument, "provided id [%s] was not a valid uuid", req.Id)
		shared.Logger().Error("GetOneVariable failed", zap.Error(err))
		return nil, err
	}

	entVar, err := shared.EntClient.Variable.Get(ctx, varUuid)
	if err != nil {
		if ent.IsNotFound(err) {
			err = status.Errorf(codes.NotFound, "Variable with id [%s] was not found", req.Id)
		}

		shared.Logger().Error("GetOneVariable failed", zap.Error(err))
		return nil, err
	}

	pbVar := pb_grpc_environment_manager.PbVariable(entVar)

	return &pb.GetOneVariableResponse{
		Variable: pbVar,
	}, nil
}

func (s *server) QueryVariables(ctx context.Context, req *pb.QueryVariablesRequest) (*pb.QueryVariablesResponse, error) {
	projUuid, err := uuid.Parse(req.ProjectId)
	if err != nil {
		err := status.Errorf(codes.InvalidArgument, "provided project id [%s] was not a valid uuid", req.ProjectId)
		shared.Logger().Error("QueryVariables failed", zap.Error(err))
		return nil, err
	}

	filters := make([]predicate.Variable, 0)
	filters = append(
		filters,
		variable.HasEnvironmentWith(environment.HasProjectWith(project.ID(projUuid))),
	)

	if req.EnvironmentKind != pb.Environment_KIND_UNSPECIFIED {
		filters = append(
			filters,
			variable.HasEnvironmentWith(
				environment.KindEQ(
					pb_grpc_environment_manager.EntEnvironmentKind(req.EnvironmentKind),
				),
			),
		)
	}

	entVars, err := shared.EntClient.Variable.Query().
		Where(filters...).
		All(ctx)

	if err != nil {
		shared.Logger().Error("QueryVariables failed", zap.Error(err))
		return nil, err
	}

	pbVars := pb_grpc_environment_manager.PbVariables(entVars)

	return &pb.QueryVariablesResponse{
		Variables: pbVars,
	}, nil
}
