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

package cdc

import (
	"errors"
	"fmt"
	"sync"
)

type EventType int

const (
	Create EventType = iota
	Update
	Delete
)

func (e EventType) String() string {
	switch e {
	case Create:
		return "create"
	case Update:
		return "update"
	case Delete:
		return "delete"
	}
	return "unknown"
}

type Event struct {
	Type EventType
	Row  interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("%s : %+v", e.Type, e.Row)
}

type EventHandler interface {
	OnEvents(transaction []Event)
}

type Watcher interface {
	RemoveEventHandler(handler EventHandler)
	AddEventHandler(handler EventHandler)
	Start() error
	Close() error
}

var drivers = map[string]func(source string) (Watcher, error){}
var lock = sync.Mutex{}

var ErrUnknownDriver = errors.New("unknown driver")

func Build(name DriverName, source Source) (Watcher, error) {
	lock.Lock()
	var builder = drivers[string(name)]
	lock.Unlock()
	if builder == nil {
		return nil, ErrUnknownDriver
	}
	return builder(string(source))
}

type DriverName string
type Source string

func Register(name string, builder func(source string) (Watcher, error)) {
	lock.Lock()
	drivers[name] = builder
	lock.Unlock()
}
