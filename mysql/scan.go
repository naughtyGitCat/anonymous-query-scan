/**
 * Created by zhangruizhi on 2024/7/17
 */

package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"slices"
)

// ScanAnonymousMappedRows scan anonymous rows without predefined struct, using simply converter match with sql types
// cols type related with time, datetime, date, timestamp will be converted to utc time, timestamp will be datetime
// return format likes: [{'col1': number, 'col2': 'string', 'col3': '0000-00-00T00:00:00Z',...}...]
func ScanAnonymousMappedRows(rows *sql.Rows) ([]map[string]any, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("get colTypes failed, %w", err)
	}
	// get col names
	var colNames []string
	for _, colType := range colTypes {
		if slices.Contains(colNames, colType.Name()) {
			return nil, fmt.Errorf("duplicate column name %s", colType.Name())
		}
		colNames = append(colNames, colType.Name())
	}
	// get col types
	var colScanTypes []reflect.Type
	for _, colType := range colTypes {
		colScanTypes = append(colScanTypes, colType.ScanType())
	}
	// prepare for scan
	scanArgs := make([]interface{}, len(colTypes))
	values := make([]*string, len(colTypes))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var allRows []map[string]any

	// native scan rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}
		var mappedRow = map[string]any{}
	ValueLoop:
		for i, stringV := range values {
			for _, converter := range SimplySQLTypeConverters {
				if colScanTypes[i] == converter.InputType {
					convertedValue, err := converter.ReplaceFunc(stringV)
					if err != nil {
						return nil, fmt.Errorf("convert value failed, %w", err)
					}
					mappedRow[colNames[i]] = convertedValue
					continue ValueLoop
				}
			}
			mappedRow[colNames[i]] = stringV
		}
		allRows = append(allRows, mappedRow)
	}
	return allRows, nil
}

// ScanAnonymousRows scan anonymous rows without predefined struct, using simply converter match with sql types
// cols type related with time, datetime, date, timestamp will be converted to utc time, timestamp will be datetime
// return format likes: [{number, 'string', '0000-00-00T00:00:00Z',...}...]
func ScanAnonymousRows(rows *sql.Rows) ([][]any, error) {
	// get col types for type convert
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("get colTypes failed, %w", err)
	}
	var colScanTypes []reflect.Type
	for _, colType := range colTypes {
		colScanTypes = append(colScanTypes, colType.ScanType())
	}
	// prepare for scan
	scanArgs := make([]interface{}, len(colTypes))
	values := make([]*string, len(colTypes))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var allValues [][]any

	// native scan rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, fmt.Errorf("scan row failed, %w", err)
		}
		typedValues := make([]interface{}, len(colTypes))
	ValueLoop:
		for i, stringV := range values {
			for _, converter := range SimplySQLTypeConverters {
				if colScanTypes[i] == converter.InputType {
					convertedValue, err := converter.ReplaceFunc(stringV)
					if err != nil {
						return nil, fmt.Errorf("convert value failed, %w", err)
					}
					typedValues[i] = convertedValue
					continue ValueLoop
				}
			}
			typedValues[i] = stringV
		}
		allValues = append(allValues, typedValues)
	}
	return allValues, nil
}

// ScanAnonymousMappedRowsExt scan anonymous rows without predefined struct, using grafana converter match with mysql types
// cols type related with time, datetime, date, timestamp will be converted to local time, timestamp will be number
// return format likes: [{number, 'string', '0000-00-00T00:00:00±0:00',...}...]
func ScanAnonymousMappedRowsExt(rows *sql.Rows) ([]map[string]any, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("get colTypes failed, %w", err)
	}

	var colNames []string
	for _, colType := range colTypes {
		if slices.Contains(colNames, colType.Name()) {
			return nil, fmt.Errorf("duplicate column name %s", colType.Name())
		}
		colNames = append(colNames, colType.Name())
	}

	var colMySQLTypes []string
	for _, colType := range colTypes {
		// useful debug line:
		// fmt.Println(fmt.Sprintf("type %s %s %s", colType.Name(), colType.DatabaseTypeName(), colType.ScanType()))
		colMySQLTypes = append(colMySQLTypes, colType.DatabaseTypeName())
	}

	scanArgs := make([]interface{}, len(colTypes))
	values := make([]*string, len(colTypes))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var allRows []map[string]any

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, fmt.Errorf("scan row failed, %w", err)
		}
		var mappedRow = map[string]any{}
	ValueLoop:
		for i, stringV := range values {
			for _, converter := range mysqlTypeConverters {
				if colMySQLTypes[i] == converter.MySQLType {
					convertedValue, err := converter.ReplaceFunc(stringV)
					if err != nil {
						return nil, fmt.Errorf("convert value failed, %w", err)
					}
					mappedRow[colNames[i]] = convertedValue
					continue ValueLoop
				}
			}
			mappedRow[colNames[i]] = stringV
		}
		allRows = append(allRows, mappedRow)
	}
	return allRows, nil
}
