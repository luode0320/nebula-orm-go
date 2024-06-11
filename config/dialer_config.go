package config

import (
	"nebula-orm-go/constants"
	"time" // 引入time包，用于处理时间相关的功能，如超时和空闲时间设定
)

// DialerConfig 用于配置与Nebula Graph建立连接的参数
//
// @Author: 罗德
// @Date: 2024/5/24
type DialerConfig struct {
	Addresses       []string      `json:"addresses" yaml:"addresses"`                   // Nebula Graph集群的地址列表
	Space           string        `json:"space" yaml:"space"`                           // 初始化图空间名称
	Username        string        `json:"username" yaml:"username"`                     // 认证用户名，支持JSON和YAML格式化
	Password        string        `json:"password" yaml:"password"`                     // 认证密码
	Timeout         time.Duration `json:"timeout" yaml:"timeout"`                       // 连接和读写超时时间
	IdleTime        time.Duration `json:"idle_time" yaml:"idle_time"`                   // 连接空闲时间，超过则关闭
	MaxConnPoolSize int           `json:"max_conn_pool_size" yaml:"max_conn_pool_size"` // 连接池最大连接数
	MinConnPoolSize int           `json:"min_conn_pool_size" yaml:"min_conn_pool_size"` // 连接池最小连接数
	InitSql         string        `json:"init_sql" yaml:"init_sql"`                     // 初始化sql(立刻尝试使用刚创建的空间、点、边可能会失败，因为创建是异步实现的)
}

// LoadDefault 方法用于加载默认配置项到DialerConfig实例中，如果相应配置项未被显式设置
//
// @Author: 罗德
// @Date: 2024/5/24
func (config *DialerConfig) LoadDefault() {
	if config.Timeout <= 0 {
		config.Timeout = constants.DefaultTimeout // 若超时未设置或设置为非正值，则使用默认超时时间
	}
	if config.IdleTime <= 0 {
		config.IdleTime = constants.DefaultIdleTime // 若空闲时间未设置或设置为非正值，则使用默认空闲时间
	}
	if config.MaxConnPoolSize <= 0 {
		config.MaxConnPoolSize = constants.DefaultMaxConnPoolSize // 若最大连接池大小未设置或设置不合理，则使用默认值
	}
}
