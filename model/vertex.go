package model

import "nebula-orm-go/constants"

// ITag 接口定义了获取标签名称的方法，适用于任何需要标识名称的模型。
//
// @Author: 罗德
// @Date: 2024/5/24
type ITag interface {
	// TagName 返回标签的名称。
	TagName() string
}

// IVertex 接口继承了ITag接口，并扩展了顶点模型所需的方法，
// 允许获取顶点ID及其关联的ID生成策略。任何实现IVertex接口的类型都应当是一个完整的顶点定义。
//
// @Author: 罗德
// @Date: 2024/5/24
type IVertex interface {
	ITag
	// GetVid 返回顶点ID。
	GetVid() interface{}
	// GetPolicy 返回顶点ID的生成策略。
	GetPolicy() constants.Policy
}

// VModel 结构体代表一个顶点模型，封装了顶点ID及其ID策略。
// 字段`Vid`和`Policy`被标记为`nebula:"-"`，指示ORM框架不应自动处理这些字段。
//
// @Author: 罗德
// @Date: 2024/5/24
type VModel struct {
	Vid    interface{}      `nebula:"-"`
	Policy constants.Policy `nebula:"-"`
}

// 实现IVertex接口，确保VModel类型符合顶点模型的要求。
var _ IVertex = new(VModel)

// TagName 方法当前实现为抛出panic，实际场景下需要替换为返回该模型对应的标签名称，
// 建议采用蛇形命名法（snake_case）作为返回值。
func (v VModel) TagName() string {
	// TODO: 实现此处逻辑以返回模型名称，转换为蛇形风格
	panic("请在此实现具体的标签名称返回逻辑")
}

// GetVid 实现接口方法，返回顶点ID。
func (v VModel) GetVid() interface{} {
	return v.Vid
}

// GetPolicy 实现接口方法，返回顶点ID的策略。
func (v VModel) GetPolicy() constants.Policy {
	return v.Policy
}
