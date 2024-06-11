package main

import (
	"encoding/json"
	"fmt"
	"log"
	nebula_orm_go "nebula-orm-go"
	"nebula-orm-go/config"
	"nebula-orm-go/examples/inserts"
	"nebula-orm-go/examples/models"
	"nebula-orm-go/examples/sql"
	"nebula-orm-go/model"
	"nebula-orm-go/orm"

	"nebula-orm-go/dialectors"
)

func main() {
	address := fmt.Sprintf("%s:%d", "192.168.1.13", 9669)
	dialer := dialectors.MustNewNebulaDialer(config.DialerConfig{
		Addresses:       []string{address},   // Nebula服务地址列表
		Space:           sql.NebulaSpaceName, // 目标空间
		Username:        "root",              // 用户名
		Password:        "123456",            // 密码
		MaxConnPoolSize: 10,                  // 连接池大小
		InitSql:         sql.NebulaInitSql,   // 初始化sql(立刻尝试使用刚创建的空间、点、边可能会失败，因为创建是异步实现的)
	})

	db := nebula_orm_go.MustOpen(dialer, config.Config{})
	defer db.Close()

	// 立刻尝试使用刚创建的空间、点、边可能会失败，因为创建是异步实现的, 你可能需要重新启动此main程序
	//time.Sleep(10 * time.Second)

	inset(db)
	get(db)
	upsert(db)
	update(db)
	del(db)
}

// 插入数据(存在则默认更新)
//
// @Author: 罗德
// @Date: 2024/6/11
func inset(db *orm.DB) {
	// 插入点
	err := inserts.InsertVertex(db)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	err = inserts.InsertBatchVertex(db)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	// 插入边
	err = inserts.InsertEdge(db)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	err = inserts.InsertBatchEdge(db)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
}

// 删除数据
//
// @Author: 罗德
// @Date: 2024/6/11
func del(db *orm.DB) {
	// 查询测试删除边的点下级
	results, err := db.Debug().GetNextVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "测试删除根节点",
		},
	}, models.SdkEdge{}, 1)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	results.PrintResult("查询点[测试删除根节点]的下级, 并准备删除一条边")

	// 删除一条边
	err = db.Debug().DeleteEdge(models.SdkEdge{
		EModel: model.EModel{
			Src: "测试删除根节点",
			Dst: "测试删除根节点的第一个节点",
		},
	})
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}

	// 查询测试删除边的点下级
	results, err = db.Debug().GetNextVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "测试删除根节点",
		},
	}, models.SdkEdge{}, 1)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	results.PrintResult("再次查询点[测试删除根节点]的下级, 查看是否删除边[测试删除根节点] -> [测试删除根节点的第一个节点], 并准备删除多条边")

	var edges []model.IEdge
	edges = append(edges,
		models.SdkEdge{
			EModel: model.EModel{
				Src: "测试删除根节点",
				Dst: "测试删除根节点的第二个节点",
			},
		},
		models.SdkEdge{
			EModel: model.EModel{
				Src: "测试删除根节点",
				Dst: "测试删除根节点的第三个节点",
			},
		})

	// 删除多条边
	err = db.Debug().DeleteEdgeBatch(edges)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}

	// 查询测试删除边的点下级
	results, err = db.Debug().GetNextVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "测试删除根节点",
		},
	}, models.SdkEdge{}, 1)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	results.PrintResult("再次查询点[测试删除根节点]的下级, 查看是否删除点[测试删除根节点] -> [测试删除根节点的第二个节点,测试删除根节点的第三个节点]的边")

	// 查询点
	vertexmap, err := db.Debug().GetBothVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "测试删除根节点的第四个节点",
		},
	}, models.SdkEdge{}, 10)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	fmt.Println("查询点[测试删除根节点的第四个节点]与其关联的点, 准备删除点[测试删除根节点的第四个节点]:")
	for _, vertex := range vertexmap {
		// 将 map 转换为 JSON 格式的字节切片
		jsonData, err := json.MarshalIndent(vertex, "", "  ")
		if err != nil {
			fmt.Println("转换为 JSON 格式时出错:", err)
			log.Panicf("异常 -> [%s]", err.Error())
		}
		// 将 JSON 格式的字节切片转换为字符串并打印
		fmt.Println(string(jsonData))
	}
	fmt.Println()

	// 删除一个点
	err = db.Debug().DeleteVertex(models.SdkVertex{
		VModel: model.VModel{
			Vid: "测试删除根节点的第四个节点",
		},
	})
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}

	// 查询点
	vertexmap, err = db.Debug().GetBothVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "测试删除根节点的第四个节点",
		},
	}, models.SdkEdge{}, 10)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	fmt.Println("查询点[测试删除根节点的第四个节点]是否删除, 是否还有与其关联的点:")
	for _, vertex := range vertexmap {
		// 将 map 转换为 JSON 格式的字节切片
		jsonData, err := json.MarshalIndent(vertex, "", "  ")
		if err != nil {
			fmt.Println("转换为 JSON 格式时出错:", err)
			log.Panicf("异常 -> [%s]", err.Error())
		}
		// 将 JSON 格式的字节切片转换为字符串并打印
		fmt.Println(string(jsonData))
	}
	fmt.Println()

	// 查询点
	vertexmap, err = db.Debug().GetBothVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "测试删除根节点",
		},
	}, models.SdkEdge{}, 10)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	fmt.Println("查询点[测试删除根节点]剩余与其关联的点, 准备全部删除:")
	for _, vertex := range vertexmap {
		// 将 map 转换为 JSON 格式的字节切片
		jsonData, err := json.MarshalIndent(vertex, "", "  ")
		if err != nil {
			fmt.Println("转换为 JSON 格式时出错:", err)
			log.Panicf("异常 -> [%s]", err.Error())
		}
		// 将 JSON 格式的字节切片转换为字符串并打印
		fmt.Println(string(jsonData))
	}
	fmt.Println()

	// 删除多个点
	var vertexs []model.IVertex
	vertexs = append(vertexs,
		models.SdkVertex{
			VModel: model.VModel{
				Vid: "测试删除根节点",
			},
		},
		models.SdkVertex{
			VModel: model.VModel{
				Vid: "测试删除根节点的第五个节点",
			},
		})
	err = db.Debug().DeleteVertexBatch(vertexs)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}

	// 查询点
	vertexmap, err = db.Debug().GetBothVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "测试删除根节点",
		},
	}, models.SdkEdge{}, 10)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	fmt.Println("查询点[测试删除根节点]剩余与其关联的点, 是否全部删除:")
	for _, vertex := range vertexmap {
		// 将 map 转换为 JSON 格式的字节切片
		jsonData, err := json.MarshalIndent(vertex, "", "  ")
		if err != nil {
			fmt.Println("转换为 JSON 格式时出错:", err)
			log.Panicf("异常 -> [%s]", err.Error())
		}
		// 将 JSON 格式的字节切片转换为字符串并打印
		fmt.Println(string(jsonData))
	}
	fmt.Println()
}

// 更新数据(存在则更新, 不存在则插入, 更新性能低于update)
//
// @Author: 罗德
// @Date: 2024/6/11
func upsert(db *orm.DB) {
	result, err := db.Debug().UpsertVertex(models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点",
		},
	}, "parent_key = 'O(∩_∩)O'", "parent_key == '无'")
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("将点[根节点]的parent_key字段更新为[O(∩_∩)O]")

	result, err = db.UpsertVertex(models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点",
		},
	}, "parent_key = '无'", "parent_key == 'O(∩_∩)O'")
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("将点[根节点]的parent_key字段还原更新为[无]")

	result, err = db.Debug().UpsertEdge(models.SdkEdge{
		EModel: model.EModel{
			Src: "测试删除根节点",
			Dst: "测试删除根节点的第一个节点",
		},
	}, "test = 'O(∩_∩)O'", "test == '测试删除根节点的第一个节点'")
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("将边[测试删除根节点] -> [测试删除根节点的第一个节点]的test字段更新为[O(∩_∩)O]")

	result, err = db.Debug().UpsertEdge(models.SdkEdge{
		EModel: model.EModel{
			Src: "测试删除根节点",
			Dst: "测试删除根节点的第一个节点",
		},
	}, "test = '测试删除根节点的第一个节点'", "test == 'O(∩_∩)O'")
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("将边[测试删除根节点] -> [测试删除根节点的第一个节点]的test字段还原更新为[测试删除根节点的第一个节点]")
}

// 更新数据(存在则更新, 不存在忽略)
//
// @Author: 罗德
// @Date: 2024/6/11
func update(db *orm.DB) {
	result, err := db.Debug().UpdateVertex(models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点",
		},
	}, "parent_key = 'O(∩_∩)O'", "parent_key == '无'")
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("将点[根节点]的parent_key字段更新为[O(∩_∩)O]")

	result, err = db.UpdateVertex(models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点",
		},
	}, "parent_key = '无'", "parent_key == 'O(∩_∩)O'")
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("将点[根节点]的parent_key字段还原更新为[无]")

	result, err = db.Debug().UpdateEdge(models.SdkEdge{
		EModel: model.EModel{
			Src: "测试删除根节点",
			Dst: "测试删除根节点的第一个节点",
		},
	}, "test = 'O(∩_∩)O'", "test == '测试删除根节点的第一个节点'")
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("将边[测试删除根节点] -> [测试删除根节点的第一个节点]的test字段更新为[O(∩_∩)O]")

	result, err = db.Debug().UpdateEdge(models.SdkEdge{
		EModel: model.EModel{
			Src: "测试删除根节点",
			Dst: "测试删除根节点的第一个节点",
		},
	}, "test = '测试删除根节点的第一个节点'", "test == 'O(∩_∩)O'")
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("将边[测试删除根节点] -> [测试删除根节点的第一个节点]的test字段还原更新为[测试删除根节点的第一个节点]")
}

// 查询数据
//
// @Author: 罗德
// @Date: 2024/6/11
func get(db *orm.DB) {
	// 查询点
	result, err := db.Debug().GetVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点",
		},
	})
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("查询点[根节点]")

	// 查询点下级
	maps, err := db.Debug().GetNextVertexMapByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点",
		},
	}, models.SdkEdge{}, 3)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	fmt.Println("查询点[根节点]的下级:")
	for _, vertex := range maps {
		// 将 map 转换为 JSON 格式的字节切片
		jsonData, err := json.MarshalIndent(vertex, "", "  ")
		if err != nil {
			fmt.Println("转换为 JSON 格式时出错:", err)
			log.Panicf("异常 -> [%s]", err.Error())
		}
		// 将 JSON 格式的字节切片转换为字符串并打印
		fmt.Println(string(jsonData))
	}
	fmt.Println()

	// 查询点上级
	results, err := db.Debug().GetUpVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点的第一个子节点",
		},
	}, models.SdkEdge{}, 3)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	results.PrintResult("查询点[根节点的第一个子节点]的上级")

	// 查询点上+下级点
	// (包含原点 + 原点的下级, 原点的下级的下级 + 原点的上级, 原点的上级的上级)
	vertexs, err := db.Debug().GetBothVertexByVid(models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点的第一个子节点",
		},
	}, models.SdkEdge{}, 2)
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	fmt.Println("查询点[根节点的第一个子节点]的上+下级点:")
	for _, vertex := range vertexs {
		// 将 map 转换为 JSON 格式的字节切片
		jsonData, err := json.MarshalIndent(vertex, "", "  ")
		if err != nil {
			fmt.Println("转换为 JSON 格式时出错:", err)
			log.Panicf("异常 -> [%s]", err.Error())
		}
		// 将 JSON 格式的字节切片转换为字符串并打印
		fmt.Println(string(jsonData))
	}
	fmt.Println()
}
