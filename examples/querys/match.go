package querys

import (
	"fmt"
	"log"
	"nebula-orm-go/orm"

	"nebula-orm-go/examples/models"
)

// MatchSingle 根据ID匹配单个用户顶点并打印其信息
//
// @Author: 罗德
// @Date: 2024/5/27
func MatchSingle(db *orm.DB, sdkVertex *models.SdkVertex) error {
	// 定义Cypher查询语句，匹配标签%s顶点，其中id等于%s，并返回所有字段
	nsql := fmt.Sprintf("match(v:%s) where id(v)=='%s' return v", sdkVertex.TagName(), sdkVertex.GetVid())
	// 定义一个SdkVertex类型的变量用于接收查询结果
	user := models.SdkVertex{}
	// 执行查询语句并解析结果到user变量中，如果执行出错则输出错误信息并终止程序
	err := db.Debug().ExecuteAndParse(nsql, &user)
	if err != nil {
		log.Fatalf("执行查询 %s 时发生错误: %v", nsql, err)
		return err
	}

	// 成功后打印查询到的用户信息
	log.Printf("%+v", user)
	return nil
}

// MatchMulti 根据ID匹配多个用户顶点并打印其信息
//
// @Author: 罗德
// @Date: 2024/6/11
func MatchMulti(db *orm.DB) {
	// 定义Cypher查询语句，匹配标签为'user'的顶点，其中id为'user_101'或'user_100'，并返回id和created字段
	nsql := "match(v:user) where id(v)==hash('user_101') or id(v)==hash('user_100')" +
		" return v.id as id,v.created as created"
	// 定义一个SdkVertex类型的切片用于接收多个查询结果
	users := []models.SdkVertex{}
	// 执行查询语句并解析结果到users切片中，如果执行出错则输出错误信息并终止程序
	err := db.Debug().ExecuteAndParse(nsql, &users)
	if err != nil {
		log.Fatalf("执行查询 %s 时发生错误: %v", nsql, err)
		panic(err)
	}
	// 成功后打印查询到的所有用户信息
	log.Printf("%+v", users)
}
