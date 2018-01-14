package querydb

var dblog DbLog

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
