package pb_grpc_environment_manager

import (
	"github.com/bitcrshr/envmgr/api/ent"
	pb "github.com/bitcrshr/envmgr/api/proto/compiled/go"
	"github.com/bitcrshr/envmgr/api/shared"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func PbProject(entObj *ent.Project) *pb.Project {
	return &pb.Project{
		Id:          entObj.ID.String(),
		DisplayName: entObj.DisplayName,
	}
}

func PbProjects(entObjs []*ent.Project) []*pb.Project {
	pbObjs := make([]*pb.Project, 0)

	for _, entObj := range entObjs {
		pbObjs = append(pbObjs, PbProject(entObj))
	}

	return pbObjs
}

func EntProject(pbObj *pb.Project) (*ent.Project, error) {
	projUuid, err := uuid.Parse(pbObj.Id)
	if err != nil {
		shared.Logger().Error("EntProject failed", zap.Error(err))
		return nil, err
	}

	return &ent.Project{
		ID:          projUuid,
		DisplayName: pbObj.DisplayName,
	}, nil
}

func EntProjects(pbObjs []*pb.Project) ([]*ent.Project, error) {
	entObjs := make([]*ent.Project, 0)

	for _, pbObj := range pbObjs {
		entObj, err := EntProject(pbObj)
		if err != nil {
			shared.Logger().Error("EntProjects failed", zap.Error(err))
			return nil, err
		}

		entObjs = append(entObjs, entObj)
	}

	return entObjs, nil
}
