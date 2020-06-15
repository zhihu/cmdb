{{range $tableName,$table:= $.Tables}}
	type Rich{{$table.Name}} struct{
    {{TypeName $table.Type}}
    {{- range $i,$relation := $table.Relations -}}
        {{- if eq $relation.Type 0 -}}{{/* has_one */}}
        {{$relation.ToTable.Name}} *Rich{{$relation.ToTable.Name}}
        {{- end -}}
        {{- if eq $relation.Type 1 -}}{{/* has_many */}}
		// has_many
        {{$relation.ToTable.Name}} map[Index{{$relation.ToTable.Name}}{{$relation.ToTable.ID.Name}}]*Rich{{$relation.ToTable.Name}}
        {{- end -}}
        {{- if eq $relation.Type 2 -}}{{/* belongs_to */}}
        {{$relation.ToTable.Name}} *Rich{{$relation.ToTable.Name}}
        {{- end -}}
    {{- end -}}
	}

	type {{ $table.Name }}Table struct{
    {{- range $name,$index :=  $table.Indexes -}}
        {{$name}} map[Index{{$table.Name}}{{$name}}]{{if $index.Unique}}*Rich{{$table.Name}}{{else}}[]*Rich{{$table.Name}} {{end}}
    {{end}}
	}

	func (t *{{$table.Name}}Table)Init(){
    {{- range $name,$index :=  $table.Indexes -}}
		t.{{$name}} = map[Index{{$table.Name}}{{$name}}]{{if $index.Unique}}*Rich{{$table.Name}}{{else}}[]*Rich{{$table.Name}} {{end}}{}
    {{end}}
	}

    {{range $name,$index :=  $table.Indexes}}
		type Index{{$table.Name}}{{$name}} struct{
        {{- range $i,$t := $index.Fields -}}
            {{$t.Name}} {{TypeName $t.Type}}
        {{end}}
		}
    {{end}}

    {{- range $name,$index :=  $table.Indexes -}}
        {{- if $index.Unique -}}
			func (t *{{$table.Name}}Table)GetBy{{$name}}({{range $i,$f:= $index.Fields}}  {{if ne $i 0}},{{end}} {{$f.Name}} {{TypeName $f.Type}} {{end}})(row *Rich{{$table.Name}},ok bool){
			row, ok = t.{{$name}}[Index{{$table.Name}}{{$name}} {
            {{- range $i,$f:= $index.Fields -}}
                {{$f.Name}}:{{$f.Name}},
            {{- end}}
			}]
			return
			}
        {{- else -}}
			func (t *{{$table.Name}}Table)FilterBy{{$name}}({{range $i,$f:= $index.Fields}}  {{if ne $i 0}},{{end}} {{$f.Name}} {{TypeName $f.Type}} {{end}})(rows []*Rich{{$table.Name}}){
			rows, _ = t.{{$name}}[Index{{$table.Name}}{{$name}} {
            {{- range $i,$f:= $index.Fields -}}
                {{$f.Name}}:{{$f.Name}},
            {{- end}}
			}]
			return
			}
        {{end}}
    {{end}}
{{end}}



type Database struct{
{{- range $tableName,$table:= $.Tables}}
    {{$tableName}}Table {{$tableName}}Table
{{end -}}
}

func(d *Database)Init(){
{{- range $tableName,$table:= $.Tables}}
	d. {{$tableName}}Table.Init()
{{- end}}
}

{{range $tableName,$table:= $.Tables}}
	func(d *Database)Insert{{$tableName}}(row *{{TypeName $table.Type}})(ok bool){
    {{if $table.SoftDelete -}}
		if row.{{$table.SoftDelete.Name}} != nil{
		return false
		}
    {{end -}}
	var richRow = &Rich{{$tableName}}{
    {{$table.Type.Name}}:*row ,
    {{- range $i,$relation := $table.Relations -}}
        {{- if eq $relation.Type 1 -}}{{/* has_many */}}
        {{$relation.ToTable.Name}} : map[Index{{$relation.ToTable.Name}}{{$relation.ToTable.ID.Name}}]*Rich{{$relation.ToTable.Name}}{},
        {{- end -}}
    {{- end }}
	}
    {{/* insert to all index:*/}}
    {{- range $name,$index :=  $table.Indexes -}}
        {{/** block to avoid variable names conflict */}}
		{
		var index = Index{{$table.Name}}{{$name}} {
        {{- range $i,$f:= $index.Fields -}}
			row.{{$f.Name}},
        {{- end -}}
		};
        {{if $index.Unique -}}
			_,ok := d.{{$tableName}}Table.{{$name}}[index]
			if ok{
			return false
			}
			d.{{$tableName}}Table.{{$name}}[index] = richRow
        {{else}}
			list := d.{{$tableName}}Table.{{$name}}[index]
			list = append(list,richRow)
			d.{{$tableName}}Table.{{$name}}[index] = list
        {{- end}}
		}
    {{end}}
    {{/* insert to all relations:*/}}
    {{range $i,$relation:= $table.Relations}}
		{
        {{- if eq $relation.Type 0 -}}{{/* has_one */}}
		richRow.{{$relation.ToTable.Name}} , ok :=  d.{{$relation.ToTable.Name}}Table.GetBy{{$relation.To.Name}}(row.{{$relation.From.Name}})
		if ok{
		richRow.{{$relation.ToTable.Name}}.{{$tableName}} = richRow
		}
        {{- end -}}
        {{- if eq $relation.Type 1 -}}
			// has_many
			richRow.{{$relation.ToTable.Name}} =  map[Index{{$relation.ToTable.Name}}{{$relation.ToTable.ID.Name}}]*Rich{{$relation.ToTable.Name}}{}
			var list =d.{{$relation.ToTable.Name}}Table.FilterBy{{$relation.To.Name}}(row.{{$relation.From.Name}})
			for _,item:=range list{
			item.{{$tableName}} = richRow
			richRow.{{$relation.ToTable.Name}}[Index{{$relation.ToTable.Name}}{{$relation.ToTable.ID.Name}}{
            {{- range $k,$f:= $relation.ToTable.ID.Fields -}}
                {{$f.Name}} : item.{{$f.Name}},
            {{- end -}}
			}] = item
			}
        {{- end -}}
        {{- if eq $relation.Type 2 -}}
			//belongs_to
			richRow.{{$relation.ToTable.Name}} ,_ =  d.{{$relation.ToTable.Name}}Table.GetBy{{$relation.To.Name}}(row.{{$relation.From.Name}})
			if richRow.{{$relation.ToTable.Name}}!=nil{
            {{- $reverse := ($relation.ToTable.GetReverseRelation $relation) -}}
            {{if eq $reverse.Type 0 -}}
				richRow.{{$relation.ToTable.Name}}.{{$tableName}} = richRow
            {{- end}}
            {{- if eq $reverse.Type 1 -}}
				richRow.{{$relation.ToTable.Name}}.{{$tableName}}[Index{{$table.Name}}{{$table.ID.Name}} {
                {{- range $i,$f:= $table.ID.Fields -}}
					row.{{$f.Name}},
                {{- end -}} }] = richRow
            {{- end}}
			}
        {{- end -}}
		}
    {{end}}
	return true
	}

    {{$idIndex :=  $table.ID -}}
	func(d *Database)Update{{$tableName}}(row *{{TypeName $table.Type}})(ok bool){
    {{if $table.SoftDelete -}}
		if row.{{$table.SoftDelete.Name}} != nil{
		return d.Delete{{$tableName}}(row)
		}
    {{end -}}
	var index = Index{{$table.Name}}{{$idIndex.Name}} {
    {{- range $i,$f:= $idIndex.Fields -}}
		row.{{$f.Name}},
    {{- end -}}
	};
	origin,ok := d.{{$tableName}}Table.{{$idIndex.Name}}[index]
	if !ok{
	return d.Insert{{$tableName}}(row)
	}
	origin.{{$table.Type.Name}} = *row
	return true
	}

	func(d *Database)Delete{{$tableName}}(row *{{TypeName $table.Type}})(ok bool){
	var index = Index{{$table.Name}}{{$idIndex.Name}} {
    {{- range $i,$f:= $idIndex.Fields -}}
		row.{{$f.Name}},
    {{- end -}}
	};
	richRow,ok := d.{{$tableName}}Table.{{$idIndex.Name}}[index]
	if !ok{
	return false
	}
    {{/* delete from all index:*/}}
    {{- range $name,$index :=  $table.Indexes -}}
        {{/** block to avoid variable names conflict */}}
		{
		var index = Index{{$table.Name}}{{$name}} {
        {{- range $i,$f:= $index.Fields -}}
			row.{{$f.Name}},
        {{- end -}}
		};
        {{if $index.Unique -}}
			delete(d.{{$tableName}}Table.{{$name}},index)
        {{else}}
			list := d.{{$tableName}}Table.{{$name}}[index]
			var newList = make([]*Rich{{$tableName}},0,len(list) - 1)
			for _,item := range list{
			if {{range $k,$f := $idIndex.Fields}}{{if ne $k 0}} || {{end}} item.{{$table.Type.Name}}.{{$f.Name}} != row.{{$f.Name}} {{end}}{
			newList = append(newList,item)
			}
			}
			d.{{$tableName}}Table.{{$name}}[index] = newList
        {{- end}}
		}
    {{end}}

    {{range $i,$relation:= $table.Relations}}
		{
        {{- if eq $relation.Type 0 -}}{{/* has_one */}}
		if richRow.{{$relation.ToTable.Name}} !=nil{
		richRow.{{$relation.ToTable.Name}}.{{$tableName}} = nil
		}
        {{- end -}}
        {{- if eq $relation.Type 1 -}}
			for _,item:=range richRow.{{$relation.ToTable.Name}} {
			item.{{$tableName}} = nil
			}
        {{- end -}}
        {{- if eq $relation.Type 2 -}}
			if richRow.{{$relation.ToTable.Name}}!=nil{
            {{$reverse := ($relation.ToTable.GetReverseRelation $relation) -}}
            {{if eq $reverse.Type 0 -}}
				richRow.{{$relation.ToTable.Name}}.{{$tableName}} = nil
            {{- end}}
            {{- if eq $reverse.Type 1 -}}
				delete(richRow.{{$relation.ToTable.Name}}.{{$tableName}},Index{{$table.Name}}{{$table.ID.Name}} {
                {{- range $i,$f:= $table.ID.Fields -}}
					row.{{$f.Name}},
                {{- end -}} })
            {{- end}}
			}
        {{- end -}}
		}
    {{end}}
	return true
	}

{{- end}}