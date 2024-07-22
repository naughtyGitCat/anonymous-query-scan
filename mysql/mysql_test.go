/**
 * Created by zhangruizhi on 2024/7/22
 */

package mysql

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"slices"
	"testing"
)

const mysqlConnectionStr = "root:******@tcp(127.0.0.1:3306)/test?charset=utf8"

func newConnection() (*sql.DB, error) {
	return sql.Open("mysql", mysqlConnectionStr)
}

func createTable() error {
	db, err := newConnection()
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS `t1`")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE t1(id int, name varchar(18))")
	return err
}

func insertData() error {
	db, err := sql.Open("mysql", mysqlConnectionStr)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO t1(id, name) VALUES (1, 'mysql')")
	return err
}

func prepare() error {
	err := createTable()
	if err != nil {
		return err
	}
	err = insertData()
	return err
}

func query() (*sql.Rows, error) {
	db, err := newConnection()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select * from t1")
	return rows, err
}

// TestMySQLAnonymousScan
// query row: [1, 'mysql']
// expect result: [[1, 'mysql']]
// actually result: [[1,"bXlzcWw="]]
func TestMySQLAnonymousScan(t *testing.T) {
	err := prepare()
	if err != nil {
		t.Errorf("prepare() failed: %v", err)
	}
	rows, err := query()
	if err != nil {
		t.Errorf("query() failed: %v", err)
	}
	cols, err := rows.Columns()
	if err != nil {
		t.Errorf("rows.Columns() failed: %v", err)
	}
	values := make([]interface{}, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var allValues [][]interface{}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			t.Errorf("rows.Scan() failed: %v", err)
		}
		allValues = append(allValues, slices.Clone(values))
	}
	bytes, err := json.Marshal(allValues)
	if err != nil {
		t.Errorf("json.Marshal() failed: %v", err)
	}
	t.Errorf(string(bytes))
}
