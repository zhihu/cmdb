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
	"testing"

	"github.com/zhihu/cmdb/pkg/query/ast"
	"github.com/zhihu/cmdb/pkg/query/lexer"
	"github.com/zhihu/cmdb/pkg/query/parser"
)

func TestAST(t *testing.T) {
	l := lexer.NewLexer([]byte(`x == "12askjdhkjajskd" && !kk &&  y in (1,23,4,5)`))
	res, err := parser.NewParser().Parse(l)
	if err != nil {
		t.Fatal(err)
	}
	requirements := res.([]*ast.Requirement)
	for _, requirement := range requirements {
		t.Logf("%#v", requirement)
	}
}
