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

// Reversely 设置反向遍历边。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Reversely() (tx *DB) {
	tx = db.getInstance()
	tx.sql += " " + constants.DirectionReversely + " "
	return
}

// Bidirect 设置双向遍历边。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Bidirect() (tx *DB) {
	tx = db.getInstance()
	tx.sql += " " + constants.DirectionBidirect + " "
	return
}

// Go 设置查询步数，允许从起始顶点出发经过指定步数的边进行搜索。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Go(step int) (tx *DB) {
	tx = db.getInstance()

	if step > 1 {
		tx.sql += fmt.Sprintf(" go %d step ", step)
	}

	return
}

// From 指定查询的起始顶点，接受一系列实现了IVertex接口的实例。
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
	tx.sql += fmt.Sprintf(" over %s ", sql)
	return
}

// Limit 限制查询结果的数量。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Limit(limit int) (tx *DB) {
	tx = db.getInstance()
	tx.sql += fmt.Sprintf(" limit %d ", limit)
	return
}

// Where 添加查询条件。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Where(sql string) (tx *DB) {
	tx = db.getInstance()
	tx.sql += fmt.Sprintf(" where %s ", sql)
	return
}

// Yield 指定查询结果需要返回的列或表达式。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Yield(sql string) (tx *DB) {
	tx = db.getInstance()
	tx.sql += fmt.Sprintf(" yield %s ", sql)
	return
}

// Group 对查询结果进行分组。
//
// @Author: 罗德
// @Date: 2024/5/27
func (db *DB) Group(fields ...string) (tx *DB) {
	tx = db.getInstance()
	for i := range fields {
		fields[i] = "$-." + fields[i]
	}
	sql := strings.Join(fields, ",")
	tx.sql += fmt.Sprintf(" group by %s ", sql)
	return
}
