package dialectors

import (
	"fmt"
	"github.com/pkg/errors"
	nebula "github.com/vesoft-inc/nebula-go/v3" // 导入Nebula Go客户端库
	"nebula-orm-go/config"
	"strconv"
	"strings"
)

var _ IDialer = new(NebulaDialer)

// IDialer 接口定义了与Nebula Graph交互的基本方法，包括执行SQL和关闭连接池
//
// @Author: 罗德
// @Date: 2024/5/24
type IDialer interface {
	Execute(sql string) (*ResultSet, error) // 执行SQL语句并返回结果集及潜在错误
	Close()                                 // 关闭连接池
}

// NebulaDialer 结构体，用于管理Nebula Graph的连接和会话
//
// @Author: 罗德
// @Date: 2024/5/24
type NebulaDialer struct {
	pool     *nebula.ConnectionPool // 连接池实例
	username string                 // 认证用户名
	password string                 // 认证密码
	space    string                 // 当前操作的图空间名
}

// NewNebulaDialer 创建一个新的NebulaDialer实例
//
// @Author: 罗德
// @Date: 2024/5/24
func NewNebulaDialer(cfg config.DialerConfig) (*NebulaDialer, error) {
	// 加载默认配置项
	cfg.LoadDefault()
	// 解析地址列表为Nebula所需的HostAddress格式
	nAddresses, err := parseAddresses(cfg.Addresses)
	if err != nil {
		return nil, err
	}

	// 初始化Nebula连接池配置
	nConfig := nebula.PoolConfig{
		TimeOut:         cfg.Timeout,
		IdleTime:        cfg.IdleTime,
		MaxConnPoolSize: cfg.MaxConnPoolSize,
		MinConnPoolSize: cfg.MinConnPoolSize,
	}

	// 使用配置创建连接池
	nPool, err := nebula.NewConnectionPool(nAddresses, nConfig, nebula.DefaultLogger{})
	if err != nil {
		return nil, errors.Wrap(err, "连接Nebula失败")
	}

	// 返回NebulaDialer实例
	return &NebulaDialer{
		pool:     nPool,
		username: cfg.Username,
		password: cfg.Password,
	}, nil
}

// MustNewNebulaDialer 确保新建NebulaDialer实例，如果失败则panic
//
// @Author: 罗德
// @Date: 2024/5/24
func MustNewNebulaDialer(cfg config.DialerConfig) *NebulaDialer {
	// 创建实例
	dialer, err := NewNebulaDialer(cfg)
	if err != nil {
		panic(err)
	}

	// 初始化sql(立刻尝试使用刚创建的空间、点、边可能会失败，因为创建是异步实现的)
	if cfg.InitSql != "" {
		if _, err := dialer.Execute(cfg.InitSql); err != nil {
			panic(err)
		}
	}

	dialer.space = cfg.Space

	return dialer
}

// Execute 执行SQL语句，内部处理了连接获取、空间切换、执行语句及结果检查
//
// @Author: 罗德
// @Date: 2024/5/24
func (d *NebulaDialer) Execute(sql string) (*ResultSet, error) {
	// 获取会话
	session, err := d.getSession()
	if err != nil {
		return &ResultSet{}, err
	}
	defer session.Release()

	// 使用指定的图空间
	if d.space != "" {
		_, err = session.Execute("use " + d.space)
		if err != nil {
			return &ResultSet{}, err
		}
	}

	// 执行SQL语句
	result, err := session.Execute(sql)
	if err != nil {
		return &ResultSet{}, err
	}

	// 检查结果集是否执行成功
	if err = checkResultSet(result); err != nil {
		return &ResultSet{}, err
	}

	// 封装并返回结果集
	return &ResultSet{result}, nil
}

// CreateSpace 创建图空间
//
// @Author: 罗德
// @Date: 2024/5/24
func (d *NebulaDialer) CreateSpace(space string) error {
	// 参数校验
	if space == "" {
		return fmt.Errorf("初始化图空间失败：space 参数无效")
	}

	// 创建空间
	createSpaceSQL := fmt.Sprintf("CREATE SPACE IF NOT EXISTS %s(VID_TYPE=FIXED_STRING(64));", space)
	if _, err := d.Execute(createSpaceSQL); err != nil {
		return err
	}

	d.space = space

	return nil
}

// getSession 从连接池获取会话
//
// @Author: 罗德
// @Date: 2024/5/24
func (d *NebulaDialer) getSession() (*nebula.Session, error) {
	return d.pool.GetSession(d.username, d.password)
}

// Close 关闭连接池
//
// @Author: 罗德
// @Date: 2024/5/24
func (d *NebulaDialer) Close() {
	d.pool.Close()
}

// checkResultSet 检查结果集是否执行成功，根据错误码判断
//
// @Author: 罗德
// @Date: 2024/5/24
func checkResultSet(nSet *nebula.ResultSet) error {
	if nSet.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
		return errors.New(fmt.Sprintf("code: %d, msg: %s",
			nSet.GetErrorCode(), nSet.GetErrorMsg()))
	}
	if !nSet.IsSucceed() {
		return errors.New(fmt.Sprintf("code: %d, msg: %s",
			nSet.GetErrorCode(), nSet.GetErrorMsg()))
	}
	return nil
}

// parseAddresses 解析地址字符串列表为Nebula所需的HostAddress格式
//
// @Author: 罗德
// @Date: 2024/5/24
func parseAddresses(addresses []string) ([]nebula.HostAddress, error) {
	hostAddresses := make([]nebula.HostAddress, len(addresses))
	for i, addr := range addresses {
		list := strings.Split(addr, ":")
		if len(list) < 2 {
			return []nebula.HostAddress{},
				errors.New(fmt.Sprintf("address %s invalid", addr))
		}
		port, err := strconv.ParseInt(list[1], 10, 64)
		if err != nil {
			return []nebula.HostAddress{},
				errors.New(fmt.Sprintf("address %s invalid", addr))
		}
		hostAddresses[i] = nebula.HostAddress{
			Host: list[0],
			Port: int(port),
		}
	}
	return hostAddresses, nil
}
