package pb_grpc_environment_manager

import (
	"github.com/bitcrshr/envmgr/api/ent"
	"github.com/bitcrshr/envmgr/api/ent/schema/gotype"
	pb "github.com/bitcrshr/envmgr/proto/compiled/go"
	"github.com/bitcrshr/envmgr/api/shared"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func PbEnvironmentKind(entObj gotype.EnvironmentKind) pb.Environment_Kind {
	switch entObj {
	case gotype.EnvironmentKind_Development:
		return pb.Environment_KIND_DEVELOPMENT
	case gotype.EnvironmentKind_Staging:
		return pb.Environment_KIND_STAGING
	case gotype.EnvironmentKind_Production:
		return pb.Environment_KIND_PRODUCTION
	default:
		return pb.Environment_KIND_UNSPECIFIED
	}
}

func EntEnvironmentKind(pbObj pb.Environment_Kind) gotype.EnvironmentKind {
	switch pbObj {
	case pb.Environment_KIND_DEVELOPMENT:
		return gotype.EnvironmentKind_Development
	case pb.Environment_KIND_STAGING:
		return gotype.EnvironmentKind_Staging
	case pb.Environment_KIND_PRODUCTION:
		return gotype.EnvironmentKind_Production
	default:
		return gotype.EnvironmentKind_Unspecified
	}
}

func PbEnvironment(entObj *ent.Environment) *pb.Environment {
	return &pb.Environment{
		Id:   entObj.ID.String(),
		Name: entObj.Name,
		Kind: PbEnvironmentKind(entObj.Kind),
	}
}

func PbEnvironments(entObjs []*ent.Environment) []*pb.Environment {
	pbObjs := make([]*pb.Environment, 0)

	for _, entObj := range entObjs {
		pbObjs = append(pbObjs, PbEnvironment(entObj))
	}

	return pbObjs
}

func EntEnvironment(pbObj *pb.Environment) (*ent.Environment, error) {
	projUuid, err := uuid.Parse(pbObj.Id)
	if err != nil {
		shared.Logger().Error("EntEnvironment failed", zap.Error(err))
		return nil, err
	}

	return &ent.Environment{
		ID:   projUuid,
		Name: pbObj.Name,
		Kind: EntEnvironmentKind(pbObj.Kind),
	}, nil
}

func EntEnvironments(pbObjs []*pb.Environment) ([]*ent.Environment, error) {
	entObjs := make([]*ent.Environment, 0)

	for _, pbObj := range pbObjs {
		entObj, err := EntEnvironment(pbObj)
		if err != nil {
			shared.Logger().Error("EntEnvironments failed", zap.Error(err))
			return nil, err
		}

		entObjs = append(entObjs, entObj)
	}

	return entObjs, nil
}
