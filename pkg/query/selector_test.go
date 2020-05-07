// Copyright 2020 Zhizhesihai (Beijing) Technology Limited.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package query

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/query/ast"
)

func TestSelector_QuerySQL(t *testing.T) {
	var metas = map[string]model.ObjectMeta{}
	metas["B"] = model.ObjectMeta{
		ID:          0,
		TypeID:      1,
		Name:        "B",
		ValueType:   model.BOOLEAN,
		Description: "",
		CreateTime:  time.Time{},
		DeleteTime:  nil,
	}
	metas["S"] = model.ObjectMeta{
		ID:          0,
		TypeID:      1,
		Name:        "S",
		ValueType:   model.STRING,
		Description: "",
		CreateTime:  time.Time{},
		DeleteTime:  nil,
	}
	metas["D"] = model.ObjectMeta{
		ID:          0,
		TypeID:      1,
		Name:        "D",
		ValueType:   model.DOUBLE,
		Description: "",
		CreateTime:  time.Time{},
		DeleteTime:  nil,
	}
	metas["I"] = model.ObjectMeta{
		ID:          0,
		TypeID:      1,
		Name:        "I",
		ValueType:   model.INTEGER,
		Description: "",
		CreateTime:  time.Time{},
		DeleteTime:  nil,
	}
	var requires []*ast.Requirement
	for _, meta := range metas {
		operators := validOperator[meta.ValueType]
		for _, operator := range operators {
			requires = append(requires, &ast.Requirement{
				Key:      meta.Name,
				Operator: operator,
				Value:    []string{randStr(meta.ValueType), randStr(meta.ValueType)},
			})
		}
	}
	selector := &Selector{requires: requires}
	sql, args, err := selector.QuerySQL(metas)
	t.Log(sql, args, err)
}

func randStr(typ int) string {
	var v string
	switch typ {
	case model.BOOLEAN:
		v = "true"
	case model.INTEGER:
		v = strconv.Itoa(rand.Intn(100))
	case model.DOUBLE:
		v = strconv.FormatFloat(rand.Float64()+float64(rand.Intn(10)), 'g', -1, 64)
	case model.STRING:
		v = String(3)
	}
	return v
}

var validOperator = map[int][]ast.Operator{
	model.STRING: {
		ast.In,
		ast.NotIn,
		ast.Equals,
		ast.NotEquals,
		ast.LessThan,
		ast.GreaterThan,
		ast.LessThanOrEquals,
		ast.GreaterThanOrEquals,
		ast.Exists,
		ast.DoesNotExist,
		ast.Positive,
		ast.Negates,
	},
	model.BOOLEAN: {
		ast.Equals,
		ast.Exists,
		ast.DoesNotExist,
		ast.Positive,
		ast.Negates,
	},
	model.INTEGER: {
		ast.In,
		ast.NotIn,
		ast.Equals,
		ast.NotEquals,
		ast.LessThan,
		ast.GreaterThan,
		ast.LessThanOrEquals,
		ast.GreaterThanOrEquals,
		ast.Exists,
		ast.DoesNotExist,
		ast.Positive,
		ast.Negates,
	},
	model.DOUBLE: {
		ast.In,
		ast.NotIn,
		ast.Equals,
		ast.NotEquals,
		ast.LessThan,
		ast.GreaterThan,
		ast.LessThanOrEquals,
		ast.GreaterThanOrEquals,
		ast.Exists,
		ast.DoesNotExist,
		ast.Positive,
		ast.Negates,
	},
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
