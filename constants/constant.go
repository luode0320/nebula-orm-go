package constants

import "time"

// Direction 是一个字符串类型，用于表示图数据库查询的方向
type Direction = string

// Policy 是一个整型类型，用于定义不同的策略选项
type Policy int

// 下面定义了Direction类型的常量，用于指示查询方向
const (
	// DirectionReversely 表示反向查询，即沿着边的方向相反地遍历
	DirectionReversely = "REVERSELY"

	// DirectionBidirect 表示双向查询，即同时考虑正向和反向的边进行遍历
	DirectionBidirect = "BIDIRECT"
)

// StructTagName 是一个常量，表示在结构体标记（tag）中用于指定特定行为或映射规则的键名
const StructTagName = "nebula"

// 下面定义了Policy类型的常量，用于选择不同的策略
const (
	// PolicyNothing 表示不做任何特殊处理的策略，默认策略
	PolicyNothing = iota // 使用iota枚举器，初始值为0

	// PolicyHash 表示应用哈希策略的选项，相较于PolicyNothing，可能涉及数据处理或优化的不同方式
	PolicyHash // 自动赋予下一个整数值，即1
)

// 常量定义了默认的超时、空闲时间和最大连接池大小、限制查询记录
const (
	DefaultTimeout         = 60 * time.Second // 默认请求超时时间为60秒
	DefaultIdleTime        = 10 * time.Minute // 默认连接空闲时间10分钟，超过此时间的空闲连接将被关闭
	DefaultMaxConnPoolSize = 20               // 默认的最大连接池大小为20
	DefaultLimit           = 1000             // 限制查询记录
)

// 常量定义了查询的返回别名
const (
	V = "v" // 默认解析 as v 的别名
	E = "e" // 默认解析 as e 的别名
)
