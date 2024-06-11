package converts

import (
	"fmt"
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings" // 字符串操作包
	"sync"
	"text/template" // 模板处理包，用于生成文本输出
)

// 定义了一个结构体，用于模板渲染批量删除 边SQL 所需的参数
//
// @Author: 罗德
// @Date: 2024/6/7
type deleteEdgeBatchStruct struct {
	Name  string             // 边的名称
	Edges []deleteBatchEdges // 源顶点和目标顶点
}

// 定义了一个结构体，用于模板渲染批量删除 边SQL 所需的参数。
//
// @Author: 罗德
// @Date: 2024/6/7
type deleteBatchEdges struct {
	Src string // 源顶点
	Dst string // 目标顶点
}

// 初始化一个模板，用于生成批量删除边的SQL语句
var deleteEdgeBatchTemplate = template.Must(template.New("delete_batch_edge").
	Parse("delete edge {{.Name}} {{range $i, $edge := .Edges}}{{if $i}}, {{end}}{{$edge.Src}} -> {{$edge.Dst}}{{end}}"))

// ConvertToDeleteEdgeBatchSql 是一个通用函数，将输入的结构转换为批量删除边的SQL语句
//
// 参数:
// edges ([]model.IEdge): 边实体切片的接口，需实现IEdge接口。
//
// 返回:
// string: 成功生成的delete edge sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/7
func ConvertToDeleteEdgeBatchSql(edges []model.IEdge) (string, error) {
	if len(edges) == 0 {
		return "", fmt.Errorf("参数为空")
	}

	// 获取边名称
	edgeName := edges[0].EdgeName()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var batch []deleteBatchEdges
	for _, edge := range edges {
		wg.Add(1)
		go func(edge model.IEdge) {
			defer wg.Done()

			// 构建源顶点, 目标顶点
			batchValue := deleteBatchEdges{
				Src: utils.GetVidWithPolicy(edge.GetVidSrc(), edge.GetVidSrcPolicy()),
				Dst: utils.GetVidWithPolicy(edge.GetVidDst(), edge.GetVidDstPolicy()),
			}

			mu.Lock()
			batch = append(batch, batchValue)
			mu.Unlock()
		}(edge)
	}
	wg.Wait()

	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := deleteEdgeBatchTemplate.Execute(buf, &deleteEdgeBatchStruct{
		Name:  edgeName,
		Edges: batch,
	})
	return buf.String(), err
}
