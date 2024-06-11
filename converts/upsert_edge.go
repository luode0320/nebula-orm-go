package converts

import (
	"fmt"
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings"
	"text/template"
)

// 定义一个结构体用于存储构建upsert边SQL所需的属性
//
// @Author: 罗德
// @Date: 2024/6/7
type upsertEdgeStruct struct {
	Name     string // 边的名称
	Src, Dst string // 源顶点和目标顶点的标识
	Set      string // 用于set更新字段的属性,一次只能修改一个字段
	Where    string // 用于补充过滤条件
	Yield    string // return返回字段
}

// 初始化一个模板，用于生成upsert边的SQL语句
//
// @Author: 罗德
// @Date: 2024/6/7
var upsertEdgeTemplate = template.Must(template.New("upsert_edge").
	Parse("upsert edge on {{.Name}} {{.Src}} -> {{.Dst}} set {{.Set}} {{if .Where}}when {{.Where}}{{end}} yield {{.Yield}}"))

// ConvertToUpsertEdgeSql 接收任意类型输入，将其转换为用于upsert边操作的SQL语句。
// 如果边存在则更新, 不存在则插入, 更新性能低于update
//
// 参数:
// edge (model.IEdge): 边实体的接口，需实现IEdge接口
// set (string): 用于set更新字段的属性,一次只能修改一个字段(例: test = 'O(∩_∩)O')
// where (string): 用于补充过滤条件(例: test == '无')
//
// 返回:
// string: 成功生成的upsert vertex sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/7
func ConvertToUpsertEdgeSql(edge model.IEdge, set string, where string) (string, error) {
	if set == "" {
		return "", fmt.Errorf("用于set更新字段的属性不能为空")
	}
	// 获取return返回字段
	clause, err := utils.GetClauseByNorm(edge)
	if err != nil {
		return "", err
	}

	// 使用缓冲区高效构建SQL字符串
	buf := new(strings.Builder)
	// 执行模板，填充upsertEdgeStruct结构体
	err = upsertEdgeTemplate.Execute(buf, &upsertEdgeStruct{
		Name:  edge.EdgeName(),
		Src:   utils.GetVidWithPolicy(edge.GetVidSrc(), edge.GetVidSrcPolicy()),
		Dst:   utils.GetVidWithPolicy(edge.GetVidDst(), edge.GetVidDstPolicy()),
		Set:   set,
		Where: where,
		Yield: clause,
	})
	// 返回构建好的SQL语句
	return buf.String(), err
}
