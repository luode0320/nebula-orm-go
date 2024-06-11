package converts

import (
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings"       // 字符串操作包
	"text/template" // 模板处理包，用于生成文本输出
)

// 定义了一个结构体，用于模板渲染删除 边SQL 所需的参数
//
// @Author: 罗德
// @Date: 2024/6/6
type deleteEdgeStruct struct {
	Name string // 边的名称
	Src  string // 源顶点
	Dst  string // 目标顶点
}

// 初始化一个模板，用于生成删除边的SQL语句
//
// @Author: 罗德
// @Date: 2024/6/6
var deleteEdgeTemplate = template.Must(template.New("delete_edge").
	Parse("delete edge {{.Name}} {{.Src}} -> {{.Dst}}"))

// ConvertToDeleteEdgeSql 是一个通用函数，将输入的结构转换为删除边的delete edge sql语句
//
// 参数:
// edges (model.IEdge): 边实体的接口，需实现IEdge接口。
//
// 返回:
// string: 成功生成的delete edge sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/6
func ConvertToDeleteEdgeSql(edge model.IEdge) (string, error) {
	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := deleteEdgeTemplate.Execute(buf, &deleteEdgeStruct{
		Name: edge.EdgeName(),
		Src:  utils.GetVidWithPolicy(edge.GetVidSrc(), edge.GetVidSrcPolicy()),
		Dst:  utils.GetVidWithPolicy(edge.GetVidDst(), edge.GetVidDstPolicy()),
	})
	return buf.String(), err
}
