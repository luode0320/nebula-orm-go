package orm

import (
	"github.com/pkg/errors"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"nebula-orm-go/config"
	"nebula-orm-go/dialectors"
)

// DB 结构体代表一个数据库连接实例，封装了与数据库交互的方法和配置。
//
// @Author: 罗德
// @Date: 2024/5/29
type DB struct {
	// 是数据库拨号器接口，用于与具体数据库建立连接。
	dialer dialectors.IDialer

	// 指向父DB实例，用于链式调用中保存状态，实现查询上下文的传递。
	parent *DB

	// 存储当前构建的SQL语句。
	sql string

	// 查询限制
	limit int

	// 是一个可执行函数，用于在某些操作完成后执行清理工作，如事务回滚或关闭连接。
	teardown func()

	// 提供日志记录功能，用于记录查询日志等。
	logger nebula.Logger

	// 标记是否开启调试模式，若为true，则会输出执行的SQL语句。
	debug bool
}

// Open 初始化并返回一个新的api.DB实例，同时根据提供的配置和选项配置数据库连接。
// 参数iDialer用于创建数据库连接，cfg包含配置信息，opts是可选配置项。
//
// @Author: 罗德
// @Date: 2024/5/29
func Open(iDialer dialectors.IDialer, cfg config.Config, opts ...config.Option) (*DB, error) {
	if iDialer == nil {
		// 如果未提供拨号器，则返回错误。
		return &DB{}, errors.New("必须创建拨号配置")
	}
	// 应用所有选项来修改配置。
	for _, opt := range opts {
		opt(&cfg)
	}
	// 加载默认配置项。
	cfg.LoadDefault()

	// 创建并返回api.DB实例。
	return &DB{
		dialer:   iDialer,
		parent:   nil,
		logger:   cfg.Logger,
		debug:    cfg.DebugMode,
		limit:    cfg.Limit,
		teardown: func() {},
	}, nil
}

// Close 关闭数据库连接。
//
// @Author: 罗德
// @Date: 2024/5/29
func (db *DB) Close() {
	db.dialer.Close()
}

// DebugMode 开启或关闭调试模式。当开启时，会在执行SQL语句前打印SQL到日志。
//
// @Author: 罗德
// @Date: 2024/5/29
func (db *DB) DebugMode() {
	db.debug = true
}

// 为当前调用链返回一个新的DB实例副本，确保每个链式调用都是独立的，
// 避免状态污染。这是实现链式调用的关键方法，每次调用都会基于当前实例创建一个新的实例。
//
// @Author: 罗德
// @Date: 2024/5/29
func (db *DB) getInstance() (tx *DB) {
	if db.parent == nil {
		// 如果当前实例没有父实例（即链的开始），则创建一个新的实例并返回。
		tx = &DB{
			dialer:   db.dialer,
			parent:   db,
			sql:      db.sql,
			logger:   db.logger,
			debug:    db.debug,
			limit:    db.limit,
			teardown: func() {},
		}
		return tx
	}
	// 若已有父实例，则返回当前实例自身作为新的链起点。
	return db
}

// execute 真正执行一个 sql
//
// @Author: 罗德
// @Date: 2024/5/29
func (db *DB) execute(sql string) (*dialectors.ResultSet, error) {
	tx := db.getInstance()

	tx.sql = sql
	if tx.debug {
		tx.logger.Info(tx.sql)
	}
	defer tx.teardown()

	result, err := tx.dialer.Execute(sql)
	if err != nil {
		return &dialectors.ResultSet{}, err
	}

	return result, nil
}
