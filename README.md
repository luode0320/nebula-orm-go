# 简介
这是一个操作nebula图数据库的orm框架。

# 示例
请优先参考[examples](examples)示例: 配置nabula图数据库连接后可立即运行测试, 帮助你更快了解orm的使用
```log 
# 部分日志仅供参考

2024/06/14 10:38:29 [INFO] insert vertex test_vertex(chain_key,parent_key) values '根节点':('根节点','无')
2024/06/14 10:38:29 [INFO] insert edge test_edge(test) values '根节点' -> '根节点的第一个子节点':('根节点的第一个子节点')
2024/06/14 10:38:29 [INFO] delete vertex '测试删除根节点的第四个节点' with edge
2024/06/14 10:38:29 [INFO] delete edge test_edge '测试删除根节点' -> '测试删除根节点的第一个节点'
2024/06/14 10:38:29 [INFO] update edge on test_edge '测试删除根节点' -> '测试删除根节点的第一个节点' set test = 'O(∩_∩)O' when test == '测试删除根节点的第一个节点' yield test as test
将边[测试删除根节点] -> [测试删除根节点的第一个节点]的test字段更新为[O(∩_∩)O]:
[
  {
    "test": "O(∩_∩)O"
  }
]
2024/06/14 10:38:29 [INFO] match(v:test_vertex) where id(v)=='根节点' return v.test_vertex.chain_key as chain_key,v.test_vertex.parent_key as parent_key
查询点[根节点]:
[
  {
    "chain_key": "根节点",
    "parent_key": "无"    
  }                       
]
```

## [连接](examples%2Fmain.go)

```go
package main
import (
	nebula_orm_go "nebula-orm-go"
	"nebula-orm-go/config"
	"nebula-orm-go/examples/inserts"
	"nebula-orm-go/examples/models"
	"nebula-orm-go/examples/sql"
	"nebula-orm-go/model"
	"nebula-orm-go/orm"
	"nebula-orm-go/dialectors"
)

// 测试方法入口
//
// @Author: 罗德
// @Date: 2024/6/11
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
}
```

## [创建空间、点、边结构](examples%2Fsql%2Fnebula.go)
```sql
create space if not exists space_luode(vid_type=fixed_string(64));

use space_luode;

create tag if not exists test_vertex(
    chain_key string,
    parent_key string
);

create edge if not exists test_edge(test string);
```

## 创建实体
[点](examples%2Fmodels%2Fsdk_vertex.go):
```go
package models

import (
	"nebula-orm-go/examples/sql"
	"nebula-orm-go/model"
)

// 确保满足IVertex定义的所有要求
//
// @Author: 罗德
// @Date: 2024/5/27
var _ model.IVertex = new(SdkVertex)

// SdkVertex 点结构体
//
// @Author: 罗德
// @Date: 2024/5/27
type SdkVertex struct {
	model.VModel
	ChainKey  string `json:"chain_key" nebula:"chain_key"`
	ParentKey string `json:"parent_key" nebula:"parent_key"`
}

// TagName 点名称, 必须实现
//
// @Author: 罗德
// @Date: 2024/5/27
func (v SdkVertex) TagName() string {
	return sql.NebulaVertexName
}

```
[边](examples%2Fmodels%2Fsdk_edge.go):
```go
package models

import (
	"nebula-orm-go/examples/sql"
	"nebula-orm-go/model"
)

// 确保满足IEdge定义的所有要求
//
// @Author: 罗德
// @Date: 2024/5/27
var _ model.IEdge = new(SdkEdge)

// SdkEdge 边结构体
//
// @Author: 罗德
// @Date: 2024/5/27
type SdkEdge struct {
	model.EModel
	Test string `json:"test" nebula:"test"`
}

// EdgeName 边名称, 必须实现
//
// @Author: 罗德
// @Date: 2024/5/27
func (v SdkEdge) EdgeName() string {
	return sql.NebulaEdgeName
}
```

## [新增点](examples%2Finserts%2Fvertex.go)

```go
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
```

## [新增边](examples%2Fmain.go)
```go
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
```

## [删除点](examples%2Fmain.go)

```go
	// 删除一个点
	err = db.Debug().DeleteVertex(models.SdkVertex{
		VModel: model.VModel{
			Vid: "测试删除根节点的第四个节点",
		},
	})
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
```
## [删除边](examples%2Fmain.go)
```go
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
```

## [更新点](examples%2Fmain.go)
```go
	result, err := db.Debug().UpdateVertex(models.SdkVertex{
		VModel: model.VModel{
			Vid: "根节点",
		},
	}, "parent_key = 'O(∩_∩)O'", "parent_key == '无'")
	if err != nil {
		log.Panicf("异常 -> [%s]", err.Error())
	}
	result.PrintResult("将点[根节点]的parent_key字段更新为[O(∩_∩)O]")
```

## [更新边](examples%2Fmain.go)
```go
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
```

## [查询点](examples%2Fmain.go)
```go
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
```

## [查询点上/下级](examples%2Fmain.go)
```go
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
```

## [更多参考](orm)
-[orm](orm)
  - [method_insert.go](orm%2Fmothod_insert.go)
  - [method_delete.go](orm%2Fmothod_delete.go)
  - [method_update.go](orm%2Fmothod_update.go)
  - [method_select.go](orm%2Fmothod_select.go)