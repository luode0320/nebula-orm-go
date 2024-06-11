package utils

import (
	"fmt"
	nebula_type "github.com/vesoft-inc/nebula-go/v3/nebula"
	"nebula-orm-go/constants"
	"reflect"
	"strings"
	"time"
)

// GetStructFieldTagMap 将 struct 中标记为 nebula 的 tag 标签的所有字段提取出来, 并记录 field 从标签名到字段索引位置的映射关系
// 这个映射关系有助于后续根据标签快速定位到结构体字段，进行数据填充或其它操作。
//
// 参数:
// - typ (reflect.Type): 要处理的结构体类型反射类型。
//
// 返回:
// - map[string]int: 一个字典，键为结构体字段上标记的标签名（使用 constants.StructTagName 定义的标签），值为该字段在结构体中的索引位置。
//
// @Author: 罗德
// @Date: 2024/5/27
func GetStructFieldTagMap(typ reflect.Type) map[string]int {
	// 初始化一个空映射以存储标签名到字段索引的对应关系
	tagMap := make(map[string]int)
	// 遍历结构体的所有字段
	for i := 0; i < typ.NumField(); i++ {
		// 获取当前字段的标记值，使用 constants.StructTagName 指定的标签名称
		tag := typ.Field(i).Tag.Get(constants.StructTagName)
		// 如果标签不存在或被设置为 "-", 表示该字段不应被处理，跳过
		if tag == "" || tag == "-" {
			continue
		}
		// 将标签名与字段索引存入映射中
		tagMap[tag] = i
	}

	// 返回构建好的映射
	return tagMap
}

// SetFieldValue 将 Nebula Graph 数据库中的值（nValue）转换并设置到 Go 语言结构体 struct.field 字段中，自动处理不同数据类型的适配。
// 这个函数根据字段类型选择合适的方法将 `nValue` 转换成相应类型并赋值给结构体字段。
//
// 参数:
// - tag (string): 字段的标签名，主要用于日志或调试信息。
// - field (reflect.Value): 结构体字段的反射值，数据将被设置到这里。
// - nValue (*nebula_type.Value): Nebula Graph 数据库中的值对象，需要转换为 Go 语言的数据类型。
//
// 返回:
// - error: 如果不支持的类型转换发生，则返回错误；若成功则返回nil。
//
// @Author: 罗德
// @Date: 2024/5/27
func SetFieldValue(field reflect.Value, nValue *nebula_type.Value) error {
	switch field.Kind() {
	case reflect.Bool: // 布尔类型
		field.SetBool(nValue.GetBVal())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: // 整型
		field.SetInt(nValue.GetIVal())

	case reflect.Float32, reflect.Float64: // 浮点型
		field.SetFloat(nValue.GetFVal())

	case reflect.String: // 字符串类型
		field.SetString(string(nValue.GetSVal()))

	case reflect.Struct: // 结构体型，目前仅处理 time.Time 类型
		switch field.Type().String() {
		case "time.Time":
			ts := nValue.GetIVal()                       // 获取整数值
			field.Set(reflect.ValueOf(time.Unix(ts, 0))) // 转换为 Time 类型并设置
		default:
			// 对于其他结构体类型，当前逻辑未实现，可根据需要扩展
		}
	default: // 未处理的类型
		// 这里可以选择打印日志或做更详细的错误处理，当前逻辑直接返回nil，表示不处理
		return nil
	}
	return nil
}

// NValueToInterface 将 nebula_type.Value 类型的值转换为 Go 语言的 interface{} 类型。
// nebula_type.Value 可能包含不同类型的数据，此函数根据 Value 内部实际设置的值类型，
// 返回相应的基本数据类型值或者复杂的结构体，以提高代码的灵活性和兼容性。
//
// 参数:
// - p (*nebula_type.Value): 来自 Nebula Graph 数据库的值对象，可能封装了不同类型的数据。
//
// 返回:
// - interface{}: 与 Nebula Value 内部类型相对应的 Go 语言值。如果 Value 未设置任何类型，则返回 nil。
//
// @Author: 罗德
// @Date: 2024/5/27
func NValueToInterface(p *nebula_type.Value) interface{} {
	// 检查并转换各种可能的数据类型
	if p.IsSetNVal() {
		return nil
	}
	if p.IsSetBVal() {
		return p.GetBVal() // 布尔值
	}
	if p.IsSetIVal() {
		return p.GetIVal() // 整数
	}
	if p.IsSetFVal() {
		return p.GetFVal() // 浮点数值
	}
	if p.IsSetSVal() {
		return string(p.GetSVal()) // 字符串
	}
	if p.IsSetDVal() {
		return p.GetDVal() // 日期
	}
	if p.IsSetTVal() {
		return p.GetTVal() // 时间
	}
	if p.IsSetDtVal() {
		return p.GetDtVal() // 日期时间值
	}
	if p.IsSetVVal() {
		return p.GetVVal() // 点值
	}
	if p.IsSetEVal() {
		return p.GetEVal() // 边值
	}
	if p.IsSetPVal() {
		return p.GetPVal() // 路径值
	}
	if p.IsSetLVal() {
		return p.GetLVal() // 列表值
	}
	if p.IsSetMVal() {
		return p.GetMVal() // 映射值
	}
	if p.IsSetUVal() {
		return p.GetUVal() // 设置值
	}
	if p.IsSetGVal() {
		return p.GetGVal() // 通用值
	}

	// 如果没有任何类型被设置，则返回nil
	return nil
}

// GetVClauseByNorm 从结构体中提取所有带有`nebula`标签的字段名，
// 并生成形如`v.vertex.key as key`的SQL片段。
//
// @Author: 罗德
// @Date: 2024/5/27
func GetVClauseByNorm(v any, vertex string) (string, error) {
	var parts []string
	val := reflect.ValueOf(v)
	typ := val.Type()

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("参数必须是一个结构体")
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(constants.StructTagName)
		if tag != "" {
			parts = append(parts, fmt.Sprintf("%s.%s.%s as %s", constants.V, vertex, tag, tag))
		}
	}

	return strings.Join(parts, ","), nil
}

// GetClauseByNorm 从结构体中提取所有带有`nebula`标签的字段名，
// 并生成形如`key as key`的SQL片段。
//
// @Author: 罗德
// @Date: 2024/5/27
func GetClauseByNorm(v any) (string, error) {
	var parts []string
	val := reflect.ValueOf(v)
	typ := val.Type()

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("参数必须是一个结构体")
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(constants.StructTagName)
		if tag != "" {
			parts = append(parts, fmt.Sprintf("%s as %s", tag, tag))
		}
	}

	return strings.Join(parts, ","), nil
}

// GetNebulaTag 从结构体中提取所有带有`nebula`标签的字段名，
// 并生成形如`v.vertex.key as key`的SQL片段。
//
// @Author: 罗德
// @Date: 2024/5/27
func GetNebulaTag(v any) ([]string, []string) {
	var fields []string
	var values []string
	val := reflect.ValueOf(v)
	typ := val.Type()

	if val.Kind() != reflect.Struct {
		return fields, values
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(constants.StructTagName)
		if tag != "" {
			// 获取字段值
			fieldValue := val.Field(i)
			// 这里简单地使用字符串形式表示字段值，实际情况可能需要更复杂的处理，特别是对于非基本类型
			valueStr := fmt.Sprintf("%s", GetVidWithPolicy(fieldValue.Interface(), 0))
			fields = append(fields, tag)
			values = append(values, valueStr)
		}
	}

	return fields, values
}

// GetVidWithPolicy 根据给定的顶点ID（vid）和ID策略（policy），生成符合Nebula图数据库要求的ID字符串表示形式。
// 这个函数首先会根据vid的实际类型进行格式化处理，然后根据策略调整最终的字符串形式。
//
// @Author: 罗德
// @Date: 2024/5/27
func GetVidWithPolicy(vid interface{}, policy constants.Policy) string {
	// 初始化一个空字符串用于存放处理后的vid值
	vidStr := ""

	// 使用类型断言确定vid的具体类型，并据此格式化vid
	switch vid.(type) {
	case int, int8, int32, int64, float32, float64:
		// 对于数字类型，直接转换为字符串
		vidStr = fmt.Sprint(vid)
	case string:
		// 对于字符串类型，添加单引号包裹
		vidStr = "'" + vid.(string) + "'"
	default:
		// 其他类型也用单引号包裹，确保兼容性
		vidStr += "'" + fmt.Sprint(vid) + "'"
	}

	// 根据策略调整vidStr
	switch policy {
	case constants.PolicyHash:
		// 如果策略为哈希，将vid用hash函数包裹
		vidStr = "hash(" + vidStr + ")"
	}

	// 返回处理后的vid字符串
	return vidStr
}
