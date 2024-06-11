package querys

import (
	"nebula-orm-go/model"

	"nebula-orm-go/examples/models"
)

var (
	userExample = &models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点",
		},
		ChainKey:  "根节点",
		ParentKey: "无",
	}

	voteExample = &models.SdkEdge{
		EModel: model.EModel{
			Src: "根节点",
			Dst: "根节点的第一个子节点",
		},
	}
)
