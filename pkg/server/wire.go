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

// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/zhihu/cmdb/pkg/storage"
	"github.com/zhihu/cmdb/pkg/storage/cache"
	"github.com/zhihu/cmdb/pkg/storage/cdc"
	"github.com/zhihu/cmdb/pkg/storage/tidb"
	"github.com/zhihu/cmdb/pkg/tools/database"
	"github.com/zhihu/cmdb/pkg/tools/pd"
)

// Set provide wire's provider set
var Set = wire.NewSet(
	wire.Struct(new(Objects), "*"),
	wire.Struct(new(Server), "*"),
	tidb.NewStorage,
	pd.NewTimestampGetter,
	wire.Struct(new(ObjectTypes), "*"),
	wire.Bind(new(storage.Storage), new(*tidb.Storage)),
	cdc.Build,
	database.MySQL,
	cache.NewCache,
)
