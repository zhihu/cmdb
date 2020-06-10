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

package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zhihu/cmdb/pkg/model/typetables"
	"github.com/zhihu/cmdb/pkg/storage/cdc"
)

type Cache struct {
	typeCache *Types
	m         sync.Map
	watcher   cdc.Watcher
	db        *sqlx.DB
}

var defaultInitTimeout = 5 * time.Second

func NewCache(watcher cdc.Watcher, db *sqlx.DB) (*Cache, error) {
	cache := NewTypes()
	watcher.AddEventHandler(cache)
	ctx, cancel := context.WithTimeout(context.Background(), defaultInitTimeout)
	defer cancel()
	err := cache.InitData(ctx, db)
	if err != nil {
		watcher.RemoveEventHandler(cache)
		return nil, err
	}
	err = watcher.Start()
	if err != nil {
		return nil, err
	}
	return &Cache{
		typeCache: cache,
		watcher:   watcher,
		db:        db,
	}, nil
}

func (c *Cache) TypeCache(fn func(d *typetables.Database)) {
	c.typeCache.Read(fn)
}

func (c *Cache) ReloadTypeCache(ctx context.Context, db *sqlx.DB) error {
	return c.typeCache.InitData(ctx, db)
}

func (c *Cache) GetObjectsCache(ctx context.Context, name string) (*Objects, error) {
	if loaded, ok := c.m.Load(name); ok {
		return loaded.(*TypeObjectsIniter).Init(ctx)
	}
	var id int
	c.typeCache.Read(func(d *typetables.Database) {
		var obj, ok = d.ObjectTypeTable.GetByName(name)
		if ok {
			id = obj.ID
		}
	})
	if id == 0 {
		return nil, errors.New("no such type")
	}
	i := &TypeObjectsIniter{
		typeCache: c.typeCache,
		id:        id,
		name:      name,
		watcher:   c.watcher,
		db:        c.db,
	}
	act, _ := c.m.LoadOrStore(name, i)
	return act.(*TypeObjectsIniter).Init(ctx)
}

type TypeObjectsIniter struct {
	typeCache *Types
	id        int
	name      string
	mutex     sync.Mutex
	inited    bool
	watcher   cdc.Watcher
	o         *Objects
	db        *sqlx.DB
}

func (i *TypeObjectsIniter) Init(ctx context.Context) (o *Objects, err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.inited {
		return i.o, nil
	}
	typeObjects := NewObjects(i.typeCache, i.id, i.name)
	i.watcher.AddEventHandler(typeObjects)
	err = typeObjects.LoadData(ctx, i.db)
	if err != nil {
		typeObjects.ResetBuffer()
		i.watcher.RemoveEventHandler(typeObjects)
		return nil, err
	}
	i.inited = true
	i.o = typeObjects
	return i.o, nil
}
