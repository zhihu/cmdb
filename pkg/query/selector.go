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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zhihu/cmdb/pkg/domain"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/query/ast"
	"github.com/zhihu/cmdb/pkg/query/lexer"
	"github.com/zhihu/cmdb/pkg/query/parser"
	"github.com/zhihu/cmdb/pkg/storage"
)

var ErrNoCondition = errors.New("query has no condition")

func Parse(query string) (storage.Selector, error) {
	l := lexer.NewLexer([]byte(query))
	res, err := parser.NewParser().Parse(l)
	if err != nil {
		return nil, err
	}
	requires, ok := res.([]*ast.Requirement)
	if !ok {
		return nil, errors.New("invalid query statements")
	}
	if len(requires) == 0 {
		return nil, ErrNoCondition
	}
	return &Selector{
		requires: requires,
	}, nil
}

type Selector struct {
	requires []*ast.Requirement
}

func (s *Selector) Match(metas map[string]*domain.Meta) bool {
	for _, require := range s.requires {
		meta, ok := metas[require.Key]
		if !ok {
			if require.Operator == ast.DoesNotExist {
				continue
			}
			if require.Operator == ast.Negates {
				continue
			}
			return false
		}
		switch require.Operator {
		case ast.Exists:
			continue
		case ast.Positive:
			if b, ok := meta.Value.(bool); ok && !b {
				return false
			}
			continue
		case ast.Negates:
			if b, ok := meta.Value.(bool); ok && b {
				return false
			}
			continue
		case ast.Equals:
			if meta.RawValue != require.Value[0] {
				return false
			}
			continue
		case ast.NotEquals:
			if meta.RawValue == require.Value[0] {
				return false
			}
			continue
		case ast.LessThan:
			switch c := meta.Value.(type) {
			case float64:
				f, err := strconv.ParseFloat(require.Value[0], 64)
				if err != nil {
					return false
				}
				if c < f {
					continue
				}
				return false
			case string:
				if c < require.Value[0] {
					continue
				}
				return false
			case int:
				f, err := strconv.ParseInt(require.Value[0], 10, 64)
				if err != nil {
					return false
				}
				if c < int(f) {
					continue
				}
				return false
			}
		case ast.GreaterThan:
			switch c := meta.Value.(type) {
			case float64:
				f, err := strconv.ParseFloat(require.Value[0], 64)
				if err != nil {
					return false
				}
				if c > f {
					continue
				}
				return false
			case string:
				if c > require.Value[0] {
					continue
				}
				return false
			case int:
				f, err := strconv.ParseInt(require.Value[0], 10, 64)
				if err != nil {
					return false
				}
				if c > int(f) {
					continue
				}
				return false
			}
		case ast.GreaterThanOrEquals:
			switch c := meta.Value.(type) {
			case float64:
				f, err := strconv.ParseFloat(require.Value[0], 64)
				if err != nil {
					return false
				}
				if c >= f {
					continue
				}
				return false
			case string:
				if c >= require.Value[0] {
					continue
				}
				return false
			case int:
				f, err := strconv.ParseInt(require.Value[0], 10, 64)
				if err != nil {
					return false
				}
				if c >= int(f) {
					continue
				}
				return false
			}
		case ast.LessThanOrEquals:
			switch c := meta.Value.(type) {
			case float64:
				f, err := strconv.ParseFloat(require.Value[0], 64)
				if err != nil {
					return false
				}
				if c <= f {
					continue
				}
				return false
			case string:
				if c <= require.Value[0] {
					continue
				}
				return false
			case int:
				f, err := strconv.ParseInt(require.Value[0], 10, 64)
				if err != nil {
					return false
				}
				if c <= int(f) {
					continue
				}
				return false
			}
		case ast.In:
			for _, v := range require.Value {
				if v != meta.RawValue {
					return false
				}
			}
			continue
		case ast.NotIn:
			for _, v := range require.Value {
				if v == meta.RawValue {
					return false
				}
			}
			continue
		}
	}
	return true
}

type sqlAndArgs struct {
	sql  string
	args []interface{}
}

func (s *Selector) QuerySQL(metas map[string]*model.ObjectMeta) (sql string, args []interface{}, err error) {
	var conditions []sqlAndArgs
	for _, require := range s.requires {
		m, ok := metas[require.Key]
		if !ok {
			continue
		}
		if m.ValueType != model.BOOLEAN {
			switch require.Operator {
			case ast.Exists, ast.Positive:
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ?", []interface{}{
					m.ID,
				}})
				continue
			case ast.DoesNotExist, ast.Negates:
				conditions = append(conditions, sqlAndArgs{"`meta_id` <> ?", []interface{}{
					m.ID,
				}})
				continue
			}
		}
		switch m.ValueType {
		case model.STRING:
			switch require.Operator {
			case ast.In:
				p := strings.Repeat(",?", len(require.Value))[1:]
				var args = make([]interface{}, 0, len(require.Value)+1)
				args = append(args, m.ID)
				for _, v := range require.Value {
					args = append(args, v)
				}
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` in (" + p + ")", args})
			case ast.NotIn:
				p := strings.Repeat(",?", len(require.Value))[1:]
				var args = make([]interface{}, 0, len(require.Value)+1)
				args = append(args, m.ID)
				for _, v := range require.Value {
					args = append(args, v)
				}
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` not in (" + p + ")", args})
			case ast.Equals,
				ast.NotEquals,
				ast.GreaterThan, ast.GreaterThanOrEquals, ast.LessThan, ast.LessThanOrEquals:
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` " + string(require.Operator) + " ?", []interface{}{
					m.ID, require.Value[0],
				}})
			default:
				return "", nil, fmt.Errorf("STRING typed not support such operator %s: %w", require.Operator, ErrInvalidOperator)
			}
			continue
		case model.BOOLEAN:
			switch require.Operator {
			case ast.Exists:
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ?", []interface{}{
					m.ID,
				}})
			case ast.DoesNotExist:
				conditions = append(conditions, sqlAndArgs{"`meta_id` <> ?", []interface{}{
					m.ID,
				}})
			case ast.Equals:
				b, err := strconv.ParseBool(require.Value[0])
				if err != nil {
					return "", nil, err
				}
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` " + string(require.Operator) + " ?", []interface{}{
					m.ID, model.BooleanValue(b),
				}})
			case ast.Positive:
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` = `true`", nil})
			case ast.Negates:
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` = `false`", nil})
			default:
				return "", nil, fmt.Errorf("BOOLEAN typed not support such operator %s: %w", require.Operator, ErrInvalidOperator)
			}
			continue
		case model.INTEGER:
			switch require.Operator {
			case ast.In:
				p := strings.Repeat(",?", len(require.Value))[1:]
				var args = make([]interface{}, 0, len(require.Value)+1)
				args = append(args, m.ID)
				for _, v := range require.Value {
					_, err := strconv.Atoi(v)
					if err != nil {
						return "", nil, err
					}
					args = append(args, v)
				}
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` in (" + p + ")", args})
			case ast.NotIn:
				p := strings.Repeat(",?", len(require.Value))[1:]
				var args = make([]interface{}, 0, len(require.Value)+1)
				args = append(args, m.ID)
				for _, v := range require.Value {
					_, err := strconv.Atoi(v)
					if err != nil {
						return "", nil, err
					}
					args = append(args, v)
				}
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` not in (" + p + ")", args})
			case ast.Equals,
				ast.NotEquals,
				ast.GreaterThan, ast.GreaterThanOrEquals, ast.LessThan, ast.LessThanOrEquals:
				v := require.Value[0]
				i, err := strconv.Atoi(v)
				if err != nil {
					return "", nil, err
				}
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and CAST(`value` as SIGNED) " + string(require.Operator) + " ?", []interface{}{
					m.ID, i,
				}})
			default:
				return "", nil, fmt.Errorf("INT typed not support such operator %s: %w", require.Operator, ErrInvalidOperator)
			}
		case model.DOUBLE:
			switch require.Operator {
			case ast.In:
				p := strings.Repeat(",?", len(require.Value))[1:]
				var args = make([]interface{}, 0, len(require.Value)+1)
				args = append(args, m.ID)
				for _, v := range require.Value {
					_, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return "", nil, err
					}
					args = append(args, v)
				}
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` in (" + p + ")", args})
			case ast.NotIn:
				p := strings.Repeat(",?", len(require.Value))[1:]
				var args = make([]interface{}, 0, len(require.Value)+1)
				args = append(args, m.ID)
				for _, v := range require.Value {
					for _, v := range require.Value {
						_, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return "", nil, err
						}
						args = append(args, v)
					}
					args = append(args, v)
				}
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and `value` not in (" + p + ")", args})
			case ast.Equals,
				ast.NotEquals,
				ast.GreaterThan, ast.GreaterThanOrEquals, ast.LessThan, ast.LessThanOrEquals:
				d, err := strconv.ParseFloat(require.Value[0], 64)
				if err != nil {
					return "", nil, err
				}
				conditions = append(conditions, sqlAndArgs{"`meta_id` = ? and CAST(`value` as DOUBLE) " + string(require.Operator) + " ?", []interface{}{
					m.ID, d,
				}})
			default:
				return "", nil, fmt.Errorf("DOUBLE typed not support such operator %s: %w", require.Operator, ErrInvalidOperator)
			}
		}
	}
	buf := bytes.NewBuffer(nil)
	for i, c := range conditions {
		args = append(args, c.args...)
		if i != len(conditions)-1 {
			buf.WriteString("select `object_id` from `object_meta_value` where delete_time is null and " + c.sql + " and object_id in (\n")
		} else {
			buf.WriteString("select `object_id` from `object_meta_value` where delete_time is null and " + c.sql)
		}
	}
	buf.Write(bytes.Repeat([]byte("\n)\n"), len(conditions)-1))
	return buf.String(), args, nil
}

var ErrInvalidOperator = errors.New("invalid operator")
