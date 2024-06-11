package converts

import (
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings"       // 字符串操作包
	"text/template" // 模板处理包，用于生成文本输出
)

// 定义了一个结构体，用于模板渲染插入 边SQL 所需的参数
//
// @Author: 罗德
// @Date: 2024/6/7
type insertEdgeStruct struct {
	Name         string // 边的名称
	Src, Dst     string // 源顶点和目标顶点
	Keys, Values string // 属性键和对应的值列表，格式化后的字符串
}

// 初始化一个模板，用于生成插入边的SQL语句
//
// @Author: 罗德
// @Date: 2024/6/7
var insertEdgeTemplate = template.Must(template.New("insert_edge").
	Parse("insert edge {{.Name}}({{.Keys}}) values {{.Src}} -> {{.Dst}}:({{.Values}})"))

// ConvertToInsertEdgeSql 是一个通用函数，将输入的结构转换为创建边的insert edge sql语语句
//
// 参数:
// edge (model.IEdge): 边实体的接口，需实现IEdge接口。
//
// 返回:
// string: 成功生成的insert edge sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/7
func ConvertToInsertEdgeSql(edge model.IEdge) (string, error) {
	// 获取边属性名称, 边属性值
	fields, values := utils.GetNebulaTag(edge)

	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := insertEdgeTemplate.Execute(buf, &insertEdgeStruct{
		Name:   edge.EdgeName(),
		Src:    utils.GetVidWithPolicy(edge.GetVidSrc(), edge.GetVidSrcPolicy()),
		Dst:    utils.GetVidWithPolicy(edge.GetVidDst(), edge.GetVidDstPolicy()),
		Keys:   strings.Join(fields, ","),
		Values: strings.Join(values, ","),
	})
	return buf.String(), err
}

// 初始化一个模板，用于生成插入忽略已存在边的SQL语句
//
// @Author: 罗德
// @Date: 2024/6/7
var insertEdgeIgnoreTemplate = template.Must(template.New("insert_edge_ignore").
	Parse("insert edge if not exists {{.Name}}({{.Keys}}) values {{.Src}} -> {{.Dst}}:({{.Values}})"))

// ConvertToInsertEdgeIgnoreSql 是一个通用函数，将输入的结构转换为创建边的insert edge if not exists sql语语句
// 检测待插入的边是否存在，只有不存在时，才会插入
// 1. if not exists 仅检测<边的类型、起始点、目的点和 rank>是否存在，不会检测属性值是否重合
// 2. if not exists 会先读取一次数据是否存在，因此对性能会有明显影响
// 3. if not exists 不支持批量插入
//
// 参数:
// edge (model.IEdge): 边实体的接口，需实现IEdge接口。
//
// 返回:
// string: 成功生成的insert edge if not exists sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/7
func ConvertToInsertEdgeIgnoreSql(edge model.IEdge) (string, error) {
	// 获取边属性名称, 边属性值
	fields, values := utils.GetNebulaTag(edge)

	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := insertEdgeIgnoreTemplate.Execute(buf, &insertEdgeStruct{
		Name:   edge.EdgeName(),
		Src:    utils.GetVidWithPolicy(edge.GetVidSrc(), edge.GetVidSrcPolicy()),
		Dst:    utils.GetVidWithPolicy(edge.GetVidDst(), edge.GetVidDstPolicy()),
		Keys:   strings.Join(fields, ","),
		Values: strings.Join(values, ","),
	})
	return buf.String(), err
}
