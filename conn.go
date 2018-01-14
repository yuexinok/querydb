package querydb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

var connections map[string]*QueryDb

func GetConn(name string, isread bool, modnum int) *QueryDb {

	config, isok := configs[name]
	if !isok {
		dblog.Fatal("DB配置:" + name + "找不到！")
	}
	//是否分实例
	name = GetBranchName(name, modnum, config.Autobranch)
	//是否读从
	readlen := len(config.Reads)
	keyname := name
	readnum := 0
	if isread && readlen > 0 {
		readnum = GetReadNumByRand(readlen)
		keyname += "_read_" + strconv.Itoa(readnum)

	}
	_, ok := connections[keyname]
	if ok {
		return connections[keyname]
	}

	//复用都配置
	if isread && readlen > 0 {
		config.CopyConfig(config.Reads[readnum])
	}

	db, err := sql.Open("mysql",
		config.Dns())

	if err != nil {
		dblog.Fatal("DB连接错误！")
	}
	if config.Maxopen > 0 {
		db.SetMaxOpenConns(config.Maxopen)
	}
	if config.Maxnum > 0 {
		db.SetMaxIdleConns(config.Maxnum)
	}
	//确保关闭
	defer db.Close()
	return &QueryDb{db: *db, modnum: modnum, configname: name, config: config}
}
