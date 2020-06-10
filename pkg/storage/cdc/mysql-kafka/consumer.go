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

package mysql_kafka

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cocotyty/forceset"
	"github.com/juju/loggo"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/storage/cdc"
	"github.com/zhihu/cmdb/pkg/tools/kafka"
)

var log = loggo.GetLogger("mysql-kafka-cdc")

var nodename, _ = os.Hostname()

type Consumer struct {
	closer       io.Closer
	version      string
	addrs        []string
	topic        string
	mutex        sync.Mutex
	transactions map[int][]*Event
	handlers     map[cdc.EventHandler]struct{}
}

func (c *Consumer) AddEventHandler(handler cdc.EventHandler) {
	c.mutex.Lock()
	if c.handlers == nil {
		c.handlers = map[cdc.EventHandler]struct{}{}
	}
	c.handlers[handler] = struct{}{}
	c.mutex.Unlock()
}

func (c *Consumer) RemoveEventHandler(handler cdc.EventHandler) {
	c.mutex.Lock()
	if c.handlers == nil {
		c.handlers = map[cdc.EventHandler]struct{}{}
	}
	delete(c.handlers, handler)
	c.mutex.Unlock()
}

func (c *Consumer) Start() error {
	closer, err := kafka.StartGroupConsume(kafka.Config{
		Addr:     c.addrs,
		Group:    "cmdb_consume_" + nodename,
		Topic:    []string{c.topic},
		Version:  c.version,
		ClientID: "cmdb_consume_" + nodename,
	}, func(msg *sarama.ConsumerMessage) {
		var evt = Event{}
		log.Infof("%s", string(msg.Value))
		err := json.Unmarshal(msg.Value, &evt)
		if err != nil {
			log.Errorf("Unmarshal value %s failed: %s", string(msg.Value), err)
			return
		}
		c.handle(&evt)
	})
	if err != nil {
		return err
	}
	c.closer = closer
	return nil
}

func (c *Consumer) handle(evt *Event) {
	if c.transactions == nil {
		c.transactions = map[int][]*Event{}
	}
	events, ok := c.transactions[evt.XID]
	events = append(events, evt)
	if !evt.Commit {
		c.transactions[evt.XID] = events
		return
	}
	if ok {
		// clear buffer
		delete(c.transactions, evt.XID)
	}
	c.handleTransaction(events)
	return
}

var convertOpt = func(opt *forceset.SetOption) {
	opt.Tag = "db"
	opt.Mappers[forceset.MapperType{
		Destination: reflect.TypeOf(time.Time{}),
		Source:      reflect.TypeOf(""),
	}] = func(dst reflect.Value, src reflect.Value, _ string) error {
		fmt := "2006-01-02 15:04:05"
		parse, err := time.Parse(fmt, src.Interface().(string))
		if err != nil {
			return err
		}
		dst.Set(reflect.ValueOf(parse))
		return nil
	}
}

func (c *Consumer) handleTransaction(evts []*Event) {
	var events []cdc.Event
	var dst interface{}
	for _, evt := range evts {
		var e = cdc.Event{}
		switch evt.Table {
		case "deleted_object":
			dst = &model.DeletedObject{}
		case "deleted_object_log":
			dst = &model.DeletedObjectLog{}
		case "deleted_object_meta_value":
			dst = &model.DeletedObjectMetaValue{}
		case "deleted_object_relation":
			dst = &model.DeletedObjectRelation{}
		case "deleted_object_relation_meta_value":
			dst = &model.DeletedObjectMetaValue{}
		case "deleted_object_relation_type":
			dst = &model.DeletedObjectRelationType{}
		case "object":
			dst = &model.Object{}
		case "object_log":
			dst = &model.ObjectLog{}
		case "object_meta":
			dst = &model.ObjectMeta{}
		case "object_meta_value":
			dst = &model.ObjectMetaValue{}
		case "object_relation":
			dst = &model.ObjectRelation{}
		case "object_relation_meta":
			dst = &model.ObjectRelationMeta{}
		case "object_relation_meta_value":
			dst = &model.ObjectRelationMetaValue{}
		case "object_relation_type":
			dst = &model.ObjectRelationType{}
		case "object_state":
			dst = &model.ObjectState{}
		case "object_status":
			dst = &model.ObjectStatus{}
		case "object_type":
			dst = &model.ObjectType{}
		}
		err := forceset.Set(dst, evt.Data, convertOpt)
		if err != nil {
			panic(err)
		}
		e.Row = dst
		switch evt.Type {
		case Insert:
			e.Type = cdc.Create
		case Update:
			e.Type = cdc.Update
		case Delete:
			e.Type = cdc.Delete
		}
		events = append(events, e)
	}
	log.Infof("%s", events)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for handler := range c.handlers {
		handler.OnEvents(events)
	}
}

func (c *Consumer) Close() error {
	return c.closer.Close()
}

type EventType string

const (
	Insert EventType = "insert"
	Update EventType = "update"
	Delete EventType = "delete"
)

type Event struct {
	Database string                 `json:"database"`
	Table    string                 `json:"table"`
	Type     EventType              `json:"type"`
	TS       int                    `json:"ts"`
	XID      int                    `json:"xid"`
	XOffset  int                    `json:"xoffset"`
	Commit   bool                   `json:"commit"`
	Data     map[string]interface{} `json:"data"`
	Old      map[string]interface{} `json:"old"`
}

const DriverName = "mysql-kafka"

func init() {
	cdc.Register(DriverName, func(source string) (w cdc.Watcher, err error) {
		c := &Consumer{}
		u, err := url.Parse(source)
		if err != nil {
			return nil, err
		}
		addrs := strings.Split(u.Host, ",")
		var values = u.Query()
		var topic = values.Get("topic")
		var version = values.Get("version")
		c.version = version
		c.topic = topic
		c.addrs = addrs
		return c, nil
	})
}
