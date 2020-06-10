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
	"context"

	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/storage/cdc"
)

type Storage interface {
	ListObjects(ctx context.Context, request *v1.ObjectListRequest) (*v1.ObjectListResponse, error)
	WatchObjects(ctx context.Context, typ string, f FilterWatcher) error
	CreateObject(ctx context.Context, obj *v1.Object) (n *v1.Object, err error)
	DeleteObject(ctx context.Context, typ, name string) (n *v1.Object, err error)
	GetObject(ctx context.Context, typ, name string) (n *v1.Object, err error)
	UpdateObject(ctx context.Context, option ObjectUpdateOption, obj *v1.Object) (n *v1.Object, err error)
	StopWatchObjects(ctx context.Context, typ string, f FilterWatcher) error

	GetObjectType(ctx context.Context, name string, consistent bool) (n *v1.ObjectType, err error)
	CreateObjectType(ctx context.Context, typ *v1.ObjectType) (n *v1.ObjectType, err error)
	ListObjectTypes(ctx context.Context, request *v1.ListObjectTypesRequest) (list []*v1.ObjectType, err error)
	UpdateObjectType(ctx context.Context, paths []string, typ *v1.ObjectType) (n *v1.ObjectType, err error)
	DeleteObjectType(ctx context.Context, name string) (n *v1.ObjectType, err error)
}

type ObjectUpdateOption struct {
	SetStatus      bool
	SetState       bool
	SetDescription bool
	SetAllMeta     bool
	MatchVersion   bool
}

type Selector interface {
	Match(metas map[string]*v1.ObjectMetaValue) bool
	QuerySQL(metas map[string]*model.ObjectMeta) (sql string, args []interface{}, err error)
}

type ObjectEvent struct {
	Object *v1.Object
	Event  cdc.EventType
}

type FilterWatcher interface {
	OnInit(object []*v1.Object)
	Filter(object *v1.Object) bool
	OnEvent(ObjectEvent)
}

type TimestampGetter interface {
	Get(ctx context.Context) (ts uint64, err error)
}
