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

package main

import (
	"io/ioutil"
	"os/exec"

	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/tools/mtables"
)

var cr = `// Copyright 2020 Zhizhesihai (Beijing) Technology Limited.
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

// DO NOT EDIT IT.
`
var cdcTPL = `

func(d *Database) OnEvents(transaction []cdc.Event) {
	for _, event := range transaction {
		switch row := event.Row.(type) {
{{range $tableName,$table:= $.Tables -}}
			case *{{TypeName $table.Type}}:
				switch event.Type {
					case cdc.Create:
						d.Insert{{$table.Type.Name}}(row)
					case cdc.Update:
						d.Update{{$table.Type.Name}}(row)
					case cdc.Delete:
						d.Delete{{$table.Type.Name}}(row)
				}
{{end -}}
		}
	}
}
`

func main() {
	var generatePath = "../tables.gen.go"
	tablesInfo, err := mtables.GenerateTablesInfo("table", "DeleteTime", model.Object{}, model.ObjectMetaValue{})
	if err != nil {
		panic(err)
	}
	data, err := mtables.Render(tablesInfo, "objects", mtables.Plugin{
		Template: cdcTPL,
		Imports: map[string]string{
			"cdc": "github.com/zhihu/cmdb/pkg/storage/cdc",
		},
	})
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(generatePath, []byte(cr+data), 0666)
	exec.Command("go", "fmt", generatePath).Run()
}
