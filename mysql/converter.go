/**
 * Created by zhangruizhi on 2024/7/22
 */

package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	dateFormat      = "2006-01-02"
	dateTimeFormat1 = "2006-01-02 15:04:05"
	dateTimeFormat2 = "2006-01-02T15:04:05Z"
)

type mysqlTypeConverter struct {
	Name        string
	MySQLType   string
	ScanType    reflect.Type
	ReplaceFunc func(*string) (any, error)
}

// grafanaMySQLTypeConverters
// copy from github.com/grafana/grafana/pkg/tsdb/mysql/mysql.go
// add timestamp converter and change time related to local time
// ref: https://dev.mysql.com/doc/refman/8.4/en/data-types.html
var mysqlTypeConverters = []mysqlTypeConverter{
	{
		Name:      "handle DOUBLE",
		ScanType:  reflect.TypeOf(sql.NullFloat64{}),
		MySQLType: "DOUBLE",
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
		Name:      "handle BIGINT",
		ScanType:  reflect.TypeOf(sql.NullInt64{}),
		MySQLType: "BIGINT",
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
		Name:      "handle DECIMAL",
		ScanType:  reflect.TypeOf(sql.NullFloat64{}),
		MySQLType: "DECIMAL",
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
		Name:      "handle TIMESTAMP",
		ScanType:  reflect.TypeOf(sql.NullInt64{}),
		MySQLType: "TIMESTAMP",
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			v, err := time.Parse(dateTimeFormat1, *in)
			if err == nil {
				return (&v).Local().Unix(), nil
			}
			v, err = time.Parse(dateTimeFormat2, *in)
			if err == nil {
				return (&v).Local().Unix(), nil
			}
			return nil, err
		},
	},
	{
		Name:      "handle DATETIME",
		ScanType:  reflect.TypeOf(sql.NullTime{}),
		MySQLType: "DATETIME",
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			v, err := time.Parse(dateTimeFormat1, *in)
			if err == nil {
				return (&v).Local(), nil
			}
			v, err = time.Parse(dateTimeFormat2, *in)
			if err == nil {
				return (&v).Local(), nil
			}

			return nil, err
		},
	},
	{
		Name:      "handle DATE",
		ScanType:  reflect.TypeOf(sql.NullTime{}),
		MySQLType: "DATE",
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			v, err := time.Parse(dateFormat, *in)
			if err == nil {
				return (&v).Local(), nil
			}
			v, err = time.Parse(dateTimeFormat1, *in)
			if err == nil {
				return (&v).Local(), nil
			}
			v, err = time.Parse(dateTimeFormat2, *in)
			if err == nil {
				return (&v).Local(), nil
			}
			return nil, err
		},
	},
	{
		Name:      "handle YEAR",
		ScanType:  reflect.TypeOf(sql.NullInt64{}),
		MySQLType: "YEAR",
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
		Name:      "handle TINYINT",
		ScanType:  reflect.TypeOf(sql.NullInt64{}),
		MySQLType: "TINYINT",
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
		Name:      "handle SMALLINT",
		ScanType:  reflect.TypeOf(sql.NullInt64{}),
		MySQLType: "SMALLINT",
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
		Name:      "handle INT",
		ScanType:  reflect.TypeOf(sql.NullInt64{}),
		MySQLType: "INT",
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
		Name:      "handle FLOAT",
		ScanType:  reflect.TypeOf(sql.NullFloat64{}),
		MySQLType: "FLOAT",
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
		Name:      "handle JSON",
		ScanType:  reflect.TypeOf(sql.NullString{}),
		MySQLType: "JSON",
		ReplaceFunc: func(in *string) (any, error) {
			if in == nil {
				return nil, nil
			}
			var j any
			err := json.Unmarshal([]byte(*in), &j)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
			}
			return &j, nil
		},
	},
}

// goSQLTypeConverter a simplified goSQLTypeConverter
type goSQLTypeConverter struct {
	Name        string
	InputType   reflect.Type
	ReplaceFunc func(*string) (any, error)
}

// SimplySQLTypeConverters
// simplified grafana mysql converters
var SimplySQLTypeConverters = []goSQLTypeConverter{
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
