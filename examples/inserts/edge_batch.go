package inserts

import (
	"nebula-orm-go/examples/models"
	"nebula-orm-go/model"
	"nebula-orm-go/orm"
)

// InsertBatchEdge 批量插入边
//
// @Author: 罗德
// @Date: 2024/5/27
func InsertBatchEdge(db *orm.DB) error {
	var sdkEdges []model.IEdge
	sdkEdges = append(sdkEdges,
		models.SdkEdge{
			EModel: model.EModel{
				Src: "测试删除根节点",
				Dst: "测试删除根节点的第一个节点",
			},
			Test: "测试删除根节点的第一个节点",
		},
		models.SdkEdge{
			EModel: model.EModel{
				Src: "测试删除根节点",
				Dst: "测试删除根节点的第二个节点",
			},
			Test: "测试删除根节点的第二个节点",
		},
		models.SdkEdge{
			EModel: model.EModel{
				Src: "测试删除根节点",
				Dst: "测试删除根节点的第三个节点",
			},
			Test: "测试删除根节点的第三个节点",
		},
		models.SdkEdge{
			EModel: model.EModel{
				Src: "测试删除根节点",
				Dst: "测试删除根节点的第四个节点",
			},
			Test: "测试删除根节点的第四个节点",
		},
		models.SdkEdge{
			EModel: model.EModel{
				Src: "测试删除根节点",
				Dst: "测试删除根节点的第五个节点",
			},
			Test: "测试删除根节点的第五个节点",
		})

	err := db.Debug().InsertEdgeBatch(sdkEdges)
	if err != nil {
		return err
	}

	return nil
}
