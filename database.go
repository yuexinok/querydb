package querydb

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Rows = sql.Rows
type Result = sql.Result
type Connection interface {
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) (*Rows, error)
	NewQuery() *QueryBuilder
	GetLastSql() Sql
}
type Sql struct {
	Sql      string
	Args     []interface{}
	CostTime time.Duration
}
type QueryDb struct {
	db         sql.DB
	modnum     int
	configname string
	config     Config
	lastsql    Sql
}

func ToMap(rows *Rows) []map[string]interface{} {
	cols, err := rows.Columns()
	if err != nil {
		dblog.Println(err)
		return nil
	}
	count := len(cols)

	var data []map[string]interface{}
	vals := make([]string, count)
	ptr := make([]interface{}, count)
	for i := 0; i < count; i++ {
		ptr[i] = &vals[i]
	}
	defer rows.Close()
	for rows.Next() {
		//字段
		entry := make(map[string]interface{}, count)

		err = rows.Scan(ptr...)
		if err != nil {
			dblog.Println(err)
		}
		for i, col := range cols {
			entry[col] = vals[i]
		}
		data = append(data, entry)
	}
	if err = rows.Err(); err != nil {
		dblog.Println(err)
	}
	return data
}

//生成一个新的查询构造器
func (querydb *QueryDb) NewQuery() *QueryBuilder {
	return &QueryBuilder{connection: querydb}
}

//查询构造器快速调用
func (querydb *QueryDb) Table(tablename ...string) *QueryBuilder {
	query := &QueryBuilder{connection: querydb}
	return query.Table(tablename...)
}

//开启一个事务
func (querydb *QueryDb) Begin() (*QueryTx, error) {
	tx, err := querydb.db.Begin()
	if err != nil {
		return nil, err
	}
	return &QueryTx{tx: *tx, modnum: querydb.modnum, configname: querydb.configname, config: querydb.config}, nil
}

//复用执行语句
func (querydb *QueryDb) Exec(query string, args ...interface{}) (Result, error) {
	querydb.lastsql.Sql = query
	querydb.lastsql.Args = args
	start := time.Now()
	defer func() {
		querydb.lastsql.CostTime = time.Since(start)
		WriteExecLog(querydb.lastsql)
	}()
	return querydb.db.Exec(query, args...)
}

//复用查询语句
func (querydb *QueryDb) Query(query string, args ...interface{}) (*Rows, error) {
	querydb.lastsql.Sql = query
	querydb.lastsql.Args = args
	start := time.Now()
	defer func() {
		querydb.lastsql.CostTime = time.Since(start)
		WriteExecLog(querydb.lastsql)

	}()
	return querydb.db.Query(query, args...)
}

func (querydb *QueryDb) GetLastSql() Sql {

	return querydb.lastsql
}

type QueryTx struct {
	tx         sql.Tx
	modnum     int
	configname string
	config     Config
	lastsql    Sql
}

func (querytx *QueryTx) Commit() error {
	return querytx.tx.Commit()
}

func (querytx *QueryTx) Rollback() error {
	return querytx.tx.Rollback()
}
func (querytx *QueryTx) NewQuery() *QueryBuilder {
	return &QueryBuilder{connection: querytx}
}
func (querytx *QueryTx) Table(tablename ...string) *QueryBuilder {
	query := &QueryBuilder{connection: querytx}
	return query.Table(tablename...)
}
func (querytx *QueryTx) Exec(query string, args ...interface{}) (Result, error) {
	querytx.lastsql.Sql = query
	querytx.lastsql.Args = args
	start := time.Now()
	defer func() {
		querytx.lastsql.CostTime = time.Since(start)
		WriteExecLog(querytx.lastsql)

	}()
	return querytx.Exec(query, args...)
}
func (querytx *QueryTx) Query(query string, args ...interface{}) (*Rows, error) {
	querytx.lastsql.Sql = query
	querytx.lastsql.Args = args
	start := time.Now()
	defer func() {
		querytx.lastsql.CostTime = time.Since(start)
		WriteExecLog(querytx.lastsql)

	}()
	return querytx.Query(query, args...)
}
func (querytx *QueryTx) GetLastSql() Sql {
	return querytx.lastsql
}

func (sql Sql) ToString() string {
	s := sql.Sql
	for _, v := range sql.Args {
		val := fmt.Sprintf("%v", v)
		val = strconv.Quote(val)
		s = strings.Replace(s, "?", val, 1)
	}
	return s
}
func (sql Sql) ToJson() string {
	return fmt.Sprintf(`{"sql":%s,"costtime":"%s"}`, strconv.Quote(sql.ToString()), sql.CostTime)
}
