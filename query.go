package querydb

import (
	"strconv"
	"strings"
)

const (
	BETWEEN    = "BETWEEN"
	NOTBETWEEN = "NOT BETWEEN"
	IN         = "IN"
	NOTIN      = "NOT IN"
	AND        = "AND"
	OR         = "OR"
	ISNULL     = "IS NULL"
	ISNOTNULL  = "IS NOT NULL"
	EQUAL      = "="
	NOTEQUAL   = "!="
	LIKE       = "LIKE"
	JOIN       = "JOIN"
	LEFTJOIN   = "LEFT JOIN"
	RIGHTJOIN  = "RIGHT JOIN"
	UNION      = "UNION"
	UNIONALL   = "UNION ALL"
	DESC       = "DESC"
	ASC        = "ASC"
)

type QueryBuilder struct {
	connection Connection
	table      []string
	columns    []string
	where      []w
	orders     []string
	groups     []string
	limit      int
	offset     int
	distinct   bool
	binds      []string
	joins      []join
	unions     []union
	unlimmit   int
	unoffset   int
	unorders   []string

	args  []interface{}
	datas []map[string]interface{}
}
type join struct {
	table    string
	on       string
	operator string
}
type union struct {
	query    QueryBuilder
	operator string
}
type w struct {
	column   string
	operator string
	valuenum int
	do       string
}

//设置操作的表名称
func (query *QueryBuilder) Table(tablename ...string) *QueryBuilder {
	query.table = tablename
	return query
}

func (query *QueryBuilder) Select(columns ...string) *QueryBuilder {
	query.columns = columns
	return query
}

func (query *QueryBuilder) Where(column string, value ...interface{}) *QueryBuilder {
	if len(value) == 0 { //一个参数直接where
		query.toWhere(column, "", 0, AND)
	} else if len(value) == 1 { //2个参数直接where =
		query.toWhere(column, EQUAL, 1, AND)
		query.addArg(value[0])
	} else { //3个参数
		switch v := value[0].(type) {
		case string:
			query.toWhere(column, v, 1, AND)
			query.addArg(value[1])
		}
	}
	return query
}
func (query *QueryBuilder) OrWhere(column string, value ...interface{}) *QueryBuilder {
	if len(value) == 0 { //一个参数直接where
		query.toWhere(column, "", 0, OR)
	} else if len(value) == 1 { //2个参数直接where =
		query.toWhere(column, EQUAL, 1, OR)
		query.addArg(value[0])
	} else {
		switch v := value[0].(type) {
		case string:
			query.toWhere(column, v, 1, OR)
			query.addArg(value[1])
		}
	}
	return query
}
func (query *QueryBuilder) Equal(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, EQUAL, 1, AND)
	query.addArg(value)
	return query
}
func (query *QueryBuilder) OrEqual(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, EQUAL, 1, OR)
	query.addArg(value)
	return query
}
func (query *QueryBuilder) NotEqual(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, NOTEQUAL, 1, AND)
	query.addArg(value)
	return query
}
func (query *QueryBuilder) OrNotEqual(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, NOTEQUAL, 1, OR)
	query.addArg(value)
	return query
}
func (query *QueryBuilder) Between(column string, value1 interface{}, value2 interface{}) *QueryBuilder {
	query.toWhere(column, BETWEEN, 2, AND)
	query.addArg(value1, value2)
	return query
}
func (query *QueryBuilder) OrBetween(column string, value1 interface{}, value2 interface{}) *QueryBuilder {
	query.toWhere(column, BETWEEN, 2, OR)
	query.addArg(value1, value2)
	return query
}
func (query *QueryBuilder) NotBetween(column string, value1 interface{}, value2 interface{}) *QueryBuilder {
	query.toWhere(column, NOTBETWEEN, 2, AND)
	query.addArg(value1, value2)
	return query
}
func (query *QueryBuilder) NotOrBetween(column string, value1 interface{}, value2 interface{}) *QueryBuilder {
	query.toWhere(column, NOTBETWEEN, 2, OR)
	query.addArg(value1, value2)
	return query
}
func (query *QueryBuilder) In(column string, value ...interface{}) *QueryBuilder {
	query.toWhere(column, IN, len(value), AND)
	query.addArg(value...)
	return query
}
func (query *QueryBuilder) OrIn(column string, value ...interface{}) *QueryBuilder {
	query.toWhere(column, IN, len(value), OR)
	query.addArg(value...)
	return query
}
func (query *QueryBuilder) NotIn(column string, value ...interface{}) *QueryBuilder {
	query.toWhere(column, NOTIN, len(value), AND)
	query.addArg(value...)
	return query
}
func (query *QueryBuilder) OrNotIn(column string, value ...interface{}) *QueryBuilder {
	query.toWhere(column, NOTIN, len(value), OR)
	query.addArg(value...)
	return query
}
func (query *QueryBuilder) IsNULL(column string) *QueryBuilder {
	query.toWhere(column, ISNULL, 0, AND)
	return query
}
func (query *QueryBuilder) OrIsNULL(column string) *QueryBuilder {
	query.toWhere(column, ISNULL, 0, OR)
	return query
}
func (query *QueryBuilder) IsNotNULL(column string) *QueryBuilder {
	query.toWhere(column, ISNOTNULL, 0, AND)
	return query
}
func (query *QueryBuilder) OrIsNotNULL(column string) *QueryBuilder {
	query.toWhere(column, ISNOTNULL, 0, OR)
	return query
}
func (query *QueryBuilder) Like(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, LIKE, 1, AND)
	query.addArg(value)
	return query
}
func (query *QueryBuilder) OrLike(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, LIKE, 1, OR)
	query.addArg(value)
	return query
}
func (query *QueryBuilder) Join(tablename string, on string) *QueryBuilder {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: JOIN})
	return query
}
func (query *QueryBuilder) LeftJoin(tablename string, on string) *QueryBuilder {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: LEFTJOIN})
	return query
}
func (query *QueryBuilder) RightJoin(tablename string, on string) *QueryBuilder {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: RIGHTJOIN})
	return query
}
func (query *QueryBuilder) Union(unions ...QueryBuilder) *QueryBuilder {
	for i, len := 0, len(unions); i < len; i++ {
		query.unions = append(query.unions, union{query: unions[i], operator: UNION})
		query.addArg(unions[i].args...)
	}
	return query
}
func (query *QueryBuilder) UnionOffset(offset int) *QueryBuilder {
	query.unoffset = offset
	return query
}
func (query *QueryBuilder) UnionLimit(limit int) *QueryBuilder {
	query.unlimmit = limit
	return query
}
func (query *QueryBuilder) UnionOrderBy(column string, direction string) *QueryBuilder {
	if strings.ToUpper(direction) == DESC {
		column += " " + DESC
	} else {
		column += " " + ASC
	}
	query.unorders = append(query.unorders, column)
	return query
}
func (query *QueryBuilder) UnionAll(unions ...QueryBuilder) *QueryBuilder {
	for i, len := 0, len(unions); i < len; i++ {
		query.unions = append(query.unions, union{query: unions[i], operator: UNIONALL})
		query.addArg(unions[i].args...)
	}
	return query
	return query
}

func (query *QueryBuilder) Distinct() *QueryBuilder {
	query.distinct = true
	return query
}
func (query *QueryBuilder) GroupBy(groups ...string) *QueryBuilder {
	query.groups = groups
	return query
}
func (query *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder {
	if strings.ToUpper(direction) == DESC {
		column += " " + DESC
	} else {
		column += " " + ASC
	}
	query.orders = append(query.orders, column)
	return query
}
func (query *QueryBuilder) Offset(offset int) *QueryBuilder {
	query.offset = offset
	return query
}
func (query *QueryBuilder) Skip(offset int) *QueryBuilder {
	query.offset = offset
	return query
}
func (query *QueryBuilder) Limit(limit int) *QueryBuilder {
	query.limit = limit
	return query
}

func (query *QueryBuilder) ToSql(method string) string {
	grammar := Grammar{builder: query, method: method}
	return grammar.ToSql()
}
func (query *QueryBuilder) toWhere(column string, operator string, valuenum int, do string) *QueryBuilder {
	query.where = append(
		query.where,
		w{column: column, operator: operator, valuenum: valuenum, do: do})
	return query
}
func (query *QueryBuilder) addArg(value ...interface{}) {
	query.args = append(query.args, value...)
}

func (query *QueryBuilder) beforeArg(value ...interface{}) {
	query.args = append(value, query.args...)
}

func (query *QueryBuilder) setData(datas ...map[string]interface{}) {
	query.datas = datas

}

func (query *QueryBuilder) GetRows() (*Rows, error) {
	grammar := Grammar{builder: query}
	sql := grammar.Select()
	rows, err := query.connection.Query(sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		dblog.Println(err.Error())
		return nil, err
	}
	return rows, nil
}
func (query *QueryBuilder) GetRow(dest ...interface{}) error {
	query.Limit(1)
	query.Offset(0)
	rows, err := query.GetRows()
	if err != nil || rows == nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			err = NewDBError(err.Error(), query.connection.GetLastSql())
			dblog.Println(err.Error())
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		dblog.Println(err.Error())
	}
	return nil
}
func (query *QueryBuilder) Count() (int64, error) {
	query.Select("count(*) as _C")
	row, err := query.GetRows()
	if err != nil || row == nil {
		return 0, err
	}
	d := ToMap(row)
	if len(d) < 1 {
		return 0, nil
	}

	switch v := d[0]["_C"].(type) {
	case int:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 0)
	default:
		return 0, nil
	}
}

func (query *QueryBuilder) Insert(datas ...map[string]interface{}) (int64, error) {
	query.setData(datas...)
	grammar := Grammar{builder: query}
	sql := grammar.Insert()
	result, err := query.connection.Exec(sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		dblog.Println(err.Error())
		return 0, err
	}
	return result.RowsAffected()
}
func (query *QueryBuilder) InsertGetId(datas map[string]interface{}) (int64, error) {
	query.setData(datas)
	grammar := Grammar{builder: query}
	sql := grammar.Insert()
	result, err := query.connection.Exec(sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		dblog.Println(err.Error())
		return 0, err
	}
	return result.LastInsertId()
}
func (query *QueryBuilder) Replace(datas ...map[string]interface{}) (int64, error) {
	query.setData(datas...)
	grammar := Grammar{builder: query}
	sql := grammar.Replace()
	result, err := query.connection.Exec(sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		dblog.Println(err.Error())
		return 0, err
	}
	return result.RowsAffected()
}
func (query *QueryBuilder) Update(datas map[string]interface{}) (int64, error) {
	query.setData(datas)
	grammar := Grammar{builder: query}
	sql := grammar.Update()
	result, err := query.connection.Exec(sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		dblog.Println(err.Error())
		return 0, err
	}
	return result.RowsAffected()
}

func (query *QueryBuilder) InsertUpdate(insert map[string]interface{}, update map[string]interface{}) (int64, error) {
	query.setData(insert, update)
	grammar := Grammar{builder: query}
	sql := grammar.InsertUpdate()
	result, err := query.connection.Exec(sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		dblog.Println(err.Error())
		return 0, err
	}
	return result.RowsAffected()
}
func (query *QueryBuilder) Delete() (int64, error) {
	grammar := Grammar{builder: query}
	sql := grammar.Delete()
	result, err := query.connection.Exec(sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		dblog.Println(err.Error())
		return 0, err
	}
	return result.RowsAffected()
}
