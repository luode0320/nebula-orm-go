package orm

import (
	"fmt"
	"math"
	"nebula-orm-go/constants"
	"nebula-orm-go/dialectors"
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
)

// GetVertexByVid 根据ID匹配单个用户顶点
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
//
// 返回:
// ResultSet: 成功返回查询结构。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) GetVertexByVid(vertex model.IVertex) (*dialectors.ResultSet, error) {
	vid := utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy())
	clause, err := utils.GetVClause(vertex, vertex.TagName())
	if err != nil {
		return nil, err
	}
	db.sql = fmt.Sprintf("match(v:%s) where id(v)==%s return %s", vertex.TagName(), vid, clause)
	result, err := db.ReturnRow()
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetNextVertexByVid 根据点ID、边名称、匹配下级点, 仅返回下级列表数据
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
// edge (model.IEdge): 边实体的接口，需实现IEdge接口
// level (int): 需要查询的层级
//
// 返回:
// ResultSet: 成功返回下级列表数据。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) GetNextVertexByVid(vertex model.IVertex, edge model.IEdge, level int) (*dialectors.ResultSet, error) {
	vid := utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy())
	clause, err := utils.GetVClause(vertex, vertex.TagName())
	if err != nil {
		return nil, err
	}
	db.sql = fmt.Sprintf("match p=(v)<-[:%s*1..%d]-(n) where id(n)==%s return %s", edge.EdgeName(), level, vid, clause)
	result, err := db.Limit(db.limit).ReturnRow()
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetNextVertexMapByVid 根据点ID、边名称、匹配下级点, 按照层级map返回包含 下一级...下level级 的层级列表
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
// edge (model.IEdge): 边实体的接口，需实现IEdge接口
// level (int): 需要查询的层级
//
// 返回:
// map: 成功按照层级map返回包含 下一级...下level级 的层级列表。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) GetNextVertexMapByVid(vertex model.IVertex, edge model.IEdge, level int) ([][]map[string]interface{}, error) {
	vid := utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy())
	// 查询下级, 因为get subgraph默认返回原点本身作为第一层, 所以查询level至少+1
	db.sql = fmt.Sprintf(`get subgraph with prop %d steps from %s out %s yield vertices as %s`, level+1, vid, edge.EdgeName(), constants.V)
	resultOut, err := db.ReturnRow()
	if err != nil {
		return nil, err
	}
	vertexs, err := resultOut.ToVertexs()
	// 移除当前查询的点本身
	if len(vertexs) <= 1 {
		return make([][]map[string]interface{}, 0), err
	}
	return vertexs[1:], err
}

// GetUpVertexByVid 根据点ID、边、匹配上级点, 仅返回上级列表数据
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
// edge (model.IEdge): 边实体的接口，需实现IEdge接口
// level (int): 需要查询的层级
//
// 返回:
// ResultSet: 成功返回上级列表数据。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) GetUpVertexByVid(vertex model.IVertex, edge model.IEdge, level int) (*dialectors.ResultSet, error) {
	vid := utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy())
	clause, err := utils.GetVClause(vertex, vertex.TagName())
	if err != nil {
		return nil, err
	}
	db.sql = fmt.Sprintf("match p=(n)<-[:%s*1..%d]-(v) where id(n)==%s return %s", edge.EdgeName(), level, vid, clause)
	result, err := db.Limit(db.limit).ReturnRow()
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetUpVertexMapByVid 根据点ID、边、匹配上级点, 按照层级map返回包含 上一级...上level级 的层级列表
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
// edge (model.IEdge): 边实体的接口，需实现IEdge接口
// level (int): 需要查询的层级
//
// 返回:
// map: 成功按照层级map返回包含 上一级...上level级 的层级列表。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) GetUpVertexMapByVid(vertex model.IVertex, edge model.IEdge, level int) ([][]map[string]interface{}, error) {
	vid := utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy())
	// 查询上级, 因为get subgraph默认返回原点本身作为第一层, 所以查询level至少+1
	db.sql = fmt.Sprintf(`get subgraph with prop %d steps from %s in %s yield vertices as %s`, level+1, vid, edge.EdgeName(), constants.V)
	resultOut, err := db.ReturnRow()
	if err != nil {
		return nil, err
	}
	vertexs, err := resultOut.ToVertexs()
	// 移除当前查询的点本身
	if len(vertexs) <= 1 {
		return make([][]map[string]interface{}, 0), err
	}
	return vertexs[1:], err
}

// GetBothAllVertexByVid 根据点ID匹配原点+上+下级点
// (包含原点 + 原点的下级, 原点的下级的下级 + 原点的上级, 原点的上级的上级, 原点的上级的下级)
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
// edge (model.IEdge): 边实体的接口，需实现IEdge接口
// level (int): 需要查询的层级
//
// 返回:
// map: 成功返回查询点map结构, 按照上下级数组层级返回。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/5
func (db *DB) GetBothAllVertexByVid(vertex model.IVertex, edge model.IEdge, level int) ([][]map[string]interface{}, error) {
	vid := utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy())
	// 查询上下级, 因为get subgraph默认返回原点本身作为第一层, 所以查询level至少+1
	db.sql = fmt.Sprintf("get subgraph with prop %d steps from %s both %s yield vertices as %s", level+1, vid, edge.EdgeName(), constants.V)
	result, err := db.ReturnRow()
	if err != nil {
		return nil, err
	}

	vertexs, err := result.ToVertexs()
	if len(vertexs) > 0 && vertexs[0] == nil {
		vertexs = make([][]map[string]interface{}, 0)
	}
	return vertexs, err
}

// GetBothVertexByVid 根据点ID匹配原点+上+下级点
// (包含原点 + 原点的下级, 原点的下级的下级 + 原点的上级, 原点的上级的上级)
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
// edge (model.IEdge): 边实体的接口，需实现IEdge接口
// level (int): 需要查询的层级
//
// 返回:
// map: 成功返回查询点map结构, 按照上下级数组层级返回。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/5
func (db *DB) GetBothVertexByVid(vertex model.IVertex, edge model.IEdge, level int) ([][]map[string]interface{}, error) {
	vid := utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy())
	// 查询上级, 因为get subgraph默认返回原点本身作为第一层, 所以查询level至少+1
	db.sql = fmt.Sprintf(`get subgraph with prop %d steps from %s out %s yield vertices as %s`, level+1, vid, edge.EdgeName(), constants.V)
	resultOut, err := db.ReturnRow()
	if err != nil {
		return nil, err
	}
	vertexsOut, err := resultOut.ToVertexs()
	if err != nil {
		return nil, err
	}

	// 查询下级, 因为get subgraph默认返回原点本身作为第一层, 所以查询level至少+1
	db.sql = fmt.Sprintf(`get subgraph with prop %d steps from %s in %s yield vertices as %s`, level+1, vid, edge.EdgeName(), constants.V)
	resultIn, err := db.ReturnRow()
	if err != nil {
		return nil, err
	}
	vertexsIn, err := resultIn.ToVertexs()
	if err != nil {
		return nil, err
	}

	// 合并两个上下级 [][]map[string]interface{} 结构
	var mergedData = make([][]map[string]interface{}, int(math.Max(float64(len(vertexsOut)), float64(len(vertexsIn)))))
	for i := range mergedData {
		if len(vertexsOut) >= i+1 {
			mergedData[i] = append(mergedData[i], vertexsOut[i]...)
		}
		if len(vertexsIn) >= i+1 && i != 0 {
			mergedData[i] = append(mergedData[i], vertexsIn[i]...)
		}
	}
	if len(mergedData) > 0 && mergedData[0] == nil {
		mergedData = make([][]map[string]interface{}, 0)
	}
	return mergedData, nil
}
