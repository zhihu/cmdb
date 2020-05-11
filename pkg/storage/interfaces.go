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

package storage

import (
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/domain"
	"github.com/zhihu/cmdb/pkg/model"
)

type Storage interface {
	ListObjects(*v1.ObjectListRequest) (*v1.ObjectListResponse, error)
}

type Selector interface {
	Match(metas map[string]*domain.Meta) bool
	QuerySQL(metas map[string]*model.ObjectMeta) (sql string, args []interface{}, err error)
}
