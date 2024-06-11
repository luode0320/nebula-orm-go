package orm

import (
	converts2 "nebula-orm-go/converts"
	"nebula-orm-go/dialectors"
	"nebula-orm-go/model"
)

// UpsertVertex 根据给定的顶点模型执行Upsert操作, 如果顶点存在则更新, 不存在则插入, 更新性能低于update。
//
// 参数:
// vertex (model.IVertex): 点实体的接口, 需实现IVertex接口
// set (string): 用于set更新字段的属性, 一次只能修改一个字段(例: parent_key = 'O(∩_∩)O')
// where (string): 用于补充过滤条件(例: parent_key == '无')
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) UpsertVertex(vertex model.IVertex, set string, where string) (*dialectors.ResultSet, error) {
	sql, err := converts2.ConvertToUpsertVertexSql(vertex, set, where)
	if err != nil {
		return nil, err
	}
	return db.execute(sql)
}

// UpdateVertex 根据给定的顶点模型执行Upsert操作, 如果顶点存在则更新, 不存在则忽略。
//
// 参数:
// vertex (model.IVertex): 点实体的接口, 需实现IVertex接口
// set (string): 用于set更新字段的属性, 一次只能修改一个字段(例: parent_key = 'O(∩_∩)O')
// where (string): 用于补充过滤条件(例: parent_key == '无')
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) UpdateVertex(vertex model.IVertex, set string, where string) (*dialectors.ResultSet, error) {
	sql, err := converts2.ConvertToUpdateVertexSql(vertex, set, where)
	if err != nil {
		return nil, err
	}
	return db.execute(sql)
}

// UpsertEdge 根据给定的边模型执行Upsert操作, 如果边存在则更新, 不存在则插入, 更新性能低于update。
//
// 参数:
// edge (model.IEdge): 边实体的接口, 需实现IEdge接口
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) UpsertEdge(edge model.IEdge, set string, where string) (*dialectors.ResultSet, error) {
	sql, err := converts2.ConvertToUpsertEdgeSql(edge, set, where)
	if err != nil {
		return nil, err
	}
	return db.execute(sql)
}

// UpdateEdge 根据给定的边模型执行Upsert操作, 如果边存在则更新, 不存在则忽略。
//
// 参数:
// edge (model.IEdge): 边实体的接口, 需实现IEdge接口
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) UpdateEdge(edge model.IEdge, set string, where string) (*dialectors.ResultSet, error) {
	sql, err := converts2.ConvertToUpdateEdgeSql(edge, set, where)
	if err != nil {
		return nil, err
	}
	return db.execute(sql)
}
