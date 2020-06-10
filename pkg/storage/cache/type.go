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
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/model/typetables"
	. "github.com/zhihu/cmdb/pkg/storage/cdc"
)

type Types struct {
	database       typetables.Database
	locker         sync.RWMutex
	loaded         bool
	bufferedEvents [][]Event
}

func NewTypes() *Types {
	return &Types{}
}

func (c *Types) OnEvents(transaction []Event) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if !c.loaded {
		c.bufferedEvents = append(c.bufferedEvents, transaction)
		return
	}
	c.consumeEvents(transaction)
}

func (c *Types) consumeEvents(transaction []Event) {
	c.database.OnEvents(transaction)
}

func (c *Types) Read(fn func(d *typetables.Database)) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	fn(&c.database)
}

func (c *Types) InitData(ctx context.Context, db *sqlx.DB) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()
	var types []*model.ObjectType
	var status []*model.ObjectStatus
	var state []*model.ObjectState
	var meta []*model.ObjectMeta
	err = tx.SelectContext(ctx, &types, `select * from object_type where delete_time is null`)
	if err != nil {
		return err
	}
	err = tx.SelectContext(ctx, &status, `select * from object_status where delete_time is null`)
	if err != nil {
		return err
	}
	err = tx.SelectContext(ctx, &state, `select * from object_state where delete_time is null`)
	if err != nil {
		return err
	}
	err = tx.SelectContext(ctx, &meta, `select * from object_meta where delete_time is null`)
	if err != nil {
		return err
	}
	_ = tx.Commit()

	c.locker.Lock()
	defer c.locker.Unlock()
	c.database.Init()
	c.loaded = true
	for _, row := range types {
		c.database.InsertObjectType(row)
	}
	for _, row := range status {
		c.database.InsertObjectStatus(row)
	}
	for _, row := range state {
		c.database.InsertObjectState(row)
	}
	for _, row := range meta {
		c.database.InsertObjectMeta(row)
	}
	for _, events := range c.bufferedEvents {
		c.database.OnEvents(events)
	}
	return nil
}
