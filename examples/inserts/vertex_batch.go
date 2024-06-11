package inserts

import (
	"nebula-orm-go/examples/models"
	"nebula-orm-go/model"
	"nebula-orm-go/orm"
)

// InsertBatchVertex 插入多点
//
// @Author: 罗德
// @Date: 2024/5/27
func InsertBatchVertex(db *orm.DB) error {
	var sdkVertexs []model.IVertex
	sdkVertexs = append(sdkVertexs,
		models.SdkVertex{
			VModel: model.VModel{
				Vid: "测试删除根节点",
			},
			ChainKey:  "测试删除根节点",
			ParentKey: "无",
		},
		models.SdkVertex{
			VModel: model.VModel{
				Vid: "测试删除根节点的第一个节点",
			},
			ChainKey:  "测试删除根节点的第一个节点",
			ParentKey: "测试删除根节点",
		},
		models.SdkVertex{
			VModel: model.VModel{
				Vid: "测试删除根节点的第二个节点",
			},
			ChainKey:  "测试删除根节点的第二个节点",
			ParentKey: "测试删除根节点",
		},
		models.SdkVertex{
			VModel: model.VModel{
				Vid: "测试删除根节点的第三个节点",
			},
			ChainKey:  "测试删除根节点的第三个节点",
			ParentKey: "测试删除根节点",
		},
		models.SdkVertex{
			VModel: model.VModel{
				Vid: "测试删除根节点的第四个节点",
			},
			ChainKey:  "测试删除根节点的第四个节点",
			ParentKey: "测试删除根节点",
		},
		models.SdkVertex{
			VModel: model.VModel{
				Vid: "测试删除根节点的第五个节点",
			},
			ChainKey:  "测试删除根节点的第五个节点",
			ParentKey: "测试删除根节点",
		})

	err := db.Debug().InsertVertexBatch(sdkVertexs)
	if err != nil {
		return err
	}
	return nil
}
