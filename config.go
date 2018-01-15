package querydb

import (
	"time"
)

var configs map[string]Config
var exctlog bool

type Config struct {
	Username   string   //账号 root
	Password   string   //密码
	Host       string   //host localhost
	Port       string   //端口 3306
	Charset    string   //字符编码 utf8mb4
	Database   string   //默认连接数据库
	Autobranch string   //是否分库 分库数量 默认为0
	Reads      []Config //是否有从库 []
	Maxopen    int      //打开连接数 默认0
	Maxnum     int      //最大连接数 默认2
}

//设置DB配置
func SetConfig(cfgs map[string]Config) {
	configs = cfgs
}

func (config Config) Dns() string {
	return config.Username + ":" +
		config.Password + "@tcp(" +
		config.Host + ":" +
		config.Port + ")/" +
		config.Database + "?charset=" +
		config.Charset + "&loc=" + time.Local.String()
}
func (config *Config) CopyConfig(c Config) {
	if c.Username != "" {
		config.Username = c.Username
	}
	if c.Password != "" {
		config.Password = c.Password
	}
	if c.Host != "" {
		config.Host = c.Host
	}
	if c.Port != "" {
		config.Port = c.Port
	}
}
func SetExecLog(b bool) {
	exctlog = b
}
func WriteExecLog(b bool, sql Sql) {
	if b || exctlog {
		dblog.Println(sql.ToString())
	}
}
