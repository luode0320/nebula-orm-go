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
