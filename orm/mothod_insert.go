package orm

import (
	converts2 "nebula-orm-go/converts"
	"nebula-orm-go/model"
)

// InsertVertex 根据给定的顶点模型插入一个顶点到图数据库
// 存在则默认覆盖结构体所有属性, 参数 vertex 没有被明确赋值的属性也会被覆盖为零值
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) InsertVertex(vertex model.IVertex) error {
	sql, err := converts2.ConvertToInsertVertexSql(vertex)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}

// InsertVertexIgnore 根据给定的顶点模型插入一个顶点到图数据库
// 检测待插入的 VID 是否存在，只有不存在时，才会插入，如果已经存在，不会进行修改
// 1. if not exists 仅检测 VID + Tag 的值是否相同，不会检测属性值
// 2. if not exists 会先读取一次数据是否存在，因此对性能会有明显影响
// 3. if not exists 不支持批量插入
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) InsertVertexIgnore(vertex model.IVertex) error {
	sql, err := converts2.ConvertToInsertVertexIgnoreSql(vertex)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}

// InsertVertexBatch 根据给定的顶点模型插入多个顶点到图数据库
// 存在则默认覆盖结构体所有属性, 参数 vertex 没有被明确赋值的属性也会被覆盖为零值
//
// 参数:
// vertexs (model.IVertex): 点实体切片的接口，需实现IVertex接口
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) InsertVertexBatch(vertexs []model.IVertex) error {
	sql, err := converts2.ConvertToInsertVertexBatchSql(vertexs)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}

// InsertEdge 根据给定的边模型插入一条边到图数据库
// (存在则默认覆盖结构体所有属性, 不推荐作为更新使用)
//
// 参数:
// edge (model.IEdge): 边实体的接口，需实现IEdge接口
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) InsertEdge(edge model.IEdge) error {
	sql, err := converts2.ConvertToInsertEdgeSql(edge)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}

// InsertEdgeIgnore 根据给定的边模型插入一条边到图数据库
// 检测待插入的边是否存在，只有不存在时，才会插入
// 1. if not exists 仅检测<边的类型、起始点、目的点和 rank>是否存在，不会检测属性值是否重合
// 2. if not exists 会先读取一次数据是否存在，因此对性能会有明显影响
// 3. if not exists 不支持批量插入
//
// 参数:
// edge (model.IEdge): 边实体的接口，需实现IEdge接口
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) InsertEdgeIgnore(edge model.IEdge) error {
	sql, err := converts2.ConvertToInsertEdgeIgnoreSql(edge)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}

// InsertEdgeBatch 根据给定的边模型批量插入边到图数据库
// (存在则默认覆盖结构体所有属性, 不推荐作为更新使用)
//
// 参数:
// edges (model.IEdge): 边实体切片的接口，需实现IEdge接口
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) InsertEdgeBatch(edges []model.IEdge) error {
	sql, err := converts2.ConvertToInsertEdgeBatchSql(edges)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}
