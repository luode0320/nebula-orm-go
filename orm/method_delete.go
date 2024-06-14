package orm

import (
	converts2 "nebula-orm-go/converts"
	"nebula-orm-go/model"
)

// DeleteVertex 删除一个点, 及其关联的边
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
//
// @Author: 罗德
// @Date: 2024/6/6
func (db *DB) DeleteVertex(vertex model.IVertex) error {
	sql, err := converts2.ConvertToDeleteVertexSql(vertex)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}

// DeleteVertexBatch 删除多个点, 及其关联的边
//
// 参数:
// vertexs ([]model.IVertex): 点实体切片的接口，需实现IVertex接口
//
// @Author: 罗德
// @Date: 2024/6/6
func (db *DB) DeleteVertexBatch(vertexs []model.IVertex) error {
	sql, err := converts2.ConvertToDeleteVertexBatchSql(vertexs)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}

// DeleteEdge 删除一条边
//
// 参数:
// edge (model.IEdge): 边实体的接口，需实现IEdge接口。
//
// @Author: 罗德
// @Date: 2024/6/6
func (db *DB) DeleteEdge(edge model.IEdge) error {
	sql, err := converts2.ConvertToDeleteEdgeSql(edge)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}

// DeleteEdgeBatch 删除多条边
//
// 参数:
// edges ([]model.IEdge): 边实体切片的接口，需实现IEdge接口。
//
// @Author: 罗德
// @Date: 2024/6/6
func (db *DB) DeleteEdgeBatch(edges []model.IEdge) error {
	sql, err := converts2.ConvertToDeleteEdgeBatchSql(edges)
	if err != nil {
		return err
	}
	_, err = db.execute(sql)
	return err
}
