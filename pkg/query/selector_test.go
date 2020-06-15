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
	"bytes"
	"encoding/json"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/query/ast"
)

func TestSelector_QuerySQL(t *testing.T) {
	var metas = map[string]*model.ObjectMeta{}
	metas["B"] = &model.ObjectMeta{
		ID:          0,
		TypeID:      1,
		Name:        "B",
		ValueType:   model.BOOLEAN,
		Description: "",
		CreateTime:  time.Time{},
		DeleteTime:  nil,
	}
	metas["S"] = &model.ObjectMeta{
		ID:          0,
		TypeID:      1,
		Name:        "S",
		ValueType:   model.STRING,
		Description: "",
		CreateTime:  time.Time{},
		DeleteTime:  nil,
	}
	metas["D"] = &model.ObjectMeta{
		ID:          0,
		TypeID:      1,
		Name:        "D",
		ValueType:   model.DOUBLE,
		Description: "",
		CreateTime:  time.Time{},
		DeleteTime:  nil,
	}
	metas["I"] = &model.ObjectMeta{
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
	if err != nil {
		t.Fatalf("generate SQL failed: %s", err)
	}
	t.Log(sql, args)
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
		v = RandString(3)
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

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func RandStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandString(length int) string {
	return RandStringWithCharset(length, charset)
}

type testObjectMetas struct {
	Age   int     `json:"age"`
	IDC   string  `json:"idc"`
	Valid bool    `json:"valid"`
	Size  float64 `json:"size"`
}

type testMatch struct {
	metas testObjectMetas
	match bool
}

var matchTestTable = []struct {
	query  string
	tables []testMatch
}{
	{
		query: `age = 12 AND idc in ("idc01","idc02")`,
		tables: []testMatch{
			{testObjectMetas{Age: 12, IDC: "idc01"}, true},
			{testObjectMetas{Age: 12, IDC: "idc01"}, true},
			{testObjectMetas{Age: 12, IDC: "idc03"}, false},
			{testObjectMetas{Age: 14, IDC: "idc01"}, false},
		},
	},
}

func TestSelector_Match(t *testing.T) {
	for _, r := range matchTestTable {
		s, err := Parse(r.query)
		if err != nil {
			t.Fatalf("parse %s error: %s", r.query, err)
		}
		for _, meta := range r.tables {
			match := s.Match(metas(meta.metas))
			if match != meta.match {
				t.Fatalf("metas: %v ,except: %v got: %v", meta.metas, match, meta.match)
			}
		}
	}

}

func metas(o interface{}) map[string]*v1.ObjectMetaValue {
	data, _ := json.Marshal(o)
	var target = map[string]json.RawMessage{}
	_ = json.Unmarshal(data, &target)
	var dest = map[string]*v1.ObjectMetaValue{}
	for k, message := range target {
		var v = v1.ObjectMetaValue{}
		_ = (&jsonpb.Unmarshaler{}).Unmarshal(bytes.NewBuffer(message), &v)
		dest[k] = &v
	}
	return dest
}
