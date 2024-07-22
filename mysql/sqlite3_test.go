/**
 * Created by zhangruizhi on 2024/7/22
 */

package mysql

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"slices"
	"testing"
)

const supposedSqlite3Result = "[[1,\"sqlite3\"]]"

// TestSqlite3AnonymousQueryScan
// sqlite3 given correct anonymous query scan result
// sql rows: [[1,"sqlite3"]]
// expect marshaled result: [[1,"sqlite3"]]
func TestSqlite3AnonymousQueryScan(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test_sqlite3.db")
	if err != nil {
		t.Errorf("error opening sqlite3 db: %v", err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS t1")
	if err != nil {
		t.Errorf("error dropping table: %v", err)
	}
	_, err = db.Exec("CREATE TABLE t1 (id int, name string)")
	if err != nil {
		t.Errorf("error creating table: %v", err)
	}
	_, err = db.Exec("INSERT INTO t1 (id, name) VALUES (1, 'sqlite3');")
	if err != nil {
		t.Errorf("error inserting into database: %v", err)
	}
	rows, err := db.Query("select * from t1")
	if err != nil {
		t.Errorf("select * from t1: %v", err)
	}
	cols, err := rows.Columns()
	if err != nil {
		t.Errorf("get row columns failed: %v", err)
	}
	values := make([]*interface{}, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var allValues [][]*interface{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			t.Errorf("error scanning row: %v", err)
		}
		allValues = append(allValues, slices.Clone(values))
	}
	bytes, err := json.Marshal(allValues)
	if err != nil {
		return
	}
	if string(bytes) != supposedSqlite3Result {
		t.Failed()
	}
}
