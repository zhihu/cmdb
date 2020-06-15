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

package mtables

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var versionString = regexp.MustCompile(`^v\d+$`)

var nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9]`)

func PkgShortName(path string) string {
	p := strings.Split(path, "/")
	var last = p[len(p)-1]
	if len(p) > 1 && versionString.MatchString(last) {
		last = p[len(p)-2]
	}
	last = strings.ToLower(last)
	return nonAlphanumeric.ReplaceAllString(last, "_")
}

func GenerateTablesInfo(tagName string, softDeleteFieldName string, tables ...interface{}) (map[string]*Table, error) {
	var m = map[string]*Table{}
	var relations = map[string][]relation{}
	for _, table := range tables {
		typ := reflect.TypeOf(table)
		table := &Table{
			Name:      typ.Name(),
			ID:        Index{},
			Type:      typ,
			Relations: nil,
			Indexes:   map[string]Index{},
		}
		var id *reflect.StructField

		typ.FieldByNameFunc(func(s string) bool {
			f, _ := typ.FieldByName(s)
			if f.PkgPath != "" {
				return false
			}

			if strings.ToUpper(s) == "ID" {
				id = &f
			}
			tag := splitTag(f.Tag.Get(tagName))
			var idx, ok = tag["index"]
			if ok && idx == "" {
				idx = s
			}
			if idx != "" {
				parts := strings.Split(idx, ",")
				var name = strings.TrimSpace(parts[0])
				var ii = table.Indexes[name]
				if len(parts) > 1 && strings.TrimSpace(parts[1]) == "unique" {
					ii.Unique = true
				}
				ii.Fields = append(ii.Fields, Field{
					Name: f.Name,
					Type: f.Type,
				})
				ii.Name = name
				table.Indexes[name] = ii
			}
			r := relationTag(tag["belongsTo"], f)
			if r.tableName != "" {
				relations[table.Name] = append(relations[table.Name], r)
			}
			if softDeleteFieldName != "" && softDeleteFieldName == s {
				table.SoftDelete = &Field{
					Name: f.Name,
					Type: f.Type,
				}
			}
			return false
		})
		if _, ok := table.Indexes["ID"]; !ok {
			if id == nil {
				return nil, errors.New("not have id index: " + typ.Name())
			}
			table.Indexes["ID"] = Index{Name: "ID", Unique: true, Fields: []Field{
				{Name: id.Name, Type: id.Type},
			}}
		}
		table.ID = table.Indexes["ID"]
		m[table.Name] = table
	}
	for name, rs := range relations {
		table := m[name]
		for _, r := range rs {
			target, ok := m[r.tableName]
			if !ok {
				return nil, errors.New("no such table: " + r.tableName)
			}

			field, ok := target.Type.FieldByName(r.fieldsName)
			if !ok {
				return nil, fmt.Errorf("no such field %s in table: %s ", r.tableName, r.fieldsName)
			}
			var current = Relation{
				From: Field{
					Name: r.f.Name,
					Type: r.f.Type,
				},
				To: Field{
					Name: field.Name,
					Type: field.Type,
				},
				ToTable: target,
				Type:    BelongsTo,
			}
			var reverseRelation = Relation{
				To: Field{
					Name: r.f.Name,
					Type: r.f.Type,
				},
				From: Field{
					Name: field.Name,
					Type: field.Type,
				},
				ToTable: table,
				Type:    r.reverseType,
			}
			table.Relations = append(table.Relations, current)
			target.Relations = append(target.Relations, reverseRelation)
			table.Indexes[current.From.Name] = Index{
				Name:   current.From.Name,
				Unique: false,
				Fields: []Field{current.From},
			}
		}
	}
	return m, nil
}

type relation struct {
	f           reflect.StructField
	reverseType RelationType
	tableName   string
	fieldsName  string
}

func relationTag(value string, f reflect.StructField) (r relation) {
	values := strings.Split(value, ",")
	r.fieldsName = "ID"
	r.f = f
	for i, v := range values {
		v = strings.TrimSpace(v)
		switch i {
		case 0:
			switch v {
			case "one":
				r.reverseType = HasOne
			case "many":
				r.reverseType = HasMany
			}
		case 1:
			r.tableName = v
		case 2:
			r.fieldsName = v
		}
	}
	return r
}

// index: indexName;belongsTo: one,tableName,fieldsName;
func splitTag(tag string) map[string]string {
	m := map[string]string{}
	list := strings.Split(tag, ";")
	for _, s := range list {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		kv := strings.Split(s, ":")
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		} else {
			m[kv[0]] = ""
		}
	}
	return m
}
