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

package kafka

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/juju/loggo"
)

var log = loggo.GetLogger("kafka")

type GroupConsumer struct {
	Config
}

type Config struct {
	Addr     []string
	Group    string
	Topic    []string
	Version  string
	ClientID string
}

func StartGroupConsume(cfg Config, consume func(message *sarama.ConsumerMessage)) (io.Closer, error) {
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, err
	}
	if !version.IsAtLeast(sarama.V0_10_2_0) {
		c := cluster.NewConfig()
		c.Consumer.Offsets.CommitInterval = time.Second
		c.ClientID = cfg.ClientID
		c.Version = version
		consumer, err := cluster.NewConsumer(cfg.Addr, cfg.Group, cfg.Topic, c)
		if err != nil {
			return nil, err
		}
		go func() {
			for msg := range consumer.Messages() {
				consume(msg)
			}
		}()
		go func() {
			for err := range consumer.Errors() {
				log.Errorf("error: %s", err)
			}
		}()
		return consumer, nil
	}
	config := sarama.NewConfig()
	config.Version = version
	config.ClientID = cfg.ClientID
	group, err := sarama.NewConsumerGroup(cfg.Addr, cfg.Group, config)
	if err != nil {
		return nil, err
	}
	h := &handler{consume: consume}
	err = group.Consume(context.Background(), cfg.Topic, h)
	go func() {
		for err := range group.Errors() {
			log.Errorf("error: %s", err)
		}
	}()
	return group, nil
}

type handler struct {
	mutex   sync.Mutex
	consume func(message *sarama.ConsumerMessage)
}

func (h *handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) recoverPanic() {
	if n := recover(); n != nil {
		log.Errorf("consume panic: %s", n)
	}
}

func (h *handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	go func() {
		for msg := range claim.Messages() {
			var msg = msg
			func() {
				defer h.mutex.Unlock()
				defer h.recoverPanic()
				h.mutex.Lock()
				h.consume(msg)
			}()
		}
	}()
	return nil
}
