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

package tidb_kafka

import (
	"errors"
	"reflect"
	"strings"
	"unicode"

	"github.com/cocotyty/forceset"
	cmodel "github.com/pingcap/ticdc/cdc/model"
	"github.com/zhihu/cmdb/pkg/model"
)

var globalConverter = converter{}

func init() {
	globalConverter.Register(
		model.Object{},
		model.ObjectMeta{},
		model.ObjectMetaValue{},
		model.ObjectStatus{},
		model.ObjectState{},
		model.ObjectType{},
		model.ObjectRelation{},
		model.ObjectRelationType{},
		model.ObjectRelationMeta{},
		model.ObjectRelationMetaValue{},
	)
}

type converter struct {
	tables map[string]table
}

type table struct {
	typ  reflect.Type
	cols map[string]int
}

func (c *converter) Register(types ...interface{}) {
	if c.tables == nil {
		c.tables = map[string]table{}
	}
	for _, t := range types {
		c.register(t)
	}
}

func (c *converter) register(o interface{}) {
	typ := reflect.TypeOf(o)
	var name = typ.Name()
	var tableName = toSnake(name)
	num := typ.NumField()
	var cols = map[string]int{}
	for i := 0; i < num; i++ {
		var f = typ.Field(i)
		if f.PkgPath != "" {
			continue
		}
		var dbTag = f.Tag.Get("db")
		var colName = strings.TrimSpace(strings.Split(dbTag, ";")[0])
		if colName == "" {
			colName = toSnake(f.Name)
		}
		cols[colName] = i
	}
	c.tables[tableName] = table{
		typ:  typ,
		cols: cols,
	}
}

var ErrUnknownTable = errors.New("unknown table")

func Convert(event *cmodel.RowChangedEvent) (row interface{}, err error) {
	return globalConverter.Convert(event)
}

func (c *converter) Convert(event *cmodel.RowChangedEvent) (row interface{}, err error) {
	var tableName = event.Table.Table
	var t, ok = c.tables[tableName]
	if !ok {
		return nil, ErrUnknownTable
	}
	var o = reflect.New(t.typ)
	var st = o.Elem()
	for name, column := range event.Columns {
		index, ok := t.cols[name]
		if !ok {
			continue
		}
		f := st.Field(index)
		err := forceset.ForceSet(f, column.Value)
		if err != nil {
			return nil, err
		}
	}
	return o.Interface(), nil
}

func toSnake(name string) string {
	var rs []rune
	for i, r := range name {
		if unicode.IsLower(r) {
			rs = append(rs, r)
		} else {
			if i != 0 {
				rs = append(rs, '_')
			}
			rs = append(rs, unicode.ToLower(r))
		}
	}
	return string(rs)
}
