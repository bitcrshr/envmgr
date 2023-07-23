package pb_grpc_environment_manager

import (
	"github.com/bitcrshr/envmgr/api/ent"
	"github.com/bitcrshr/envmgr/api/shared"
	pb "github.com/bitcrshr/envmgr/proto/compiled/go"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func PbVariable(entObj *ent.Variable) *pb.Variable {
	return &pb.Variable{
		Id:    entObj.ID.String(),
		Key:   entObj.Key,
		Value: entObj.Value,
	}
}

func PbVariables(entObjs []*ent.Variable) []*pb.Variable {
	pbObjs := make([]*pb.Variable, 0)

	for _, entObj := range entObjs {
		pbObjs = append(pbObjs, PbVariable(entObj))
	}

	return pbObjs
}

func EntVariable(pbObj *pb.Variable) (*ent.Variable, error) {
	projUuid, err := uuid.Parse(pbObj.Id)
	if err != nil {
		shared.Logger().Error("EntVariable failed", zap.Error(err))
		return nil, err
	}

	return &ent.Variable{
		ID:    projUuid,
		Key:   pbObj.Key,
		Value: pbObj.Value,
	}, nil
}

func EntVariables(pbObjs []*pb.Variable) ([]*ent.Variable, error) {
	entObjs := make([]*ent.Variable, 0)

	for _, pbObj := range pbObjs {
		entObj, err := EntVariable(pbObj)
		if err != nil {
			shared.Logger().Error("EntVariables failed", zap.Error(err))
			return nil, err
		}

		entObjs = append(entObjs, entObj)
	}

	return entObjs, nil
}
