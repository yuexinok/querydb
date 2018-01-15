# querydb

这是一个针对go mysql查询的查询构造器，支持主从配置，支持读写分离，支持分库配置。参照PHP `Laravel` 框架database编写，使用简单，并且对`database/sql`进行了简单的使用封装，如果你有使用上的问题和建议，欢迎联系我。

### 获取

git clone https://github.com/yuexinok/querydb

使用方式：**go get github.com/yuexinok/querydb**

### 配置：

参看querydb.Config 提供参数：

```go
type Config struct {
	Username   string   //账号 root
	Password   string   //密码
	Host       string   //host localhost
	Port       string   //端口 3306
	Charset    string   //字符编码 utf8mb4
	Database   string   //默认连接数据库
	Autobranch string   //是否分库 分库的个数
	Reads      []Config //是否有从库 []
	Maxopen    int //打开连接数 默认0
	Maxnum     int //最大连接数 默认2
}
```

配置范例：

```go
	//设置Logger
	querydb.SetLogger(log.New(os.Stdout, "", log.Ldate))

	//配置集合
	configs := map[string]querydb.Config{}

	reads := make([]querydb.Config, 1)

	reads[0] = querydb.Config{Username: "read", Password: "123456", Host: "127.0.0.1", Port: "3306", Charset: "utf8mb4"}

	//标准
	configs["crm"] = querydb.Config{Username: "root", Password: "123456", Host: "127.0.0.1", Port: "3306", Charset: "utf8mb4", Database: "d_ec_crmextend"}

	//分库
	configs["base"] = querydb.Config{Autobranch: "2"}
	configs["base0"] = querydb.Config{Username: "root", Password: "123456", Host: "127.0.0.1", Port: "3306", Charset: "utf8mb4",  Database: "d_ec_crm", Reads: reads}
	configs["base1"] = querydb.Config{Username: "root", Password: "123456", Host: "127.0.0.1", Port: "3306", Charset: "utf8mb4",  Database: "d_ec_crm", Reads: reads}

	//读写分离
	configs["user"] = querydb.Config{Username: "root", Password: "123456", Host: "127.0.0.1", Port: "3306", Charset: "utf8mb4", Database: "d_ec_user", Reads: reads}

	//初始化设置
	querydb.SetConfig(configs)
```



### 连接：



### 查询：



### 插入：



### 更新：



### 删除：

### 

### 事务：

### 

### 注意事项：