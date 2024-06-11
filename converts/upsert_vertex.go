package converts

import (
	"fmt"
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings"
	"text/template"
)

// 定义了一个用于构建upsert点SQL语句的结构体，
//
// @Author: 罗德
// @Date: 2024/6/7
type upsertVertexStruct struct {
	Name  string // 顶点标签名称
	Vid   string // 顶点ID，可能包含ID策略
	Set   string // 用于set更新字段的属性,一次只能修改一个字段
	Where string // 用于补充过滤条件
	Yield string // return返回字段
}

// 初始化一个SQL模板，用于构造upsert vertex语句。
//
// @Author: 罗德
// @Date: 2024/6/7
var upsertVertexTemplate = template.Must(template.New("upsert_vertex").
	Parse("upsert vertex on {{.Name}} {{.Vid}} set {{.Set}} {{if .Where}}when {{.Where}}{{end}} yield {{.Yield}}"))

// ConvertToUpsertVertexSql 是一个通用函数，接受不同类型的输入并将其转换为用于upsert vertex操作的SQL语句。
// 如果顶点存在则更新, 不存在则插入, 更新性能低于update。
//
// 参数:
// vertex (model.IVertex): 点实体的接口，需实现IVertex接口
// set (string): 用于set更新字段的属性,一次只能修改一个字段(例: parent_key = 'O(∩_∩)O')
// where (string): 用于补充过滤条件(例: parent_key == '无')
//
// 返回:
// string: 成功生成的upsert vertex sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/7
func ConvertToUpsertVertexSql(vertex model.IVertex, set string, where string) (string, error) {
	if set == "" {
		return "", fmt.Errorf("用于set更新字段的属性不能为空")
	}
	// 获取return返回字段
	clause, err := utils.GetClause(vertex)
	if err != nil {
		return "", err
	}

	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err = upsertVertexTemplate.Execute(buf, &upsertVertexStruct{
		Name:  vertex.TagName(),
		Vid:   utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy()),
		Set:   set,
		Where: where,
		Yield: clause,
	})

	return buf.String(), err
}
