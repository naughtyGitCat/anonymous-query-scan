/**
 * Created by zhangruizhi on 2024/7/22
 */

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"reflect"
	"strconv"
	"time"
)

const (
	dateFormat      = "2006-01-02"
	dateTimeFormat1 = "2006-01-02 15:04:05"
	dateTimeFormat2 = "2006-01-02T15:04:05Z"
)

// grafanaMySQLTypeConverters
// copy from github.com/grafana/grafana/pkg/tsdb/mysql/mysql.go
// add timestamp converter and change time related to local time
// ref: https://dev.mysql.com/doc/refman/8.4/en/data-types.html
var grafanaMySQLTypeConverters = []sqlutil.StringConverter{
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
					return (&v).Local().Unix(), nil
				}
				v, err = time.Parse(dateTimeFormat2, *in)
				if err == nil {
					return (&v).Local().Unix(), nil
				}
				return nil, err
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
					return (&v).Local(), nil
				}
				v, err = time.Parse(dateTimeFormat2, *in)
				if err == nil {
					return (&v).Local(), nil
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

// simplySQLTypeConverter a simplified simplySQLTypeConverter
type simplySQLTypeConverter struct {
	Name        string
	InputType   reflect.Type
	ReplaceFunc func(*string) (any, error)
}

// SimplySQLTypeConverters
// simplified grafana mysql converters
var SimplySQLTypeConverters = []simplySQLTypeConverter{
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
