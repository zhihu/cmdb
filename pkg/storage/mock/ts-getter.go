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

package mock

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/pingcap/tidb/store/tikv/oracle"
)

const maxLogical = int64(1 << 18)

func NewTimestampGetter() *TimestampGetter {
	return &TimestampGetter{}
}

type TimestampGetter struct {
	counter int64
}

func (t *TimestampGetter) Get(ctx context.Context) (uint64, error) {
	physical := time.Now().UnixNano() / int64(time.Millisecond)
	logical := atomic.AddInt64(&t.counter, 1)
	if logical > maxLogical {
		return 0, errors.New("reach max logical")
	}
	ts := oracle.ComposeTS(physical, logical)
	return ts, nil
}
