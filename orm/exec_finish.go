// Package orm 单次完整调用: db.Execute()
//
// @Author: 罗德
// @Date: 2024/5/27
package orm

import (
	"nebula-orm-go/dialectors"
)

// Execute 执行给定的SQL语句，并返回执行结果集。
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) Execute(sql string) (*dialectors.ResultSet, error) {
	return db.execute(sql)
}

// ExecuteAndParse 执行SQL语句，并将结果解析到给定的结构体中。输入可以是单个map、结构体指针、结构体切片指针。
// in 可以是 map[string]interface{}, *Strcut, *[]map, *[]struct
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) ExecuteAndParse(sql string, in interface{}) error {
	nResult, err := db.execute(sql)
	if err != nil {
		return err
	}
	return nResult.UnmarshalResultSet(in)
}

// ReturnRow 直接执行当前构建的SQL语句并返回未经处理的结果集。
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) ReturnRow() (*dialectors.ResultSet, error) {
	return db.execute(db.sql)
}

// Return 执行当前构建的SQL语句，并将结果反序列化到指定的输出结构体中。
//
// @Author: 罗德
// @Date: 2024/5/28
func (db *DB) Return(out interface{}) error {
	return db.ExecuteAndParse(db.sql, out)
}
