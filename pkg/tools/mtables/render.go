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

//go:generate sh ./generate.sh
import (
	"bytes"
	"reflect"
	"strconv"
	"text/template"
)

type RelationType int

const (
	HasOne RelationType = iota
	HasMany
	BelongsTo
)

func (r RelationType) String() string {
	switch r {
	case HasOne:
		return "has_one"
	case HasMany:
		return "has_many"
	case BelongsTo:
		return "belongs_to"
	}
	return "unknown"
}

type Field struct {
	Name string
	Type reflect.Type
}

type Relation struct {
	From    Field
	To      Field
	ToTable *Table
	Type    RelationType
}

type Index struct {
	Name   string
	Unique bool
	Fields []Field
}

type Table struct {
	Name       string
	ID         Index
	Type       reflect.Type
	Relations  []Relation
	Indexes    map[string]Index
	SoftDelete *Field
}

func (t *Table) GetReverseRelation(r Relation) Relation {
	for _, rel := range t.Relations {
		if rel.To == r.From && rel.From == r.To {
			return rel
		}
	}
	return Relation{}
}

type ImportsContext map[string]string

func (i ImportsContext) TypeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		return "*" + i.TypeName(t.Elem())
	}
	if t.Kind() == reflect.Map {
		return "[" + i.TypeName(t.Key()) + "]" + i.TypeName(t.Elem())
	}
	if t.Kind() == reflect.Slice {
		return "[]" + i.TypeName(t.Elem())
	}
	path := t.PkgPath()
	if path == "" {
		return t.Name()
	}
	name := PkgShortName(path)
	var finalName = name
	if current, ok := i[finalName]; ok && current != path {
		var num = 1
		for {
			finalName = name + "_" + strconv.Itoa(num)
			current, ok = i[finalName]
			if ok && current != path {
				num++
				continue
			}
			break
		}
	}
	i[finalName] = path
	return name + "." + t.Name()
}

type Plugin struct {
	Template string
	Imports  map[string]string // pkg name -> pkg path
}

func Render(tables map[string]*Table, pkg string, plugins ...Plugin) (sourceCode string, err error) {
	var imports = ImportsContext{}
	for _, plugin := range plugins {
		for k, v := range plugin.Imports {
			imports[k] = v
		}
	}
	var fnMap = map[string]interface{}{
		"TypeName": imports.TypeName,
	}
	body := bytes.NewBuffer(nil)

	var data = map[string]interface{}{
		"Tables": tables,
	}

	tpl := template.Must(template.New("").Funcs(fnMap).Parse(bodyTemplate))
	err = tpl.Execute(body, data)
	if err != nil {
		return "", err
	}

	for _, plugin := range plugins {
		tpl, err := template.New("").Funcs(fnMap).Parse(plugin.Template)
		if err != nil {
			return "", err
		}
		err = tpl.Execute(body, data)
		if err != nil {
			return "", err
		}
	}

	file := bytes.NewBuffer(nil)
	tpl = template.Must(template.New("").Parse(headerTemplate))
	err = tpl.Execute(file, map[string]interface{}{
		"Package": pkg,
		"Imports": imports,
	})
	if err != nil {
		return "", err
	}
	_, _ = body.WriteTo(file)
	return file.String(), nil
}
