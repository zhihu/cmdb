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
	"strings"

	"github.com/juju/loggo"
	"github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/query"
	"github.com/zhihu/cmdb/pkg/storage"
	"github.com/zhihu/cmdb/pkg/storage/cdc"
)

var log = loggo.GetLogger("server")

type Objects struct {
	Storage storage.Storage
}

func (o *Objects) Delete(ctx context.Context, request *v1.ObjectDeleteRequest) (*v1.Object, error) {
	panic("implement me")
}

func (o *Objects) Update(ctx context.Context, request *v1.ObjectUpdateRequest) (*v1.Object, error) {
	action := storage.ObjectUpdateOption{}
	var updateMetas = map[string]*v1.ObjectMetaValue{}
	for _, path := range request.UpdateMask.Paths {
		var names = strings.Split(path, ".")
		if len(names) == 0 {
			continue
		}
		switch names[0] {
		case "metas":
			if len(names) == 1 {
				action.SetAllMeta = true
			} else {
				updateMetas[names[1]] = request.Object.Metas[names[1]]
			}
			continue
		case "status":
			action.SetStatus = true
		case "state":
			action.SetState = true
		case "description":
			action.SetDescription = true
		}
	}
	if !action.SetAllMeta {
		request.Object.Metas = updateMetas
	}
	return o.Storage.UpdateObject(ctx, action, request.Object)
}

func (o *Objects) Get(ctx context.Context, request *v1.ObjectGetRequest) (*v1.Object, error) {
	return o.Storage.GetObject(ctx, request.Type, request.Name)
}

func (o *Objects) Watch(request *v1.ObjectListRequest, server v1.Objects_WatchServer) error {
	f := &FilterWatcher{
		server: server,
	}
	if request.Query != "" {
		selector, err := query.Parse(request.Query)
		if err != nil {
			return err
		}
		f.selector = selector
	}
	return o.Storage.WatchObjects(server.Context(), request.Type, f)
}

type FilterWatcher struct {
	server   v1.Objects_WatchServer
	selector storage.Selector
}

func (f *FilterWatcher) OnInit(objects []*v1.Object) {
	_ = f.server.Send(&v1.ObjectWatchEvent{
		Objects: objects,
		Type:    v1.WatchEventType_INIT,
	})
}

func (f *FilterWatcher) Filter(object *v1.Object) bool {
	if f.selector == nil {
		return true
	}
	return f.selector.Match(object.Metas)
}

func (f *FilterWatcher) OnEvent(event storage.ObjectEvent) {
	evt := &v1.ObjectWatchEvent{
		Objects: []*v1.Object{event.Object},
	}
	switch event.Event {
	case cdc.Create:
		evt.Type = v1.WatchEventType_CREATE
	case cdc.Update:
		evt.Type = v1.WatchEventType_UPDATE
	case cdc.Delete:
		evt.Type = v1.WatchEventType_DELETE
	}
	_ = f.server.Send(evt)
}

func (o *Objects) Create(ctx context.Context, object *v1.Object) (*v1.Object, error) {
	n, err := o.Storage.CreateObject(ctx, object)
	return n, err
}

func (o *Objects) List(ctx context.Context, request *v1.ObjectListRequest) (*v1.ObjectListResponse, error) {
	return o.Storage.ListObjects(ctx, request)
}
