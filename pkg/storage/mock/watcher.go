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

package mock

import (
	"sync"

	"github.com/zhihu/cmdb/pkg/storage/cdc"
)

func NewWatcher() *Watcher {
	return &Watcher{}
}

type Watcher struct {
	handlers []cdc.EventHandler
	locker   sync.RWMutex
}

func (w *Watcher) RemoveEventHandler(handler cdc.EventHandler) {
	w.locker.Lock()
	var n []cdc.EventHandler
	for _, eventHandler := range w.handlers {
		if eventHandler != handler {
			n = append(n, eventHandler)
		}
	}
	w.locker.Unlock()
}

func (w *Watcher) AddEventHandler(handler cdc.EventHandler) {
	w.locker.Lock()
	var n []cdc.EventHandler
	for _, eventHandler := range w.handlers {
		if eventHandler != handler {
			n = append(n, eventHandler)
		}
	}
	n = append(n, handler)
	w.locker.Unlock()
}

func (w *Watcher) TellAllHandler(events []cdc.Event) {
	w.locker.RLock()
	handlers := make([]cdc.EventHandler, len(w.handlers))
	copy(handlers, w.handlers)
	w.locker.RUnlock()
	for _, handler := range handlers {
		handler.OnEvents(events)
	}
}

func (w *Watcher) Start() error {
	return nil
}

func (w *Watcher) Close() error {
	return nil
}
