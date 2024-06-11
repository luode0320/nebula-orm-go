package converts

import (
	"fmt"
	"nebula-orm-go/model"
	"nebula-orm-go/utils"
	"strings"
	"sync"
	"text/template"
)

// 定义了一个结构体，用于构建创建顶点SQL语句所需的各个部分
//
// @Author: 罗德
// @Date: 2024/6/7
type insertVertexBatchStruct struct {
	Name    string         // 顶点的名称
	Keys    string         // 属性列名序列，由各顶点公共属性名组成，以逗号分隔
	Vertexs []batchVertexs // 一批顶点的ID与属性值对集合
}

// 代表单个顶点在批量插入中的数据结构，包含顶点ID及其对应的属性值序列
//
// @Author: 罗德
// @Date: 2024/6/7
type batchVertexs struct {
	Vid    string // 顶点的唯一标识符（ID）
	Values string // 该顶点的属性值序列，按照与属性列名相同的顺序排列并以逗号分隔
}

// 创建并初始化一个模板，用于生成创建顶点的SQL语句。
//
// @Author: 罗德
// @Date: 2024/6/7
var insertVertexBatchTemplate = template.Must(template.New("insert_batch_vertex").
	Parse("insert vertex {{.Name}}({{.Keys}}) values {{range $i, $value := .Vertexs}}{{if $i}}, {{end}}{{$value.Vid}}:({{$value.Values}}){{end}}"))

// ConvertToInsertVertexBatchSql 将给定的顶点切片实体模型转换为Nebula图数据库的insert vertex sql语句。
// 存在则默认覆盖结构体所有属性, 参数 vertexs 没有被明确赋值的属性也会被覆盖为零值
//
// 参数:
// vertexs ([]model.IVertex): 点实体切片的接口，需实现IVertex接口。
//
// 返回:
// string: 成功生成的insert vertex sql语句。
// error: 如果转换过程中发生错误，则返回具体的错误信息。
//
// @Author: 罗德
// @Date: 2024/6/7
func ConvertToInsertVertexBatchSql(vertexs []model.IVertex) (string, error) {
	if len(vertexs) == 0 {
		return "", fmt.Errorf("参数为空")
	}

	// 获取点名称, 点属性名称
	tagName := vertexs[0].TagName()
	fields, _ := utils.GetNebulaTag(vertexs[0])

	var wg sync.WaitGroup
	var mu sync.Mutex
	var batch []batchVertexs
	for _, vertex := range vertexs {
		wg.Add(1)
		go func(vertex model.IVertex) {
			defer wg.Done()

			// 构建顶点id, 点属性值
			_, values := utils.GetNebulaTag(vertex)
			batchValue := batchVertexs{
				Vid:    utils.GetVidWithPolicy(vertex.GetVid(), vertex.GetPolicy()),
				Values: strings.Join(values, ","),
			}

			mu.Lock()
			batch = append(batch, batchValue)
			mu.Unlock()
		}(vertex)
	}
	wg.Wait()

	// 使用模板生成最终的SQL语句
	buf := new(strings.Builder)
	err := insertVertexBatchTemplate.Execute(buf, &insertVertexBatchStruct{
		Name:    tagName,
		Keys:    strings.Join(fields, ","),
		Vertexs: batch,
	})
	return buf.String(), err
}
