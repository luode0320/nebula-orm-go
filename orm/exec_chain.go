// Package orm 链式调用: db.Debug().From().Go().Execute()
//
// @Author: 罗德
// @Date: 2024/5/27
package orm

import (
	"fmt"
	"nebula-orm-go/constants"
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings"
)

// Debug 设置调试模式，当执行nGQL时会打印查询语句。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Debug() (tx *DB) {
	tx = db.getInstance()
	tx.debug = true
	return
}

// Match 语句使用的路径类型是trail，这意味着点可以重复出现，但边不能重复
// MATCH (v)-[e:follow*1..2]->(v2) WHERE id(v) == "player100" RETURN id(v2) AS destination; # 查询 player100 1~2 跳内的朋友。
// https://docs.nebula-graph.com.cn/3.6.0/3.ngql-guide/7.general-query-statements/2.match/
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Match(match string) (tx *DB) {
	tx = db.getInstance()
	tx.sql += fmt.Sprintf(" match %d ", match)
	return
}

// Go 语句采用的路径类型是walk，即遍历时点和边都可以重复
// GO 1 TO 2 STEPS FROM "player100" OVER follow YIELD dst(edge) AS destination; # 查询 player100 1~2 跳内的朋友。
// https://docs.nebula-graph.com.cn/3.6.0/3.ngql-guide/7.general-query-statements/3.go/
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Go(level int) (tx *DB) {
	tx = db.getInstance()

	if level > 1 {
		// 遍历M~N跳的边。如果M为0，输出结果和M为1相同，即GO 0 TO 2和GO 1 TO 2是相同的。
		tx.sql += fmt.Sprintf(" go 1 TO %d step ", level)
	}

	return
}

// From 指定查询的起始顶点，接受一系列实现了IVertex接口的实例。
// GO 1 TO 2 STEPS FROM "player100" OVER follow YIELD dst(edge) AS destination; # 查询 player100 1~2 跳内的朋友。
// FIND SHORTEST PATH FROM "player102" TO "team204" OVER * YIELD path AS p; # 查找并返回 player102 到 team204 的最短路径。
// GET SUBGRAPH 1 STEPS FROM "player101" YIELD VERTICES AS nodes, EDGES AS relationships; # 查询从点player101开始、0~1 跳、所有 Edge type 的子图。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) From(vs ...model.IVertex) (tx *DB) {
	tx = db.getInstance()

	vids := make([]string, len(vs))
	for i, v := range vs {
		vids[i] = utils.GetVidWithPolicy(v.GetVid(), v.GetPolicy())
	}

	if tx.sql == "" {
		tx.sql += " go "
	}
	tx.sql += fmt.Sprintf(" from %s ", strings.Join(vids, ","))

	return
}

// Over 指定查询中要遍历的边类型，接受一系列实现了IEdge接口的实例。
// GO 1 TO 2 STEPS FROM "player100" OVER follow YIELD dst(edge) AS destination; # 查询 player100 1~2 跳内的朋友。
// FIND SHORTEST PATH FROM "player102" TO "team204" OVER * YIELD path AS p; # 查找并返回 player102 到 team204 的最短路径。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Over(edges ...model.IEdge) (tx *DB) {
	tx = db.getInstance()
	names := make([]string, len(edges))
	for i, edge := range edges {
		names[i] = edge.EdgeName()
	}

	sql := strings.Join(names, ",")
	if sql == "" {
		tx.sql += fmt.Sprintf(" over * ")
		return
	}
	tx.sql += fmt.Sprintf(" over %s ", sql)
	return
}

// Reversely 设置反向遍历边。
// 配合FROM <vertex_list> OVER <edge_type_list> [{REVERSELY | BIDIRECT}]使用
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Reversely() (tx *DB) {
	tx = db.getInstance()
	tx.sql += " " + constants.DirectionReversely + " "
	return
}

// Bidirect 设置双向遍历边。
// 配合FROM <vertex_list> OVER <edge_type_list> [{REVERSELY | BIDIRECT}]使用
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Bidirect() (tx *DB) {
	tx = db.getInstance()
	tx.sql += " " + constants.DirectionBidirect + " "
	return
}

// Where 添加查询条件。
// MATCH (v:player)-[e]->(v2)  WHERE v2.player.age < 25  RETURN v2.player.name; # 过滤点属性
// GO FROM "player100" OVER follow WHERE follow.degree > 90 YIELD dst(edge); # 过滤边属性
// https://docs.nebula-graph.com.cn/3.6.0/3.ngql-guide/8.clauses-and-options/where/
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Where(sql string) (tx *DB) {
	tx = db.getInstance()
	tx.sql += fmt.Sprintf(" where %s ", sql)
	return
}

// Yield 指定查询结果需要返回的列或表达式。在 nGQL 中，YIELD和 openCypher 中的RETURN类似。
// (nebula自己定义的, 不是通用的openCypher模式, 所以未实现, 给出官方文档参考)
// GO FROM "player100" OVER follow YIELD properties($$).name AS Friend;
// FETCH PROP ON player "player100" YIELD properties(vertex).name;
// https://docs.nebula-graph.com.cn/3.6.0/3.ngql-guide/8.clauses-and-options/yield/
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Yield(sql string) (tx *DB) {
	tx = db.getInstance()
	tx.sql += sql
	return
}

// Group 对查询结果进行分组。
// (nebula自己定义的, 不是通用的openCypher模式, 所以未实现, 给出官方文档参考)
// GO FROM "player100" OVER follow BIDIRECT YIELD properties($$).name as Name | GROUP BY $-.Name YIELD $-.Name as Player, count(*) AS Name_Count;
// https://docs.nebula-graph.com.cn/3.6.0/3.ngql-guide/8.clauses-and-options/group-by/
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Group(sql string) (tx *DB) {
	tx = db.getInstance()
	tx.sql += sql
	return
}

// Limit 限制查询结果的数量, LIMIT在原生 nGQL 语句和 openCypher 兼容语句中的用法有所不同。
// (兼容性太低, 所以未实现, 给出官方文档参考)
// LOOKUP ON player YIELD id(vertex)| LIMIT 3; # 从结果中返回最前面的 3 行数据。
// GO FROM "player100" OVER follow REVERSELY YIELD properties($$).name AS Friend LIMIT 1, 3; # 从结果中返回第 2 行开始的 3 行数据。
// https://docs.nebula-graph.com.cn/3.6.0/3.ngql-guide/8.clauses-and-options/limit/
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Limit(sql string) (tx *DB) {
	tx = db.getInstance()
	tx.sql += sql
	return
}
