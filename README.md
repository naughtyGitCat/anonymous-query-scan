# anonymous-query-scan

> `go-sql-driver/sql` just can't do marshal anonymous query scan correctly

## How-to

```go
// ScanAnonymousMappedRows scan anonymous rows without predefined struct, using simply converter match with sql types
// cols type related with time, datetime, date, timestamp will be converted to utc time, timestamp will be datetime
// return format likes: [{'col1': number, 'col2': 'string', 'col3': '0000-00-00T00:00:00Z',...}...]
func ScanAnonymousMappedRows(rows *sql.Rows) ([]map[string]any, error) {}

// ScanAnonymousRows scan anonymous rows without predefined struct, using simply converter match with sql types
// cols type related with time, datetime, date, timestamp will be converted to utc time, timestamp will be datetime
// return format likes: [{number, 'string', '0000-00-00T00:00:00Z',...}...]
func ScanAnonymousRows(rows *sql.Rows) ([][]any, error) {}

// ScanAnonymousMappedRowsExt scan anonymous rows without predefined struct, using grafana converter match with mysql types
// cols type related with time, datetime, date, timestamp will be converted to local time, timestamp will be number
// return format likes: [{number, 'string', '0000-00-00T00:00:00±0:00',...}...]
func ScanAnonymousMappedRowsExt(rows *sql.Rows) ([]map[string]any, error) {}
```

#### install
```shell
go get -u https://github.com/naughtyGitCat/anonymous-query-scan/mysql
```

#### usage
```go
import (
    "encoding/json"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    mysql "github.com/naughtyGitCat/anonymous-query-scan/mysql"
)
func main() {
    db, err := sql.Open("mysql", mysqlConnectionStr)
	if err != nil {
		panic(err)
    }
    rows, err := db.Query("select * from user")
    if err != nil {
        panic(err) 
    // ScanAnonymousMappedRowsExt return timestamp optimized and time localized rows 
    mappedRows, err := mysql.ScanAnonymousMappedRowsExt(rows)
    if err != nil {
        return nil, err
    }
    rowBytes, err := json.Marshal(mappedRows)
    if err != nil {
        return nil, err
    }
	fmt.Println(string(bytes))
}
```

## ref
- https://dev.mysql.com/doc/refman/8.4/en/data-types.html
- https://stackoverflow.com/questions/55995917/reverse-of-reflect-typeof
- https://github.com/grafana/grafana-plugin-sdk-go
- https://github.com/go-sql-driver/mysql/wiki/Examples
- [ISSUE-407 Returned values are always []byte](https://github.com/go-sql-driver/mysql/issues/407)