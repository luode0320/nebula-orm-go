package converts

import (
	"fmt"
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings"
	"sync"
	"text/template" // 模板处理包，用于生成文本输出
)

// 定义了一个结构体，用于模板渲染批量创建 边SQL 所需的参数
//
// @Author: 罗德
// @Date: 2024/6/7
type insertEdgeBatchStruct struct {
	Name  string       // 边的名称
	Keys  string       // 边的类型, 目前只能是单个, 可以对应多个value
	Edges []batchEdges // 一批边起点终点的集合
}

// 定义了一个结构体，用于模板渲染批量创建 边SQL 所需的value参数。
//
// @Author: 罗德
// @Date: 2024/6/7
type batchEdges struct {
	Src    string // 源顶点
	Dst    string // 目标顶点
	Values string // 边类型的值
}

// 初始化一个模板，用于生成插入边的SQL语句
//
// @Author: 罗德
// @Date: 2024/6/7
var insertEdgeBatchTemplate = template.Must(template.New("insert_batch_edge").
	Parse("insert edge {{.Name}}({{.Keys}}) values {{range $i, $edge := .Edges}}{{if $i}}, {{end}}{{$edge.Src}} -> {{$edge.Dst}}:({{$edge.Values}}){{end}}"))

// ConvertToInsertEdgeBatchSql 是一个通用函数，将输入的结构转换为创建边的SQL语句
//
// 参数:
// edges ([]model.IEdge): 边实体切片的接口，需实现IEdge接口。
//
// 返回:
// string: 成功生成的create edge sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/7
func ConvertToInsertEdgeBatchSql(edges []model.IEdge) (string, error) {
	if len(edges) == 0 {
		return "", fmt.Errorf("参数为空")
	}

	// 获取边名称, 边属性名称
	tagName := edges[0].EdgeName()
	fields, _ := utils.GetNebulaTag(edges[0])

	var wg sync.WaitGroup
	var mu sync.Mutex
	var batch []batchEdges
	for _, edge := range edges {
		wg.Add(1)
		go func(edge model.IEdge) {
			defer wg.Done()

			// 构建源顶点, 目标顶点, 边属性值
			_, values := utils.GetNebulaTag(edge)
			batchValue := batchEdges{
				Src:    utils.GetVidWithPolicy(edge.GetVidSrc(), edge.GetVidSrcPolicy()),
				Dst:    utils.GetVidWithPolicy(edge.GetVidDst(), edge.GetVidDstPolicy()),
				Values: strings.Join(values, ","),
			}

			mu.Lock()
			batch = append(batch, batchValue)
			mu.Unlock()
		}(edge)
	}
	wg.Wait()

	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := insertEdgeBatchTemplate.Execute(buf, &insertEdgeBatchStruct{
		Name:  tagName,
		Keys:  strings.Join(fields, ","),
		Edges: batch,
	})
	return buf.String(), err
}
