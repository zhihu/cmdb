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

package tidb_kafka

import (
	"io"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/juju/loggo"
	cmodel "github.com/pingcap/ticdc/cdc/model"
	"github.com/pingcap/ticdc/cdc/sink/codec"
	"github.com/zhihu/cmdb/pkg/storage/cdc"
	"github.com/zhihu/cmdb/pkg/tools/kafka"
)

var log = loggo.GetLogger("ticdc-kafka-cdc")

var nodename, _ = os.Hostname()

type Consumer struct {
	closer       io.Closer
	version      string
	addrs        []string
	topic        string
	mutex        sync.Mutex
	transactions map[uint64][]*Event
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
		log.Infof("msg :%s", msg.Key)
		batchDecoder, err := codec.NewJSONEventBatchDecoder(msg.Key, msg.Value)
		if err != nil {
			log.Errorf("create decoder: %s", err)
			return
		}
		for {
			tp, hasNext, err := batchDecoder.HasNext()
			if err != nil {
				log.Errorf("Decoder: %s", err)
				return
			}
			if !hasNext {
				break
			}
			switch tp {
			case cmodel.MqMessageTypeRow:
				row, err := batchDecoder.NextRowChangedEvent()
				if err != nil {
					log.Errorf("MqMessageTypeRow: %s", err)
					return
				}
				obj, err := Convert(row)
				events := c.transactions[row.CommitTs]
				var typ cdc.EventType = cdc.Update
				if row.Delete {
					typ = cdc.Delete
				}
				events = append(events, &Event{
					TS:     row.CommitTs,
					Object: obj,
					Type:   typ,
				})
				log.Debugf("%s: %s", typ, obj)
			case cmodel.MqMessageTypeResolved:
				ts, err := batchDecoder.NextResolvedEvent()
				if err != nil {
					log.Infof("MqMessageTypeResolved: %s", err)
					return
				}
				log.Infof("resolve: %d", ts)
				for _, events := range c.transactions {
					c.handleTransaction(events)
				}
				c.transactions = nil
			}
		}
	})
	if err != nil {
		return err
	}
	c.closer = closer
	return nil
}

func (c *Consumer) handleTransaction(evts []*Event) {
	var events []cdc.Event
	for _, evt := range evts {
		var e = cdc.Event{}
		e.Type = evt.Type
		e.Row = evt.Object
		events = append(events, e)
	}
	log.Infof("transaction: %s", events)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for handler := range c.handlers {
		handler.OnEvents(events)
	}
}

func (c *Consumer) Close() error {
	return c.closer.Close()
}

const DriverName = "tidb-kafka"

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

type Event struct {
	TS     uint64
	Object interface{}
	Type   cdc.EventType
}
