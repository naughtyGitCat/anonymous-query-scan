/**
 * Created by zhangruizhi on 2024/7/22
 */

package mysql

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
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

// TestScanAnonymousRows
// query row: [1, 'mysql']
// expect result: [[1,"mysql"]]
// actually result: [[1,"mysql"]]
func TestScanAnonymousRows(t *testing.T) {
	err := prepare()
	if err != nil {
		t.Errorf("prepare() failed: %v", err)
	}
	rows, err := query()
	if err != nil {
		t.Errorf("query() failed: %s", err)
	}
	anonymousRows, err := ScanAnonymousRows(rows)
	if err != nil {
		t.Errorf("ScanAnonymousRows() failed: %v", err)
	}
	bytes, err := json.Marshal(anonymousRows)
	if err != nil {
		t.Errorf("json.Marshal() failed: %v", err)
	}
	const rawJson = "[[1,\"mysql\"]]"
	if string(bytes) != rawJson {
		t.Errorf("mapped json doesn't match")
	}
}

// TestScanAnonymousMappedRows
// query row: [1, 'mysql']
// expect result: [{"id":1,"name":"mysql"}]
// actually result: [{"id":1,"name":"mysql"}]
func TestScanAnonymousMappedRows(t *testing.T) {
	err := prepare()
	if err != nil {
		t.Errorf("prepare() failed: %v", err)
	}
	rows, err := query()
	if err != nil {
		t.Errorf("query() failed: %s", err)
	}
	mappedRows, err := ScanAnonymousMappedRows(rows)
	if err != nil {
		t.Errorf("ScanAnonymousMappedRows() failed: %v", err)
	}

	bytes, err := json.Marshal(mappedRows)
	if err != nil {
		t.Errorf("json.Marshal() failed: %v", err)
	}
	const mappedJson = "[{\"id\":1,\"name\":\"mysql\"}]"
	if string(bytes) != mappedJson {
		t.Errorf("mapped json doesn't match")
	}
}

// TestMySQLScanToGrafanaFrames
// query row: [1, 'mysql']
// expect result: {"schema":{"fields":[{"name":"id","type":"number","typeInfo":{"frame":"int64","nullable":true}},{"name":"name","type":"string","typeInfo":{"frame":"string","nullable":true}}]},"data":{"values":[[1],["mysql"]]}}
// actually result: {"schema":{"fields":[{"name":"id","type":"number","typeInfo":{"frame":"int64","nullable":true}},{"name":"name","type":"string","typeInfo":{"frame":"string","nullable":true}}]},"data":{"values":[[1],["mysql"]]}}
func TestMySQLScanToGrafanaFrames(t *testing.T) {
	err := prepare()
	if err != nil {
		t.Errorf("prepare() failed: %v", err)
	}
	rows, err := query()
	if err != nil {
		t.Errorf("query() failed: %v", err)
	}
	frames, err := sqlutil.FrameFromRows(rows, 8888, sqlutil.ToConverters(grafanaMySQLTypeConverters...)...)
	if err != nil {
		t.Errorf("sqlutil.FrameFromRows() failed: %v", err)
	}
	bytes, err := json.Marshal(frames)
	if err != nil {
		t.Errorf("json.Marshal() failed: %v", err)
	}
	const grafanaTypedJson = "{\"schema\":{\"fields\":[{\"name\":\"id\",\"type\":\"number\",\"typeInfo\":{\"frame\":\"int64\",\"nullable\":true}},{\"name\":\"name\",\"type\":\"string\",\"typeInfo\":{\"frame\":\"string\",\"nullable\":true}}]},\"data\":{\"values\":[[1],[\"mysql\"]]}}"
	if string(bytes) != grafanaTypedJson {
		t.Errorf("grafana formatted json doesn't match")
	}
}
