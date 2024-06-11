package converts

import (
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings"
	"text/template"
)

// 定义了一个结构体，用于构建创建顶点SQL语句所需的各个部分
//
// @Author: 罗德
// @Date: 2024/6/7
type insertVertexStruct struct {
	Name   string // 顶点的名称
	Vid    string // 顶点ID及其策略
	Keys   string // 用于插入的属性键列表，格式化后的字符串
	Values string // 对应属性值列表，格式化后的字符串
}

// 创建并初始化一个模板，用于生成创建顶点的SQL语句。
//
// @Author: 罗德
// @Date: 2024/6/7
var insertVertexTemplate = template.Must(template.New("insert_vertex").
	Parse("insert vertex {{.Name}}({{.Keys}}) values {{.Vid}}:({{.Values}})"))

// ConvertToInsertVertexSql 将给定的顶点实体模型转换为Nebula图数据库的insert vertex sql语句。
// 存在则默认覆盖结构体所有属性, 参数 vertex 没有被明确赋值的属性也会被覆盖为零值
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口。
//
// 返回:
// string: 成功生成的insert vertex sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/7
func ConvertToInsertVertexSql(vertex model.IVertex) (string, error) {
	// 获取点属性名称, 点属性值
	fields, values := utils.GetNebulaTag(vertex)

	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := insertVertexTemplate.Execute(buf, &insertVertexStruct{
		Name:   vertex.TagName(),
		Vid:    utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy()),
		Keys:   strings.Join(fields, ","),
		Values: strings.Join(values, ","),
	})
	return buf.String(), err
}

// 创建并初始化一个模板，用于生成创建忽略已存在顶点的SQL语句。
//
// @Author: 罗德
// @Date: 2024/6/7
var insertVertexIgnoreTemplate = template.Must(template.New("insert_vertex_ignore").
	Parse("insert vertex if not exists {{.Name}}({{.Keys}}) values {{.Vid}}:({{.Values}})"))

// ConvertToInsertVertexIgnoreSql 将给定的顶点实体模型转换为Nebula图数据库的insert vertex if not exists sql语句。
// 检测待插入的 VID 是否存在，只有不存在时，才会插入，如果已经存在，不会进行修改
// 1. if not exists 仅检测 VID + Tag 的值是否相同，不会检测属性值
// 2. if not exists 会先读取一次数据是否存在，因此对性能会有明显影响
// 3. if not exists 不支持批量插入
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口。
//
// 返回:
// string: 成功生成的insert vertex if not exists sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/7
func ConvertToInsertVertexIgnoreSql(vertex model.IVertex) (string, error) {
	// 获取点属性名称, 点属性值
	fields, values := utils.GetNebulaTag(vertex)

	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := insertVertexIgnoreTemplate.Execute(buf, &insertVertexStruct{
		Name:   vertex.TagName(),
		Vid:    utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy()),
		Keys:   strings.Join(fields, ","),
		Values: strings.Join(values, ","),
	})
	return buf.String(), err
}
