package model

import "nebula-orm-go/constants"

// IEdge 接口定义了边模型所需的方法，允许获取边的名称以及源顶点和目标顶点的ID及其策略。
//
// @Author: 罗德
// @Date: 2024/5/24
type IEdge interface {
	// EdgeName 返回边的名称。
	EdgeName() string
	// GetVidSrc 返回源顶点ID。
	GetVidSrc() interface{}
	// GetVidSrcPolicy 返回源顶点ID的策略。
	GetVidSrcPolicy() constants.Policy
	// GetVidDst 返回目标顶点ID。
	GetVidDst() interface{}
	// GetVidDstPolicy 返回目标顶点ID的策略。
	GetVidDstPolicy() constants.Policy
}

// EModel 结构体代表一个边模型，封装了源顶点和目标顶点的ID及其对应的ID策略。
// 注意：`Src`, `SrcPolicy`, `Dst`, `DstPolicy` 字段标记为`nebula:"-"`意味着这些字段不由ORM自动处理。
//
// @Author: 罗德
// @Date: 2024/5/24
type EModel struct {
	Src       interface{}      `nebula:"-"`
	SrcPolicy constants.Policy `nebula:"-"`
	Dst       interface{}      `nebula:"-"`
	DstPolicy constants.Policy `nebula:"-"`
}

// 实现IEdge接口，确保EModel类型满足IEdge定义的所有要求。
var _ IEdge = new(EModel)

// EdgeName 方法当前实现为抛出panic，实际应用中应返回具体的边名称。
func (v EModel) EdgeName() string {
	panic("请在此实现具体的边名称返回逻辑")
}

// GetVidSrc 实现接口方法，返回源顶点ID。
func (v EModel) GetVidSrc() interface{} {
	return v.Src
}

// GetVidSrcPolicy 实现接口方法，返回源顶点ID的策略。
func (v EModel) GetVidSrcPolicy() constants.Policy {
	return v.SrcPolicy
}

// GetVidDst 实现接口方法，返回目标顶点ID。
func (v EModel) GetVidDst() interface{} {
	return v.Dst
}

// GetVidDstPolicy 实现接口方法，返回目标顶点ID的策略。
func (v EModel) GetVidDstPolicy() constants.Policy {
	return v.DstPolicy
}
