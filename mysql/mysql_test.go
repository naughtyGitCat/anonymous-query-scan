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
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS `t_json`")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE t_json(id int, content json)")
	return err
}

func insertData() error {
	db, err := sql.Open("mysql", mysqlConnectionStr)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO t1(id, name) VALUES (1, 'mysql')")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO t_json(id, content) VALUES (1, '{\"name\":\"mysql\", \"version\": 5.7, \"enabled\": true}'), (2, '[\"5.7\", \"8.0\"]'), (3, '[{\"name\":\"mysql\", \"version\": 5.7}]')")
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

func queryT1() (*sql.Rows, error) {
	db, err := newConnection()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select * from t1")
	return rows, err
}

func queryTJSON() (*sql.Rows, error) {
	db, err := newConnection()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select * from t_json")
	return rows, err
}

// TestMySQLAnonymousScan
// queryT1 row: [1, 'mysql']
// expect result: [[1, 'mysql']]
// actually result: [[1,"bXlzcWw="]]
func TestMySQLAnonymousScan(t *testing.T) {
	err := prepare()
	if err != nil {
		t.Errorf("prepare() failed: %v", err)
	}
	rows, err := queryT1()
	if err != nil {
		t.Errorf("queryT1() failed: %v", err)
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
// queryT1 row: [1, 'mysql']
// expect result: [[1,"mysql"]]
// actually result: [[1,"mysql"]]
func TestScanAnonymousRows(t *testing.T) {
	err := prepare()
	if err != nil {
		t.Errorf("prepare() failed: %v", err)
	}
	rows, err := queryT1()
	if err != nil {
		t.Errorf("queryT1() failed: %s", err)
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
// queryT1 row: [1, 'mysql']
// expect result: [{"id":1,"name":"mysql"}]
// actually result: [{"id":1,"name":"mysql"}]
func TestScanAnonymousMappedRows(t *testing.T) {
	err := prepare()
	if err != nil {
		t.Errorf("prepare() failed: %v", err)
	}
	rows, err := queryT1()
	if err != nil {
		t.Errorf("queryT1() failed: %s", err)
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

// TestMySQLJsonColumn
// queryTJSON rows: (1, '{"name":"mysql", "version": 5.7, "enabled": true}'), (2, '["5.7", "8.0"]'), (3, '[{"name":"mysql", "version": 5.7}]')
// expect result: [{"content":{"enabled":true,"name":"mysql","version":5.7},"id":1},{"content":["5.7","8.0"],"id":2},{"content":[{"name":"mysql","version":5.7}],"id":3}]
// actually result: [{"content":{"enabled":true,"name":"mysql","version":5.7},"id":1},{"content":["5.7","8.0"],"id":2},{"content":[{"name":"mysql","version":5.7}],"id":3}]
func TestMySQLJsonColumn(t *testing.T) {
	err := prepare()
	if err != nil {
		t.Errorf("prepare() failed: %v", err)
	}
	rows, err := queryTJSON()
	if err != nil {
		t.Errorf("queryTJSON() failed: %v", err)
	}
	mappedRows, err := ScanAnonymousMappedRowsExt(rows)
	if err != nil {
		t.Errorf("ScanAnonymousMappedRowsExt() failed: %v", err)
	}
	bytes, err := json.Marshal(mappedRows)
	if err != nil {
		t.Errorf("json.Marshal() failed: %v", err)
	}
	const rawJson = "[{\"content\":{\"enabled\":true,\"name\":\"mysql\",\"version\":5.7},\"id\":1},{\"content\":[\"5.7\",\"8.0\"],\"id\":2},{\"content\":[{\"name\":\"mysql\",\"version\":5.7}],\"id\":3}]"
	if string(bytes) != rawJson {
		t.Errorf("mapped json doesn't match")
	}
}

// TestMySQLScanToGrafanaFrames
// queryT1 row: [1, 'mysql']
// expect result: {"schema":{"fields":[{"name":"id","type":"number","typeInfo":{"frame":"int64","nullable":true}},{"name":"name","type":"string","typeInfo":{"frame":"string","nullable":true}}]},"data":{"values":[[1],["mysql"]]}}
// actually result: {"schema":{"fields":[{"name":"id","type":"number","typeInfo":{"frame":"int64","nullable":true}},{"name":"name","type":"string","typeInfo":{"frame":"string","nullable":true}}]},"data":{"values":[[1],["mysql"]]}}

//func TestMySQLScanToGrafanaFrames(t *testing.T) {
//	err := prepare()
//	if err != nil {
//		t.Errorf("prepare() failed: %v", err)
//	}
//	rows, err := queryT1()
//	if err != nil {
//		t.Errorf("queryT1() failed: %v", err)
//	}
//
//	frames, err := sqlutil.FrameFromRows(rows, 8888, sqlutil.ToConverters(grafanaMySQLTypeConverters...)...)
//	if err != nil {
//		t.Errorf("sqlutil.FrameFromRows() failed: %v", err)
//	}
//	bytes, err := json.Marshal(frames)
//	if err != nil {
//		t.Errorf("json.Marshal() failed: %v", err)
//	}
//	const grafanaTypedJson = "{\"schema\":{\"fields\":[{\"name\":\"id\",\"type\":\"number\",\"typeInfo\":{\"frame\":\"int64\",\"nullable\":true}},{\"name\":\"name\",\"type\":\"string\",\"typeInfo\":{\"frame\":\"string\",\"nullable\":true}}]},\"data\":{\"values\":[[1],[\"mysql\"]]}}"
//	if string(bytes) != grafanaTypedJson {
//		t.Errorf("grafana formatted json doesn't match")
//	}
//}
