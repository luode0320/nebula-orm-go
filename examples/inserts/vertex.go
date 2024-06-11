package inserts

import (
	"nebula-orm-go/examples/models"
	"nebula-orm-go/model"
	"nebula-orm-go/orm"
)

// InsertVertex 插入一个点
//
// @Author: 罗德
// @Date: 2024/5/27
func InsertVertex(db *orm.DB) error {
	sdkVertex := models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点",
		},
		ChainKey:  "根节点",
		ParentKey: "无",
	}

	// 执行插入
	err := db.Debug().InsertVertex(sdkVertex)
	if err != nil {
		return err
	}

	sdkVertex = models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点的第一个子节点",
		},
		ChainKey:  "根节点的第一个子节点",
		ParentKey: "根节点",
	}

	// 执行插入
	err = db.Debug().InsertVertex(sdkVertex)
	if err != nil {
		return err
	}

	sdkVertex = models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点的第二个子节点",
		},
		ChainKey:  "根节点的第二个子节点",
		ParentKey: "根节点",
	}

	// 执行插入
	err = db.Debug().InsertVertex(sdkVertex)
	if err != nil {
		return err
	}

	sdkVertex = models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点的第二个子节点的第一个子节点",
		},
		ChainKey:  "根节点的第二个子节点的第一个子节点",
		ParentKey: "根节点的第二个子节点",
	}

	// 执行插入
	err = db.Debug().InsertVertex(sdkVertex)
	if err != nil {
		return err
	}

	sdkVertex = models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点的第一个子节点的第一个子节点",
		},
		ChainKey:  "根节点的第一个子节点的第一个子节点",
		ParentKey: "根节点的第一个子节点",
	}

	// 执行插入
	err = db.Debug().InsertVertex(sdkVertex)
	if err != nil {
		return err
	}

	sdkVertex = models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点的第一个子节点的第二个子节点",
		},
		ChainKey:  "根节点的第一个子节点的第二个子节点",
		ParentKey: "根节点的第一个子节点",
	}

	// 执行插入
	err = db.Debug().InsertVertex(sdkVertex)
	if err != nil {
		return err
	}

	return nil
}
