package querydb

type DbError struct {
	msg string
	sql Sql
}

func NewDBError(msg string, sql Sql) DbError {
	return DbError{msg: msg, sql: sql}
}

func (e DbError) Error() string {
	return "DBError:" + e.msg + " " + e.sql.ToJson()
}
