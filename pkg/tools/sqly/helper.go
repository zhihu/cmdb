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

package sqly

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Execer struct {
	Ctx context.Context
	Tx  *sqlx.Tx
	Err error
}

func (e *Execer) Exec(query string, args ...interface{}) {
	if e.Err != nil {
		return
	}
	_, err := e.Tx.ExecContext(e.Ctx, query, args...)
	if err != nil {
		e.Err = err
	}
}
