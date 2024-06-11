package inserts

import (
	"nebula-orm-go/examples/models"
	"nebula-orm-go/model"
	"nebula-orm-go/orm"
)

// InsertEdge 插入一条边
//
// @Author: 罗德
// @Date: 2024/5/27
func InsertEdge(db *orm.DB) error {
	sdkEdge := models.SdkEdge{
		EModel: model.EModel{
			Src: "根节点",
			Dst: "根节点的第一个子节点",
		},
		Test: "根节点的第一个子节点",
	}

	// 执行插入
	err := db.Debug().InsertEdge(sdkEdge)
	if err != nil {
		return err
	}

	sdkEdge = models.SdkEdge{
		EModel: model.EModel{
			Src: "根节点",
			Dst: "根节点的第二个子节点",
		},
		Test: "根节点的第二个子节点",
	}

	// 执行插入
	err = db.Debug().InsertEdge(sdkEdge)
	if err != nil {
		return err
	}

	sdkEdge = models.SdkEdge{
		EModel: model.EModel{
			Src: "根节点的第二个子节点",
			Dst: "根节点的第二个子节点的第一个子节点",
		},
		Test: "根节点的第二个子节点的第一个子节点",
	}

	// 执行插入
	err = db.Debug().InsertEdge(sdkEdge)
	if err != nil {
		return err
	}

	sdkEdge = models.SdkEdge{
		EModel: model.EModel{
			Src: "根节点的第一个子节点",
			Dst: "根节点的第一个子节点的第一个子节点",
		},
		Test: "根节点的第一个子节点的第一个子节点",
	}

	// 执行插入
	err = db.Debug().InsertEdge(sdkEdge)
	if err != nil {
		return err
	}

	sdkEdge = models.SdkEdge{
		EModel: model.EModel{
			Src: "根节点的第一个子节点",
			Dst: "根节点的第一个子节点的第二个子节点",
		},
		Test: "根节点的第一个子节点的第二个子节点",
	}

	// 执行插入
	err = db.Debug().InsertEdge(sdkEdge)
	if err != nil {
		return err
	}

	return nil
}
