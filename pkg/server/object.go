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

package server

import (
	"context"

	"github.com/juju/loggo"
	"github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/storage"
)

var log = loggo.GetLogger("server")

type Objects struct {
	Storage storage.Storage
}

func (o *Objects) List(ctx context.Context, request *v1.ObjectListRequest) (*v1.ObjectListResponse, error) {
	return o.Storage.ListObjects(request)
}
