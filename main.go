/**
 * Created by zhangruizhi on 2024/7/17
 */

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	_ "github.com/mattn/go-sqlite3"

	"reflect"
	"slices"
	"strconv"
	"time"
)

const mysqlConnectionStr = "root:******@tcp(127.0.0.1:3306)/test?charset=utf8"

const (
	dateFormat      = "2006-01-02"
	dateTimeFormat1 = "2006-01-02 15:04:05"
	dateTimeFormat2 = "2006-01-02T15:04:05Z"
)

type SQLTypeConverter struct {
	Name        string
	InputType   reflect.Type
	ReplaceFunc func(*string) (any, error)
}

var sqlTypeConverters = []SQLTypeConverter{
	{
		Name:      "NullTime",
		InputType: reflect.TypeOf(sql.NullTime{}),
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			v, err := time.Parse(dateTimeFormat1, *in)
			if err == nil {
				return &v, nil
			}
			v, err = time.Parse(dateTimeFormat2, *in)
			if err == nil {
				return &v, nil
			}
			return nil, err
		},
	},
	{
		Name:      "NullString",
		InputType: reflect.TypeOf(sql.NullString{}),
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			return &in, nil
		},
	},
	{
		Name:      "NullByte",
		InputType: reflect.TypeOf(sql.NullByte{}),
		ReplaceFunc: func(in *string) (any, error) {
			return nil, errors.New("scan type NullByte is not supported")
		},
	},
	{
		Name:      "NullBool",
		InputType: reflect.TypeOf(sql.NullBool{}),
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			if *in == "1" {
				return true, nil
			} else if *in == "0" {
				return false, nil
			}
			return nil, errors.New(fmt.Sprintf("scan type bool value %s is not supported", *in))
		},
	},
	{
		Name:      "NullFloat64",
		InputType: reflect.TypeOf(sql.NullFloat64{}),
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			v, err := strconv.ParseFloat(*in, 64)
			if err != nil {
				return nil, err
			}
			return &v, nil
		},
	},
	{
		Name:      "NullInt16",
		InputType: reflect.TypeOf(sql.NullInt16{}),
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			v, err := strconv.ParseInt(*in, 10, 64)
			if err != nil {
				return nil, err
			}
			return &v, nil
		},
	},
	{
		Name:      "NullInt32",
		InputType: reflect.TypeOf(sql.NullInt32{}),
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			v, err := strconv.ParseInt(*in, 10, 64)
			if err != nil {
				return nil, err
			}
			return &v, nil
		},
	},
	{
		Name:      "NullInt64",
		InputType: reflect.TypeOf(sql.NullInt64{}),
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			v, err := strconv.ParseInt(*in, 10, 64)
			if err != nil {
				return nil, err
			}
			return &v, nil
		},
	},
}

var mysqlTypeConverters = []sqlutil.StringConverter{
	{
		Name:           "handle DOUBLE",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "DOUBLE",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableFloat64,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := strconv.ParseFloat(*in, 64)
				if err != nil {
					return nil, err
				}
				return &v, nil
			},
		},
	},
	{
		Name:           "handle BIGINT",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "BIGINT",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableInt64,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := strconv.ParseInt(*in, 10, 64)
				if err != nil {
					return nil, err
				}
				return &v, nil
			},
		},
	},
	{
		Name:           "handle DECIMAL",
		InputScanKind:  reflect.Slice,
		InputTypeName:  "DECIMAL",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableFloat64,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := strconv.ParseFloat(*in, 64)
				if err != nil {
					return nil, err
				}
				return &v, nil
			},
		},
	},
	{
		Name:           "handle DATETIME",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "DATETIME",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableTime,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := time.Parse(dateTimeFormat1, *in)
				if err == nil {
					return &v, nil
				}
				v, err = time.Parse(dateTimeFormat2, *in)
				if err == nil {
					return &v, nil
				}

				return nil, err
			},
		},
	},
	{
		Name:           "handle DATE",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "DATE",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableTime,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := time.Parse(dateFormat, *in)
				if err == nil {
					return &v, nil
				}
				v, err = time.Parse(dateTimeFormat1, *in)
				if err == nil {
					return &v, nil
				}
				v, err = time.Parse(dateTimeFormat2, *in)
				if err == nil {
					return &v, nil
				}
				return nil, err
			},
		},
	},
	{
		Name:           "handle TIMESTAMP",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "TIMESTAMP",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableInt64,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := strconv.ParseInt(*in, 10, 64)
				if err == nil {
					return &v, nil
				}
				return nil, err
			},
		},
	},
	{
		Name:           "handle YEAR",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "YEAR",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableInt64,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := strconv.ParseInt(*in, 10, 64)
				if err != nil {
					return nil, err
				}
				return &v, nil
			},
		},
	},
	{
		Name:           "handle TINYINT",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "TINYINT",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableInt64,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := strconv.ParseInt(*in, 10, 64)
				if err != nil {
					return nil, err
				}
				return &v, nil
			},
		},
	},
	{
		Name:           "handle SMALLINT",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "SMALLINT",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableInt64,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := strconv.ParseInt(*in, 10, 64)
				if err != nil {
					return nil, err
				}
				return &v, nil
			},
		},
	},
	{
		Name:           "handle INT",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "INT",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableInt64,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := strconv.ParseInt(*in, 10, 64)
				if err != nil {
					return nil, err
				}
				return &v, nil
			},
		},
	},
	{
		Name:           "handle FLOAT",
		InputScanKind:  reflect.Struct,
		InputTypeName:  "FLOAT",
		ConversionFunc: func(in *string) (*string, error) { return in, nil },
		Replacer: &sqlutil.StringFieldReplacer{
			OutputFieldType: data.FieldTypeNullableFloat64,
			ReplaceFunc: func(in *string) (any, error) {
				if in == nil {
					return nil, nil
				}
				v, err := strconv.ParseFloat(*in, 64)
				if err != nil {
					return nil, err
				}
				return &v, nil
			},
		},
	},
}

func main0() {
	var x = []uint8{1, 2, 3, 4, 5, 6}
	var y = []int8{1, 2, 3, 4, 5, 6}
	xBytes, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	yBytes, err := json.Marshal(y)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("uint8 %s, int8: %s", string(xBytes), string(yBytes)))
}

func main2() {
	db, err := sql.Open("mysql", mysqlConnectionStr)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select * from t3")
	if err != nil {
		panic(err)
	}
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
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
			panic(err)
		}
		fmt.Println(fmt.Sprintf("values: %#v", values))
		allValues = append(allValues, slices.Clone(values))
	}
	fmt.Println(fmt.Sprintf("allValues: %#v", allValues))
	//for i := range allValues {
	//	for j := range allValues[i] {
	//		fmt.Println(fmt.Sprintf("%d,%T", allValues[i][j], allValues[i][j]))
	//	}
	//}
	bytes, err := json.Marshal(allValues)
	if err != nil {
		return
	}
	fmt.Println(string(bytes))
}

func grafanaFrameFromRows() {
	db, err := sql.Open("mysql", mysqlConnectionStr)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select * from t4")
	if err != nil {
		panic(err)
	}
	// var converters = []sqlutil.Converter{sqlutil.NullStringConverter, sqlutil.NullInt64Converter, sqlutil.NullInt32Converter, sqlutil.NullDecimalConverter, sqlutil.NullInt16Converter, sqlutil.NullTimeConverter, sqlutil.NullByteConverter}
	var sConverters = []sqlutil.StringConverter{
		{
			Name:           "handle DOUBLE",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "DOUBLE",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableFloat64,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := strconv.ParseFloat(*in, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:           "handle BIGINT",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "BIGINT",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableInt64,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := strconv.ParseInt(*in, 10, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:           "handle DECIMAL",
			InputScanKind:  reflect.Slice,
			InputTypeName:  "DECIMAL",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableFloat64,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := strconv.ParseFloat(*in, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:           "handle DATETIME",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "DATETIME",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableTime,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := time.Parse(dateTimeFormat1, *in)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dateTimeFormat2, *in)
					if err == nil {
						return &v, nil
					}

					return nil, err
				},
			},
		},
		{
			Name:           "handle DATE",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "DATE",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableTime,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := time.Parse(dateFormat, *in)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dateTimeFormat1, *in)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dateTimeFormat2, *in)
					if err == nil {
						return &v, nil
					}
					return nil, err
				},
			},
		},
		{
			Name:           "handle TIMESTAMP",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "TIMESTAMP",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableTime,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := time.Parse(dateTimeFormat1, *in)
					if err == nil {
						return &v, nil
					}
					v, err = time.Parse(dateTimeFormat2, *in)
					if err == nil {
						return &v, nil
					}
					return nil, err
				},
			},
		},
		{
			Name:           "handle YEAR",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "YEAR",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableInt64,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := strconv.ParseInt(*in, 10, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:           "handle TINYINT",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "TINYINT",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableInt64,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := strconv.ParseInt(*in, 10, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:           "handle SMALLINT",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "SMALLINT",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableInt64,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := strconv.ParseInt(*in, 10, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:           "handle INT",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "INT",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableInt64,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := strconv.ParseInt(*in, 10, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
		{
			Name:           "handle FLOAT",
			InputScanKind:  reflect.Struct,
			InputTypeName:  "FLOAT",
			ConversionFunc: func(in *string) (*string, error) { return in, nil },
			Replacer: &sqlutil.StringFieldReplacer{
				OutputFieldType: data.FieldTypeNullableFloat64,
				ReplaceFunc: func(in *string) (any, error) {
					if in == nil {
						return nil, nil
					}
					v, err := strconv.ParseFloat(*in, 64)
					if err != nil {
						return nil, err
					}
					return &v, nil
				},
			},
		},
	}
	fromRows, err := sqlutil.FrameFromRows(rows, 8888, sqlutil.ToConverters(sConverters...)...)
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(fromRows)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

// typeConvert5StringWithColName do convert when scan type match
func typeConvert5StringWithColName() {
	db, err := sql.Open("mysql", mysqlConnectionStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from t4")
	if err != nil {
		panic(err)
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		panic(err)
	}

	var colNames []string
	for _, colType := range colTypes {
		if slices.Contains(colNames, colType.Name()) {
			panic("duplicate col name: " + colType.Name())
		}
		colNames = append(colNames, colType.Name())
	}

	var colScanTypes []reflect.Type
	for _, colType := range colTypes {
		fmt.Println(fmt.Sprintf("type %s %s %s", colType.Name(), colType.DatabaseTypeName(), colType.ScanType()))
		colScanTypes = append(colScanTypes, colType.ScanType())
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
			panic(err)
		}
		var mappedRow = map[string]any{}
	ValueLoop:
		for i, stringV := range values {
			for _, converter := range sqlTypeConverters {
				if colScanTypes[i] == converter.InputType {
					convertedValue, err := converter.ReplaceFunc(stringV)
					if err != nil {
						panic(err)
					}
					mappedRow[colNames[i]] = convertedValue
					continue ValueLoop
				}
			}
			mappedRow[colNames[i]] = stringV
		}
		allRows = append(allRows, mappedRow)
	}

	bytes, err := json.Marshal(allRows)
	if err != nil {
		return
	}
	fmt.Println(fmt.Sprintf("allValues: %#v, json: %s", allRows, string(bytes)))
}

func typeConvert5String() {
	db, err := sql.Open("mysql", mysqlConnectionStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from t4")
	if err != nil {
		panic(err)
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		panic(err)
	}
	var types []reflect.Type
	for _, colType := range colTypes {
		types = append(types, colType.ScanType())
	}
	for _, colType := range colTypes {
		fmt.Println(fmt.Sprintf("type %s %s %s", colType.Name(), colType.DatabaseTypeName(), colType.ScanType()))
	}
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	scanArgs := make([]interface{}, len(cols))
	values := make([]*string, len(cols))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var allValues [][]interface{}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}
		typedValues := make([]interface{}, len(cols))
	ValueLoop:
		for i, stringV := range values {
			fmt.Println(fmt.Sprintf("stringV: %v", &stringV))
			for _, converter := range sqlTypeConverters {
				if types[i] == converter.InputType {
					convertedValue, err := converter.ReplaceFunc(stringV)
					if err != nil {
						panic(err)
					}
					typedValues[i] = convertedValue
					continue ValueLoop
				}
			}
			typedValues[i] = stringV
		}
		allValues = append(allValues, typedValues)
		fmt.Println(fmt.Sprintf("values: %#v", typedValues))
		// allValues = append(allValues, slices.Clone(values))
	}
	//for i := range allValues {
	//	for j := range allValues[i] {
	//		fmt.Println(fmt.Sprintf("%d,%T", allValues[i][j], allValues[i][j]))
	//	}
	//}
	bytes, err := json.Marshal(allValues)
	if err != nil {
		return
	}
	fmt.Println(fmt.Sprintf("allValues: %#v, json: %s", allValues, string(bytes)))
}

// typeConvert5StringWithColName1 do convert when database type match
func typeConvert5StringWithColName1() {
	db, err := sql.Open("mysql", mysqlConnectionStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from t4")
	if err != nil {
		panic(err)
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		panic(err)
	}

	var colNames []string
	for _, colType := range colTypes {
		if slices.Contains(colNames, colType.Name()) {
			panic("duplicate col name: " + colType.Name())
		}
		colNames = append(colNames, colType.Name())
	}

	var colMySQLTypes []string
	for _, colType := range colTypes {
		fmt.Println(fmt.Sprintf("type %s %s %s", colType.Name(), colType.DatabaseTypeName(), colType.ScanType()))
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
			panic(err)
		}
		var mappedRow = map[string]any{}
	ValueLoop:
		for i, stringV := range values {
			for _, converter := range mysqlTypeConverters {
				if colMySQLTypes[i] == converter.InputTypeName {
					convertedValue, err := converter.Replacer.ReplaceFunc(stringV)
					if err != nil {
						panic(err)
					}
					mappedRow[colNames[i]] = convertedValue
					continue ValueLoop
				}
			}
			mappedRow[colNames[i]] = stringV
		}
		allRows = append(allRows, mappedRow)
	}

	bytes, err := json.Marshal(allRows)
	if err != nil {
		return
	}
	fmt.Println(fmt.Sprintf("allValues: %#v, json: %s", allRows, string(bytes)))
}

func main1() {
	db, err := sql.Open("sqlite3", "./awesomeProject.db")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select * from t1")
	if err != nil {
		panic(err)
	}
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
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
			panic(err)
		}
		allValues = append(allValues, slices.Clone(values))
	}
	bytes, err := json.Marshal(allValues)
	if err != nil {
		return
	}
	fmt.Println(string(bytes))
}

func main() {
	// grafanaFrameFromRows()
	// typeConvert5StringWithColName()
	typeConvert5StringWithColName1()
}
