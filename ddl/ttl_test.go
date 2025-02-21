// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ddl

import (
	"testing"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/stretchr/testify/assert"
)

func Test_getTTLInfoInOptions(t *testing.T) {
	falseValue := false
	trueValue := true

	cases := []struct {
		options   []*ast.TableOption
		ttlInfo   *model.TTLInfo
		ttlEnable *bool
		err       error
	}{
		{
			[]*ast.TableOption{},
			nil,
			nil,
			nil,
		},
		{
			[]*ast.TableOption{
				{
					Tp:            ast.TableOptionTTL,
					ColumnName:    &ast.ColumnName{Name: model.NewCIStr("test_column")},
					Value:         ast.NewValueExpr(5, "", ""),
					TimeUnitValue: &ast.TimeUnitExpr{Unit: ast.TimeUnitYear},
				},
			},
			&model.TTLInfo{
				ColumnName:       model.NewCIStr("test_column"),
				IntervalExprStr:  "5",
				IntervalTimeUnit: int(ast.TimeUnitYear),
				Enable:           true,
			},
			nil,
			nil,
		},
		{
			[]*ast.TableOption{
				{
					Tp:        ast.TableOptionTTLEnable,
					BoolValue: false,
				},
				{
					Tp:            ast.TableOptionTTL,
					ColumnName:    &ast.ColumnName{Name: model.NewCIStr("test_column")},
					Value:         ast.NewValueExpr(5, "", ""),
					TimeUnitValue: &ast.TimeUnitExpr{Unit: ast.TimeUnitYear},
				},
			},
			&model.TTLInfo{
				ColumnName:       model.NewCIStr("test_column"),
				IntervalExprStr:  "5",
				IntervalTimeUnit: int(ast.TimeUnitYear),
				Enable:           false,
			},
			&falseValue,
			nil,
		},
		{
			[]*ast.TableOption{
				{
					Tp:        ast.TableOptionTTLEnable,
					BoolValue: false,
				},
				{
					Tp:            ast.TableOptionTTL,
					ColumnName:    &ast.ColumnName{Name: model.NewCIStr("test_column")},
					Value:         ast.NewValueExpr(5, "", ""),
					TimeUnitValue: &ast.TimeUnitExpr{Unit: ast.TimeUnitYear},
				},
				{
					Tp:        ast.TableOptionTTLEnable,
					BoolValue: true,
				},
			},
			&model.TTLInfo{
				ColumnName:       model.NewCIStr("test_column"),
				IntervalExprStr:  "5",
				IntervalTimeUnit: int(ast.TimeUnitYear),
				Enable:           true,
			},
			&trueValue,
			nil,
		},
	}

	for _, c := range cases {
		ttlInfo, ttlEnable, err := getTTLInfoInOptions(c.options)

		assert.Equal(t, c.ttlInfo, ttlInfo)
		assert.Equal(t, c.ttlEnable, ttlEnable)
		assert.Equal(t, c.err, err)
	}
}
