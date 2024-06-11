package converts

import (
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings"       // 字符串操作包
	"text/template" // 模板处理包，用于生成文本输出
)

// 定义了一个结构体，用于模板渲染删除 点SQL 所需的参数
//
// @Author: 罗德
// @Date: 2024/6/6
type deleteVertexStruct struct {
	Vid string // 顶点ID及其策略
}

// 初始化一个模板，用于生成删除点的SQL语句
//
// @Author: 罗德
// @Date: 2024/6/6
var deleteVertexTemplate = template.Must(template.New("delete_vertex").
	Parse("delete vertex {{.Vid}} with edge"))

// ConvertToDeleteVertexSql 是一个通用函数，将输入的结构转换为删除点的delete vertex sql语句
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
//
// 返回:
// string: 成功生成的delete vertex sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/6
func ConvertToDeleteVertexSql(vertex model.IVertex) (string, error) {
	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := deleteVertexTemplate.Execute(buf, &deleteVertexStruct{
		Vid: utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy()),
	})
	return buf.String(), err
}
