package querys

import (
	"fmt"
	"log"
	"nebula-orm-go/orm"
	"time"
)

// Chainable 使用 chainable 必须保证函数调用顺序, 因为逻辑上是有序的.
//
// @Author: 罗德
// @Date: 2024/6/11
func Chainable(db *orm.DB) {
	cnt := int64(0)
	queryWhere := fmt.Sprintf("%s.created > %d", voteExample.EdgeName(), time.Now().Unix())
	err := db.Debug().From(userExample).Over(voteExample).Bidirect().
		Where(queryWhere).Yield("'' as id").
		Group("id").Yield("count(1)").Return(&cnt)
	if err != nil {
		panic(err)
	}
	log.Printf("count: %d", cnt)
}
