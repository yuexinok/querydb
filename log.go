package querydb

var dblog DbLog
var exctlog bool

type DbLog interface {
	Println(v ...interface{}) //普通
	Fatal(v ...interface{})   //致命
}

func SetLogger(log DbLog) {
	dblog = log
}
func GetLogger() DbLog {
	return dblog
}

func SetExecLog(b bool) {
	exctlog = b
}
func WriteExecLog(sql Sql) {
	if exctlog {
		dblog.Println(sql.ToJson())
	}
}
