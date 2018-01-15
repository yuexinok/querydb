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

获取连接：

```go
//获取一个连接实例 name实例名称，isread是否只读，modnum如果分库的话填写分库因子 默认为0
func GetConn(name string, isread bool, modnum int) *QueryDb
```

```go
//获取读实例
dbread := querydb.GetConn("base", true, 2018)
//获取读写实例
dbwrite := querydb.GetConn("base", false, 2018)
//标准获取
crm := querydb.GetConn("crm", false, 0)
```

返回的queryDb实现了接口：

```go
type Connection interface {
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) (*Rows, error)
	NewQuery() *QueryBuilder
	GetLastSql() Sql
}
```

实例QueryDb有Begin()方法，返回一个QueryTx实例，QueryTx实例也实现了Connection接口

```go
func (querydb *QueryDb) Begin() (*QueryTx, error)
```



### 查询：



### 插入：

```go
//返回受影响行数和错误
func (query *QueryBuilder) Insert(datas ...map[string]interface{}) (int64, error)
//返回对应的自增id
func (query *QueryBuilder) InsertGetId(datas map[string]interface{}) (int64, error)
```

```go
data := map[string]interface{}{"f_title": "标题"}
//获取自增ID插入
id, err := db.Table("d_ec_user.t_tags").InsertGetId(data)
fmt.Println(id, err)

//单条插入
id, err = db.Table("d_ec_user.t_tags").Insert(data)
fmt.Println(id, err)

//批量插入
num, err := db.Table("d_ec_user.t_tags").Insert(data, data)
fmt.Println(num, err)
```



### 更新：

返回受影响行数，和错误

```go
func (query *QueryBuilder) Update(datas map[string]interface{}) (int64, error)

//插入更新ON DUPLICATE KEY UPDATE
//第一个参数为插入的数据，第一个参数为如果数据操作要更新的数据
func (query *QueryBuilder) InsertUpdate(insert map[string]interface{}, update map[string]interface{}) (int64, error)

//替换
func (query *QueryBuilder) Replace(datas ...map[string]interface{}) (int64, error)
```

```go
num1, err1 := db.Table("d_ec_user.t_tags").
   Where("f_tag_id", "<=", 8).
   Where("f_count", 0).
   Limit(2).
   OrderBy("f_tag_id", "desc").
   Update(map[string]interface{}{"f_title": `更换的表体"双引号",'单引号'`})

   //UPDATE d_ec_user.t_tags SET f_title = \"更换的表体\\\"双引号\\\",'单引号'\" WHERE  f_tag_id <= \"8\" AND f_count = \"0\" ORDER BY f_tag_id DESC LIMIT 2"
fmt.Println(num1, err1)

//insertupdate
num2, err2 := db.Table("d_ec_user.t_tags").InsertUpdate(map[string]interface{}{"f_title": "插入的数据", "f_tag_id": 100}, map[string]interface{}{"f_count": querydb.NewEpr("f_count+2")})
	fmt.Println(num2, err2)

//replace
num3, err3 := db.Table("d_ec_user.t_tags").Replace(map[string]interface{}{"f_title": `替换的数据`, "f_tag_id": 100})
	fmt.Println(num3, err3)
```

针对更新提供`querydb.NewEpr`(data string) 用于原生支持db相关操作：

```go
querydb.NewEpr("f_count+2")
```



### 删除：

返回被删除的行数和错误

```go
func (query *QueryBuilder) Delete() (int64, error)
```

```Go
deletenum, err := db.Table("d_ec_user.t_tags").In("f_tag_id", []interface{}{1, 2, 3}...).Delete()
deletenum, err = db.Table("d_ec_user.t_tags").In("f_tag_id", 1, 2, 3).Delete()
fmt.Println(deletenum, err)
```

### 事务：





### 调试：

直接调用db实例：的GetLastSql

```go
type Sql struct {
	Sql      string
	Args     []interface{}
	CostTime time.Duration
}
func (querydb *QueryDb) GetLastSql() Sql
```

```go
fmt.Println(db.GetLastSql())//{INSERT INTO d_ec_crm.t_crm_change  (f_crm_id) VALUES (?) [1236] 3.511848ms}

//或者直接输出完整sql
fmt.Println(db.GetLastSql().ToString())//INSERT INTO d_ec_crm.t_crm_change  (f_crm_id) VALUES ("1236")
//json格式
fmt.Println(db.GetLastSql().ToJson()) //{"sql":"SELECT f_crm_id FROM d_ec_crm.t_crm_change WHERE  f_crm_id = \"1236\" LIMIT 0,1","costtime":"1.050622ms"}


```

日志输出：

```go
//设置Logger
querydb.SetLogger(log.New(os.Stdout, "", log.Ldate))
//打印执行日志
querydb.SetExecLog(true)
```

### 注意事项：