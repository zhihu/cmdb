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
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jmoiron/sqlx"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/storage/cache"
	"github.com/zhihu/cmdb/pkg/storage/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var testDB *sqlx.DB

func init() {
	dsn := os.Getenv("TIDB_DSN")
	if dsn != "" {
		testDB = sqlx.MustOpen("mysql", dsn)
	}
}

var testWatcher = mock.NewWatcher()
var testCache *cache.Cache
var testStorage *Storage
var initStorageOnce = sync.Once{}

func getTestStorage(t *testing.T) *Storage {
	if testDB == nil {
		t.Skip("no database,please set TIDB_DSN env to start this test")
	}
	initStorageOnce.Do(func() {
		var err error
		tsGetter := mock.NewTimestampGetter()
		testCache, err = cache.NewCache(testWatcher, testDB)
		if err != nil {
			t.Fatal(err)
		}
		testStorage = NewStorage(testDB, tsGetter, testCache)
	})
	return testStorage
}

func TestStorage_CreateGetUpdateObjectType(t *testing.T) {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	storage := getTestStorage(t)
	_, err := storage.ForceDeleteObjectType(ctx, ExampleObjectType.Name)
	if status.Code(err) == codes.NotFound {
		err = nil
	}
	if err != nil {
		t.Fatal(err)
	}
	o, err := storage.CreateObjectType(ctx, ExampleObjectType)
	if err != nil {
		t.Fatal(err)
	}
	if !ObjectTypeEqual(ExampleObjectType, o) {
		t.Fatal("expected:", JSONStr(ExampleObjectType), " got:", JSONStr(o))
	}
	o, err = storage.GetObjectType(ctx, ExampleObjectType.Name, true)
	if !ObjectTypeEqual(ExampleObjectType, o) {
		t.Fatal("expected:", JSONStr(ExampleObjectType), " got:", JSONStr(o))
	}
	n, err := storage.UpdateObjectType(ctx, []string{
		"metas.size.description",
	}, &v1.ObjectType{Name: "resource", Metas: map[string]*v1.ObjectMeta{
		"size": {
			Name:        "size",
			Description: "halo",
		},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if n.Metas["size"].Description != "halo" {
		t.Fatal("expected: halo got:", n.Metas["size"].Description)
	}
}

var json = &runtime.JSONPb{OrigName: true, EmitDefaults: true}

func JSONStr(i interface{}) string {
	data, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(data)
}

func ObjectTypeEqual(a *v1.ObjectType, b *v1.ObjectType) bool {
	if a.Name != b.Name || a.Description != b.Description {
		return false
	}
	if len(a.Statuses) != len(b.Statuses) {
		return false
	}
	for name, status := range a.Statuses {
		bs, ok := b.Statuses[name]
		if !ok {
			return false
		}
		if status.Name != bs.Name || status.Description != bs.Description {
			return false
		}
		if len(status.States) != len(bs.States) {
			return false
		}
		for name, state := range status.States {
			bState, ok := bs.States[name]
			if !ok {
				return false
			}
			if state.Name != bState.Name || state.Description != bState.Description {
				return false
			}
		}
	}
	for name, meta := range a.Metas {
		bmeta, ok := b.Metas[name]
		if !ok {
			return false
		}
		if meta.Name != bmeta.Name || meta.Description != bmeta.Description || bmeta.ValueType != meta.ValueType {
			return false
		}
	}

	return true
}

var ExampleObjectType = v1.CheckObjectType(&v1.ObjectType{
	Name:        "resource",
	Description: "resource is just for test",
	Statuses: map[string]*v1.ObjectStatus{
		"Incomplete": {
			Name:        "Incomplete",
			Description: "Host not yet ready for use. It has been powered on and entered in Collins but the automated induction process is still being performed.",
			States: map[string]*v1.ObjectState{
				"FAILED": {
					Name:        "FAILED",
					Description: "A service in this state has encountered a problem and may not be operational. It cannot be started nor stopped.",
				},
				"NEW": {
					Name:        "NEW",
					Description: "A service in this state is inactive. It does minimal work and consumes minimal resources.",
				},
			},
		},
		"New": {
			Name:        "New",
			Description: "Host has completed the automated induction process and is waiting for an onsite tech to complete physical intake",
		},
		"Unallocated": {
			Name:        "Unallocated",
			Description: "Host has completed intake process and is ready for use (eg available resource for provisioning into a role)",
		},
	},
	Metas: map[string]*v1.ObjectMeta{
		"Address": {
			Name:        "address",
			Description: "resource's address",
			ValueType:   v1.ValueType_STRING,
		},
		"Size": {
			Name:        "size",
			Description: "resource's size",
			ValueType:   v1.ValueType_INTEGER,
		},
	},
	CreateTime: nil,
})
