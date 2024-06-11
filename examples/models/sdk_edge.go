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
