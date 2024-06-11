package config

import (
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"nebula-orm-go/constants"
)

// Option 定义了一个函数类型，该类型的函数接收一个指向Config的指针，并对其进行修改以应用特定的配置选项。
type Option func(c *Config)

// Config 结构体定义了与Nebula Graph数据库交互时的额外配置选项，包括调试模式标志、日志记录器、限制查询
//
// @Author: 罗德
// @Date: 2024/6/11
type Config struct {
	Limit     int           // Limit 限制查询记录
	DebugMode bool          // DebugMode 指示是否开启调试模式，若为true，则可能会输出额外的日志信息以帮助调试。
	Logger    nebula.Logger // Logger 提供日志记录功能的接口，用于记录与数据库交互过程中的信息。
}

// LoadDefault 方法为Config结构体提供了默认配置加载逻辑。
//
// @Author: 罗德
// @Date: 2024/6/11
func (config *Config) LoadDefault() {
	if config.Logger == nil {
		config.Logger = &nebula.DefaultLogger{} // 当用户未提供日志记录器时，使用Nebula官方的默认日志记录器。
	}
	if config.Limit <= 0 {
		config.Limit = constants.DefaultLimit // 若限制查询记录不合理，则使用默认值
	}
}
