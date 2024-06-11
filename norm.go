package nebula_orm_go

import (
	"nebula-orm-go/config"
	"nebula-orm-go/dialectors"
	"nebula-orm-go/orm"
)

// MustOpen 是Open的便捷版本，如果打开数据库时发生错误，则直接panic。
// 适合在数据库连接是程序运行前提条件的场景下使用。
//
// @Author: 罗德
// @Date: 2024/5/27
func MustOpen(iDialer dialectors.IDialer, cfg config.Config, opts ...config.Option) *orm.DB {
	db, err := orm.Open(iDialer, cfg, opts...)
	if err != nil {
		// 在遇到错误时立即终止程序。
		panic(err)
	}

	return db
}
