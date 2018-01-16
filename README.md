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

```go
//多条查询  返回的Rows和sql.Rows用户一致
func (query *QueryBuilder) GetRows() (*Rows, error)
//单条查询 dest用法和sql.QueryRow()用法一致
func (query *QueryBuilder) GetRow(dest ...interface{}) error
//总数查询
func (query *QueryBuilder) Count() (int64, error)

//提供辅助函数 把rows转化成对应的map
func ToMap(rows *Rows) []map[string]interface{}

//提供辅助函数 把rows转化成对应的 struct 切片
//TODO
```

基本用法：

```go
//多条查询
rows, err := db.Table("d_ec_user.t_tags").Where("f_tag_id", 5).GetRows()
if err == nil {
   //生成map
   list := querydb.ToMap(rows)
   fmt.Println(list)
}

//单条查询
var title string
	err = db.Table("d_ec_user.t_tags").Select("f_title").Where("f_tag_id", "=", 5).GetRow(&title)
	fmt.Println(err, title)
```

#### 复杂用法：

##### 多表查询：

```go
//设置操作的表名称
func (query *QueryBuilder) Table(tablename ...string) *QueryBuilder
```

```go

//范例
rows, err = crm.Table("d_ec_crm.t_eccrm_detail as d", "d_ec_crm.t_crm_relation as r").
		Select("d.f_name", "r.f_user_id").
		Where("d.f_crm_id=r.f_crm_id").
		Where("d.f_crm_id", 232740452).
		GetRows()
//SELECT d.f_name,r.f_user_id FROM d_ec_crm.t_eccrm_detail as d,d_ec_crm.t_crm_relation as r WHERE d.f_crm_id=r.f_crm_id AND d.f_crm_id = "232740452"
```

##### 自定义查询列：

```go
func (query *QueryBuilder) Select(columns ...string) *QueryBuilder
```

```go
var max, min int
err = crm.Table("d_ec_crm.t_eccrm_detail").
   Select("max(f_crm_id),min(f_crm_id)").
   GetRow(&max, &min)
```

即只要满足为string，select可以支持sql的各类复杂用法

##### where：

```go
//and 
func (query *QueryBuilder) Where(column string, value ...interface{}) *QueryBuilder
//or
func (query *QueryBuilder) OrWhere(column string, value ...interface{}) *QueryBuilder

//相等
func (query *QueryBuilder) Equal(column string, value interface{}) *QueryBuilder
func (query *QueryBuilder) OrEqual(column string, value interface{}) *QueryBuilder
//不相等
func (query *QueryBuilder) NotEqual(column string, value interface{}) *QueryBuilder
func (query *QueryBuilder) OrNotEqual(column string, value interface{}) *QueryBuilder
```

一个参数的时候为原生where，2个参数的时候为column等,3个参数的时候第一个参数为列，第2个是操作，第3个是值

```Go
crm.Table("d_ec_crm.t_eccrm_detail").Where("f_crm_id=230740537").GetRows()
crm.Table("d_ec_crm.t_eccrm_detail").Where("f_crm_id",230740537).GetRows()
crm.Table("d_ec_crm.t_eccrm_detail").Where("f_crm_id","=",230740537).GetRows()
crm.Table("d_ec_crm.t_eccrm_detail").Equal("f_crm_id",230740537).GetRows()
```

其中操作符可以是：>,<,>=,<=等

其他特殊操作：

```go
//Between 
func (query *QueryBuilder) Between(column string, value1 interface{}, value2 interface{}) *QueryBuilder
func (query *QueryBuilder) OrBetween(column string, value1 interface{}, value2 interface{}) *QueryBuilder
func (query *QueryBuilder) NotBetween(column string, value1 interface{}, value2 interface{}) *QueryBuilder
func (query *QueryBuilder) NotOrBetween(column string, value1 interface{}, value2 interface{}) *QueryBuilder

//in
func (query *QueryBuilder) In(column string, value ...interface{}) *QueryBuilder
func (query *QueryBuilder) OrIn(column string, value ...interface{}) *QueryBuilder
func (query *QueryBuilder) NotIn(column string, value ...interface{}) *QueryBuilder
func (query *QueryBuilder) OrNotIn(column string, value ...interface{}) *QueryBuilder

//是否是空
func (query *QueryBuilder) IsNULL(column string) *QueryBuilder
func (query *QueryBuilder) OrIsNULL(column string) *QueryBuilder
func (query *QueryBuilder) IsNotNULL(column string) *QueryBuilder
func (query *QueryBuilder) OrIsNotNULL(column string) *QueryBuilder
//like查询
func (query *QueryBuilder) Like(column string, value interface{}) *QueryBuilder
func (query *QueryBuilder) OrLike(column string, value interface{}) *QueryBuilder
```

```go
var crmname string
err = crm.Table("d_ec_crm.t_eccrm_detail").
	Select("f_name").
	Between("f_crm_id", 230740537, 230740560).
	GetRow(&crmname)
//SELECT f_name FROM d_ec_crm.t_eccrm_detail WHERE  f_crm_id BETWEEN "230740537" AND "230740560"

err = crm.Table("d_ec_crm.t_eccrm_detail").Select("f_name").In("f_crm_id", 230740537, 230740560).GetRow(&crmname)
err = crm.Table("d_ec_crm.t_eccrm_detail").Select("f_name").In("f_crm_id", []interface{}{230740537, 230740560}...).GetRow(&crmname)

err = crm.Table("d_ec_crm.t_eccrm_detail").Select("f_name").Like("f_name", "李%").GetRow(&crmname)
```

##### Limit,OrderBy,GroupBy,Skip,Distinct：

```go
func (query *QueryBuilder) Distinct() *QueryBuilder
func (query *QueryBuilder) GroupBy(groups ...string) *QueryBuilder
//可以多次调用
func (query *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder
func (query *QueryBuilder) Offset(offset int) *QueryBuilder
//同Offset
func (query *QueryBuilder) Skip(offset int) *QueryBuilder
func (query *QueryBuilder) Limit(limit int) *QueryBuilder
```



##### Join，LeftJoin RightJoin:

```go
func (query *QueryBuilder) Join(tablename string, on string) *QueryBuilder
func (query *QueryBuilder) LeftJoin(tablename string, on string) *QueryBuilder
func (query *QueryBuilder) RightJoin(tablename string, on string) *QueryBuilder
```

```go
var crmname string
var crmid int64
err = crm.Table("d_ec_crm.t_eccrm_detail as d").
	Join("d_ec_crm.t_crm_relation as r", "d.f_crm_id=r.f_crm_id").
	Select("d.f_name", "d.f_crm_id").
	Where("d.f_corp_id", 21299).
	Where("r.f_user_id", 0).
	GetRow(&crmname, &crmid)

//SELECT d.f_name,d.f_crm_id FROM d_ec_crm.t_eccrm_detail as d JOIN d_ec_crm.t_crm_relation as r ON d.f_crm_id=r.f_crm_id WHERE  d.f_corp_id = "21299" AND r.f_user_id = "0" LIMIT 0,1
```

##### Union,UnionAll

```go
func (query *QueryBuilder) Union(unions ...QueryBuilder) *QueryBuilder 
func (query *QueryBuilder) UnionAll(unions ...QueryBuilder) *QueryBuilder
```



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