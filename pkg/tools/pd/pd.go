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

package pd

import (
	"context"

	pd "github.com/pingcap/pd/v4/client"
	"github.com/pingcap/tidb/store/tikv/oracle"
	"github.com/zhihu/cmdb/pkg/storage"
)

type Config struct {
	Addr []string
	pd.SecurityOption
}

func NewTimestampGetter(conf *Config) (storage.TimestampGetter, error) {
	cli, err := pd.NewClient(conf.Addr, pd.SecurityOption{
		CAPath:   conf.CAPath,
		CertPath: conf.CertPath,
		KeyPath:  conf.KeyPath,
	})
	if err != nil {
		return nil, err
	}
	return &TimestampGetter{cli: cli}, nil
}

type TimestampGetter struct {
	cli pd.Client
}

func (g *TimestampGetter) Get(ctx context.Context) (ts uint64, err error) {
	ps, ls, err := g.cli.GetTS(ctx)
	if err != nil {
		return
	}
	ts = oracle.ComposeTS(ps, ls)
	return ts, nil
}
