package dialectors

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"nebula-orm-go/utils"
	"reflect"
)

// 定义两个错误常量，分别表示尝试向nil映射赋值的错误和记录未找到的错误
var (
	NilPointError       = errors.New("尝试向nil赋值")
	RecordNotFoundError = errors.New("未找到任何记录")
)

// ResultSet 结构体，封装了nebula.ResultSet，用于处理查询结果
//
// @Author: 罗德
// @Date: 2024/5/24
type ResultSet struct {
	*nebula.ResultSet // 内嵌nebula.ResultSet类型
}

// PrintResult 打印查询结果
//
// @Author: 罗德
// @Date: 2024/5/29
func (resultSet *ResultSet) PrintResult(log string) *ResultSet {
	mapResult := make([]map[string]interface{}, 0)
	_ = resultSet.UnmarshalResultSet(&mapResult)
	fmt.Printf("%s:\n", log)
	// 将 map 转换为 JSON 格式的字节切片
	jsonData, err := json.MarshalIndent(mapResult, "", "  ")
	if err != nil {
		fmt.Println("转换为 JSON 格式时出错:", err)
		return resultSet
	}
	// 将 JSON 格式的字节切片转换为字符串并打印
	fmt.Println(string(jsonData))
	fmt.Println()
	return resultSet
}

// ToMap 转换为map
//
// @Author: 罗德
// @Date: 2024/5/29
func (resultSet *ResultSet) ToMap() map[string]interface{} {
	mapResult := make(map[string]interface{})
	if err := resultSet.UnmarshalResultSet(&mapResult); err != nil {
		fmt.Println("转换结果集到map时发生错误:", err)
		return mapResult
	}

	return mapResult
}

// ToVertexs 解码 ResultSet 的点值到map
//
// @Author: 罗德
// @Date: 2024/5/29
func (resultSet *ResultSet) ToVertexs() ([][]map[string]interface{}, error) {
	// 初始化一个新的字典切片，大小与 ResultSet 的行数相同
	vertexMaps := make([][]map[string]interface{}, len(resultSet.GetRows()))

	// 检查 ResultSet 中是否有数据行，如果没有数据行（小于1行），说明记录未找到，返回相应错误
	if resultSet.GetRowSize() < 1 {
		fmt.Println("检查 ResultSet 中是否有数据行: 记录未找到")
		return vertexMaps, RecordNotFoundError
	}

	// 遍历 ResultSet 的每一行
	for i, row := range resultSet.GetRows() {
		// 为当前行创建一个新的字典
		vertexMaps[i] = make([]map[string]interface{}, 0)
		// 遍历列名，将当前行的每列值转换并存入字典
		for _, p := range row.Values {
			if !p.IsSetLVal() {
				continue
			}
			for _, vval := range p.GetLVal().GetValues() { // 列表值
				if !vval.IsSetVVal() {
					continue
				}
				for _, prop := range vval.GetVVal().GetTags() {
					propsMap := make(map[string]interface{})
					for key, value := range prop.GetProps() {
						// 将每一列的值转换为 interface{} 类型，并存入目标字典中，键为列名
						propsMap[key] = utils.NValueToInterface(value)
					}
					vertexMaps[i] = append(vertexMaps[i], propsMap) // 点值
				}
			}
		}
	}

	return vertexMaps, nil
}

// UnmarshalResultSet 解码 ResultSet 到给定的接口类型中。
// 支持的目标类型包括：字典(map)，字典指针(*map)，字典切片指针(*[]map)，结构体及其指针，结构体切片及其指针，以及一些基本整数类型。
//
// @Author: 罗德
// @Date: 2024/5/29
func (resultSet *ResultSet) UnmarshalResultSet(in interface{}) error {
	switch values := in.(type) {
	case *map[string]interface{}: // 解码到字典指针
		return toMap(*values, resultSet)

	case *[]map[string]interface{}: // 解码到字典切片指针
		return toMapSlice(values, resultSet)

	case *map[string]string: // 解码到字典字符串指针
		return toMapString(*values, resultSet)

	case *[]map[string]string: // 解码到字典字符串切片指针
		return toMapStringSlice(values, resultSet)

	default:
		// 获取反射值以便进一步检查类型
		val := reflect.ValueOf(values)
		switch val.Kind() {
		case reflect.Ptr:
			val = reflect.Indirect(val)
			switch val.Kind() {
			case reflect.Interface: // 解码到结构体
				val = val.Elem()
				return toStruct(val, resultSet)

			case reflect.Struct: // 解码到结构体
				return toStruct(val, resultSet)

			case reflect.Slice: // 解码到结构体切片
				return toStructSlice(val, resultSet)

			case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64: // 处理整数类型
				return toInt(val, resultSet)

			default: // 不支持的类型
				return errors.Errorf("不支持的类型. 类型是:%v", val.Kind())
			}
		default:
			return errors.New("必须是指针类型")
		}
	}
}

// 将 ResultSet 中的整型数据转换并设置给基本整数类型的反射值。
// 假设 ResultSet 只有一行一列，并且这一列的数据可以安全转化为整数类型。
//
// 参数:
// - val (reflect.Value): 整数类型的反射值，用于接收转换后的整数值。
// - resultSet (*dialectors.ResultSet): 查询结果集，应只包含单行单列的整型数据。
//
// @Author: 罗德
// @Date: 2024/5/29
func toInt(val reflect.Value, resultSet *ResultSet) (err error) {
	// 检查结构体是否为nil，如果是则返回预定义错误
	if val.Interface() == nil {
		return NilPointError
	}

	// 检查 ResultSet 中是否有数据行，如果没有数据行（小于1行），说明记录未找到，返回相应错误
	if resultSet.GetRowSize() < 1 {
		val.SetInt(0)
		return nil
	}

	// 使用defer-recover来捕获并处理潜在的运行时错误，以确保函数能安全返回错误而不是panic
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New("未知的exec错误")
			}
		}
	}()

	// 由于假定 ResultSet 只有一行一列，直接获取该行该列的整数值
	cnt := resultSet.GetRows()[0].GetValues()[0].IVal
	// 将获取的整数值设置给反射值
	val.SetInt(*cnt)

	return nil
}

// 将 ResultSet 中的数据转换并填充到结构体反射值中。
// 这个函数假设 ResultSet 只有一行数据，且结构体字段使用了相应的标签与 ResultSet 的列名对应。
//
// 参数:
// - val (reflect.Value): 结构体的反射值，数据将被填充到这里。
// - resultSet (*dialectors.ResultSet): 查询结果集，包含要转换的数据。
//
// @Author: 罗德
// @Date: 2024/5/29
func toStruct(val reflect.Value, resultSet *ResultSet) (err error) {
	// 检查结构体是否为nil，如果是则返回预定义错误
	if val.Interface() == nil {
		return NilPointError
	}

	// 检查 ResultSet 中是否有数据行，如果没有数据行（小于1行），说明记录未找到，返回相应错误
	if resultSet.GetRowSize() < 1 {
		return RecordNotFoundError
	}

	// 使用defer-recover来捕获并处理潜在的运行时错误，以确保函数能安全返回错误而不是panic
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New("未知的exec错误")
			}
		}
	}()

	// 获取 ResultSet 的第一行数据
	row := resultSet.GetRows()[0]
	// 将 struct 中标记为 nebula-orm-go 的 tag 标签的所有字段提取出来, 并记录 field 从标签名到字段索引位置的映射关系
	fieldTagMap := utils.GetStructFieldTagMap(val.Type())

	// 遍历列名，查找与之对应的结构体字段，并设置值
	for j, col := range resultSet.GetColNames() {
		// 查找该列名对应的结构体字段位置
		fieldPos, ok := fieldTagMap[col]
		if !ok {
			continue
		}

		// 获取列的值
		value := row.GetValues()[j]
		// 通过反射获取字段
		field := val.Field(fieldPos)
		// 尝试设置字段值，这里假设setFieldValue是一个处理字段赋值的函数
		err = utils.SetFieldValue(field, value)
		// 如果在设置字段值时发生错误，应立刻返回
		if err != nil {
			return err
		}
	}

	return nil
}

// 将 ResultSet 中的多行数据转换并填充到结构体切片的反射值中。
// 每个结构体实例对应 ResultSet 中的一行数据。
//
// 参数:
// - val (reflect.Value): 结构体切片的反射值，数据将被填充到这里。
// - resultSet (*dialectors.ResultSet): 查询结果集，包含要转换的数据。
//
// @Author: 罗德
// @Date: 2024/5/29
func toStructSlice(val reflect.Value, resultSet *ResultSet) (err error) {
	// 检查结构体切片是否为nil，如果是则返回预定义错误
	if val.Interface() == nil {
		return NilPointError
	}
	// 检查 ResultSet 中是否有数据行，如果没有数据行（小于1行），直接返回，不视为错误
	if resultSet.GetRowSize() < 1 {
		return nil
	}

	// 使用defer-recover来捕获并处理潜在的运行时错误，确保函数能安全返回错误而不是panic
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New("未知的exec错误")
			}
		}
	}()

	// 重置 val 为新的切片，长度和容量都为 ResultSet 的行数
	val.Set(reflect.MakeSlice(val.Type(), resultSet.GetRowSize(), resultSet.GetRowSize()))
	// TODO 获取切片中第一个元素的字段标签映射，这里似乎有误，应该是直接调用 GetStructFieldTagMap 函数并传入类型
	fieldTagMap := utils.GetStructFieldTagMap(val.Index(0).Type())
	// 遍历 ResultSet 的每一行数据
	for i, row := range resultSet.GetRows() {
		// 这里可以优化 GetColNames, 只循环两个共有的 key
		// 对于每行数据，遍历其列名
		for j, col := range resultSet.GetColNames() {
			// 查找列名对应的结构体字段位置
			fieldPos, ok := fieldTagMap[col]
			if !ok {
				continue
			}

			// 获取当前列的值
			nValue := row.GetValues()[j]
			// 获取切片中的结构体实例
			field := val.Index(i).Field(fieldPos)
			// 尝试设置字段值
			err = utils.SetFieldValue(field, nValue)
			// 如果在设置字段值时发生错误，应立刻返回
			if err != nil {
				return err
			}
		}
	}
	return
}

// 将 ResultSet 数据转换并填充到一个字典(map[string]interface{})中。
//
// 参数:
// - values: 用于存储解码后数据的目标字典。必须是非nil的，否则函数会返回错误。
// - resultSet: 包含查询结果的 ResultSet 实例。
//
// @Author: 罗德
// @Date: 2024/5/29
func toMap(values map[string]interface{}, resultSet *ResultSet) error {
	// 检查目标字典是否为nil，如果是则返回预定义错误
	if values == nil {
		return NilPointError
	}

	// 检查 ResultSet 中是否有数据行，如果没有数据行（小于1行），说明记录未找到，返回相应错误
	if resultSet.GetRowSize() < 1 {
		return RecordNotFoundError
	}

	// 获取 ResultSet 的第一行数据
	row := resultSet.GetRows()[0]
	// 遍历 ResultSet 的列名集合
	for i, col := range resultSet.GetColNames() {
		// 将每一列的值转换为 interface{} 类型，并存入目标字典中，键为列名
		values[col] = utils.NValueToInterface(row.Values[i])
	}

	// 成功处理完所有数据，返回nil表示没有错误
	return nil
}

// 将 ResultSet 数据转换并填充到一个字典(map[string]string)中。
// 如果resultSet的字段结果类型不是string, 将会赋值为""
//
// 参数:
// - values: 用于存储解码后数据的目标字典。必须是非nil的，否则函数会返回错误。
// - resultSet: 包含查询结果的 ResultSet 实例。
//
// @Author: 罗德
// @Date: 2024/5/29
func toMapString(values map[string]string, resultSet *ResultSet) error {
	// 检查目标字典是否为nil，如果是则返回预定义错误
	if values == nil {
		return NilPointError
	}

	// 检查 ResultSet 中是否有数据行，如果没有数据行（小于1行），说明记录未找到，返回相应错误
	if resultSet.GetRowSize() < 1 {
		return RecordNotFoundError
	}

	// 获取 ResultSet 的第一行数据
	row := resultSet.GetRows()[0]
	// 遍历 ResultSet 的列名集合
	for i, col := range resultSet.GetColNames() {
		// 将每一列的值转换为 string 类型，并存入目标字典中，键为列名
		valueToInterface := utils.NValueToInterface(row.Values[i])
		if str, ok := valueToInterface.(string); ok {
			values[col] = str
		}
	}

	// 成功处理完所有数据，返回nil表示没有错误
	return nil
}

// 将 ResultSet 数据转换并填充到一个字典切片(*[]map[string]interface{})中。
// 每个字典代表 ResultSet 中的一行数据。
//
// 参数:
// - values: 指向字典切片的指针，用于存储转换后的 ResultSet 数据。
// - resultSet: 包含查询结果的 ResultSet 实例。
//
// @Author: 罗德
// @Date: 2024/5/29
func toMapSlice(values *[]map[string]interface{}, resultSet *ResultSet) error {
	// 检查目标字典切片是否为nil，如果是则返回预定义错误
	if values == nil {
		return NilPointError
	}

	// 获取 ResultSet 的所有列名
	cols := resultSet.GetColNames()

	// 初始化一个新的字典切片，大小与 ResultSet 的行数相同
	_values := make([]map[string]interface{}, resultSet.GetRowSize())

	// 遍历 ResultSet 的每一行
	for i, row := range resultSet.GetRows() {
		// 为当前行创建一个新的字典
		_values[i] = make(map[string]interface{})
		// 遍历列名，将当前行的每列值转换并存入字典
		for j, col := range cols {
			_values[i][col] = utils.NValueToInterface(row.Values[j])
		}
	}

	// 将所有转换后的数据追加到原始的values切片中
	*values = append(*values, _values...)

	return nil
}

// 将 ResultSet 数据转换并填充到一个字典切片(*[]map[string]string)中。
// 如果resultSet的字段结果类型不是string, 将会赋值为""
//
// 参数:
// - values: 指向字典切片的指针，用于存储转换后的 ResultSet 数据。
// - resultSet: 包含查询结果的 ResultSet 实例。
//
// @Author: 罗德
// @Date: 2024/5/29
func toMapStringSlice(values *[]map[string]string, resultSet *ResultSet) error {
	// 检查目标字典切片是否为nil，如果是则返回预定义错误
	if values == nil {
		return NilPointError
	}

	// 获取 ResultSet 的所有列名
	cols := resultSet.GetColNames()

	// 初始化一个新的字典切片，大小与 ResultSet 的行数相同
	_values := make([]map[string]string, resultSet.GetRowSize())

	// 遍历 ResultSet 的每一行
	for i, row := range resultSet.GetRows() {
		// 为当前行创建一个新的字典
		_values[i] = make(map[string]string)
		// 遍历列名，将当前行的每列值转换并存入字典
		for j, col := range cols {
			// 将每一列的值转换为 string 类型，并存入目标字典中，键为列名
			valueToInterface := utils.NValueToInterface(row.Values[j])
			if str, ok := valueToInterface.(string); ok {
				_values[i][col] = str
			}
		}
	}

	// 将所有转换后的数据追加到原始的values切片中
	*values = append(*values, _values...)

	return nil
}
