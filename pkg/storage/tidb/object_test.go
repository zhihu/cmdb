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

package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/storage/cdc"
)

var testObjectType = v1.CheckObjectType(&v1.ObjectType{
	Name: "object_test_resource",
	Statuses: map[string]*v1.ObjectStatus{
		"CREATED": {
			States: map[string]*v1.ObjectState{
				"READY": {},
			},
		},
	},
	Metas: map[string]*v1.ObjectMeta{
		"size": {
			ValueType: v1.ValueType_INTEGER,
		},
	},
})

func prepareObjectTests(t *testing.T) {
	storage := getTestStorage(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, _ = storage.ForceDeleteObjectType(ctx, testObjectType.Name)
	testWatcher.TellAllHandler([]cdc.Event{
		{Row: &model.ObjectType{}},
	})
	_, err := storage.CreateObjectType(ctx, testObjectType)
	if err != nil {
		t.Fatalf("unexcepted error happened when execute CreateObjectType: %s", err)
	}
	err = testCache.ReloadTypeCache(ctx, testDB)
	if err != nil {
		t.Fatalf("unexcepted error happened when execute ReloadTypeCache: %s", err)
	}
	return
}

func TestStorage_CreateGetUpdateDeleteObject(t *testing.T) {
	prepareObjectTests(t)
	storage := getTestStorage(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	o := &v1.Object{
		Type:   testObjectType.Name,
		Name:   "test",
		Status: "CREATED",
		State:  "READY",
		Metas: map[string]*v1.ObjectMetaValue{
			"size": {
				Value:     "1200",
				ValueType: v1.ValueType_INTEGER,
			},
		},
	}
	_, err := storage.CreateObject(ctx, o)
	if err != nil {
		t.Fatalf("unexcepted error happened when execute CreateObject: %s", err)
	}
	obj, err := storage.GetObject(ctx, testObjectType.Name, o.Name)
	if err != nil {
		t.Fatalf("unexcepted error happened when execute GetObject: %s", err)
	}
	if !objectEquals(o, obj) {
		t.Fatalf("except created object equals to got object: %s %s", objectString(o), objectString(obj))
	}
	objects, err := storage.ListObjects(ctx, &v1.ObjectListRequest{
		Type:        testObjectType.Name,
		View:        v1.ObjectView_NORMAL,
		ShowDeleted: true,
	})
	if err != nil {
		t.Fatalf("unexcepted error happened when execute ListObjects: %s", err)
	}
	if len(objects.Objects) != 1 {
		t.Fatalf("")
	}
	if !objectEquals(o, objects.Objects[0]) {
		t.Fatalf("except created object equals to got object: %s %s", objectString(o), objectString(objects.Objects[0]))
	}

	_, err = storage.DeleteObject(ctx, testObjectType.Name, o.Name)
	if err != nil {
		t.Fatalf("unexcepted error happened when execute DeleteObject: %s", err)
	}

}

func objectEquals(o *v1.Object, o2 *v1.Object) bool {
	if o.Name != o2.Name ||
		o.Type != o2.Type ||
		o.Status != o2.Status ||
		o.State != o2.State ||
		len(o.Metas) != len(o2.Metas) {
		return false
	}
	for name, value := range o.Metas {
		v2, ok := o2.Metas[name]
		if !ok {
			return false
		}
		if v2.Value != value.Value {
			return false
		}
	}
	return true
}

func objectString(o *v1.Object) string {
	m := jsonpb.Marshaler{Indent: "  "}
	str, _ := m.MarshalToString(o)
	return str
}
