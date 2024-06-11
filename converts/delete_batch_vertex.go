package converts

import (
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings" // 字符串操作包
	"sync"
	"text/template" // 模板处理包，用于生成文本输出
)

// 定义了一个结构体，用于模板渲染批量删除 点SQL 所需的参数
//
// @Author: 罗德
// @Date: 2024/6/6
type deleteVertexBatchStruct struct {
	Vids string // 顶点ID, 以逗号分隔
}

// 初始化一个模板，用于生成批量删除边的SQL语句
//
// @Author: 罗德
// @Date: 2024/6/6
var deleteVertexBatchTemplate = template.Must(template.New("delete_batch_vertex").
	Parse("delete vertex {{.Vids}} with edge"))

// ConvertToDeleteVertexBatchSql 是一个通用函数，将输入的结构转换为批量删除点的delete vertex sql语句
//
// 参数:
// vertex ([]model.IVertex): 点实体切片的接口，需实现IVertex接口
//
// 返回:
// string: 成功生成的delete vertex sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/6
func ConvertToDeleteVertexBatchSql(vertexs []model.IVertex) (string, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var batch []string
	for _, vertex := range vertexs {
		wg.Add(1)
		go func(vertex model.IVertex) {
			defer wg.Done()

			mu.Lock()
			batch = append(batch, utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy()))
			mu.Unlock()
		}(vertex)
	}
	wg.Wait()

	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := deleteVertexBatchTemplate.Execute(buf, &deleteVertexBatchStruct{
		Vids: strings.Join(batch, ","),
	})
	return buf.String(), err
}
