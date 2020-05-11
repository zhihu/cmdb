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

package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/juju/loggo"
)

func Setup(conf string) error {
	err := loggo.ConfigureLoggers(conf)
	if err != nil {
		return err
	}
	_, err = loggo.DefaultContext().ReplaceWriter(loggo.DefaultWriterName, loggo.NewSimpleWriter(os.Stderr, func(entry loggo.Entry) string {
		ts := entry.Timestamp.In(time.Local).Format("2006-01-02 15:04:05")
		// Just get the basename from the filename
		filename := filepath.Base(entry.Filename)
		return fmt.Sprintf("%s %s %s %s:%d %s", ts, entry.Level, entry.Module, filename, entry.Line, entry.Message)
	}))
	return err
}
