package goqrm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

type Modeler interface {
	Select(...string) *Model
	Get() []map[string]interface{}
	Where(column, operator, value string) *Model
	OrderBy(column, order string) *Model
	Limit(num int) *Model
}

func Connect(connString string) (*sql.DB, error) {

	db, err := sql.Open("mysql", connString)

	Conn = db
	return db, err
}

func NewModel(modeler Modeler) *Model {
	val := reflect.Indirect(reflect.ValueOf(modeler))
	tableName := val.FieldByName("Table")
	return &Model{
		db:             Conn,
		table:          tableName.String(),
		whereCondition: []string{"Where 1 = 1 "},
	}
}

type Model struct {
	db             *sql.DB
	table          string
	columns        []string
	whereCondition []string
	orderBy        string
	limit          string
}

func (m *Model) Select(str ...string) *Model {
	m.columns = str
	return m
}

func (m *Model) Where(column, operator, value string) *Model {
	m.whereCondition = append(m.whereCondition, fmt.Sprint("and ", column, operator, "'", value, "'"))
	return m
}
func (m *Model) OrderBy(column, order string) *Model {
	m.orderBy = fmt.Sprint(" Order by ", column, order)
	return m
}
func (m *Model) Limit(num int) *Model {
	m.limit = fmt.Sprint(" limit ", num)
	return m
}
func (m *Model) ToSql() string {
	columnsString := " * "
	if len(m.columns) > 0 {
		columnsString = strings.Join(m.columns, " , ")
	}

	sql := fmt.Sprint("Select ", columnsString, " from ", m.table, " ", strings.Join(m.whereCondition, " "), " ", m.orderBy, " ", m.limit)

	return sql
}
func (m *Model) Get() []map[string]interface{} {

	rows, err := m.db.Query(m.ToSql())

	if err != nil {
		fmt.Println("Query Error")
		panic(err)
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	// below code is taken from stackoverflow post
	//https://stackoverflow.com/questions/53435455/handling-dynamic-queries-cant-scan-into-struct
	var allMaps []map[string]interface{}
	for rows.Next() {
		values := make([]string, len(columns))
		pointers := make([]interface{}, len(columns))
		for i, _ := range values {
			pointers[i] = &values[i]
		}
		rows.Scan(pointers...)
		resultMap := make(map[string]interface{})
		for i, val := range values {
			resultMap[columns[i]] = val
		}
		allMaps = append(allMaps, resultMap)
	}
	return allMaps
}
