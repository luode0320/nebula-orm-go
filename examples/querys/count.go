package querys

import (
	"log"
	"nebula-orm-go/orm"
)

// Count 统计
//
// @Author: 罗德
// @Date: 2024/6/11
func Count(db *orm.DB) {
	nsql := "match(v:user)-[e]->(v2) where id(v)==hash('user_101') return count(e)"
	cnt := 0
	err := db.Debug().ExecuteAndParse(nsql, &cnt)
	if err != nil {
		log.Fatalf("exec %s error: %v", nsql, err)
		panic(err)
	}
	log.Printf("count: %d", cnt)
}
