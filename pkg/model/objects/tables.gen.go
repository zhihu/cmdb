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

// DO NOT EDIT IT.

package objects

import (
	cdc "github.com/zhihu/cmdb/pkg/storage/cdc"

	model "github.com/zhihu/cmdb/pkg/model"
)

type RichObject struct {
	model.Object
	// has_many
	ObjectMetaValue map[IndexObjectMetaValueID]*RichObjectMetaValue
}

type ObjectTable struct {
	ID map[IndexObjectID]*RichObject
}

func (t *ObjectTable) Init() {
	t.ID = map[IndexObjectID]*RichObject{}

}

type IndexObjectID struct {
	ID int
}

func (t *ObjectTable) GetByID(ID int) (row *RichObject, ok bool) {
	row, ok = t.ID[IndexObjectID{ID: ID}]
	return
}

type RichObjectMetaValue struct {
	model.ObjectMetaValue
	Object *RichObject
}

type ObjectMetaValueTable struct {
	ID       map[IndexObjectMetaValueID]*RichObjectMetaValue
	ObjectID map[IndexObjectMetaValueObjectID][]*RichObjectMetaValue
}

func (t *ObjectMetaValueTable) Init() {
	t.ID = map[IndexObjectMetaValueID]*RichObjectMetaValue{}
	t.ObjectID = map[IndexObjectMetaValueObjectID][]*RichObjectMetaValue{}

}

type IndexObjectMetaValueID struct {
	ObjectID int
	MetaID   int
}

type IndexObjectMetaValueObjectID struct {
	ObjectID int
}

func (t *ObjectMetaValueTable) GetByID(ObjectID int, MetaID int) (row *RichObjectMetaValue, ok bool) {
	row, ok = t.ID[IndexObjectMetaValueID{ObjectID: ObjectID, MetaID: MetaID}]
	return
}
func (t *ObjectMetaValueTable) FilterByObjectID(ObjectID int) (rows []*RichObjectMetaValue) {
	rows, _ = t.ObjectID[IndexObjectMetaValueObjectID{ObjectID: ObjectID}]
	return
}

type Database struct {
	ObjectTable ObjectTable

	ObjectMetaValueTable ObjectMetaValueTable
}

func (d *Database) Init() {
	d.ObjectTable.Init()
	d.ObjectMetaValueTable.Init()
}

func (d *Database) InsertObject(row *model.Object) (ok bool) {
	if row.DeleteTime != nil {
		return false
	}
	var richRow = &RichObject{
		Object:          *row,
		ObjectMetaValue: map[IndexObjectMetaValueID]*RichObjectMetaValue{},
	}

	{
		var index = IndexObjectID{row.ID}
		_, ok := d.ObjectTable.ID[index]
		if ok {
			return false
		}
		d.ObjectTable.ID[index] = richRow

	}

	{ // has_many
		richRow.ObjectMetaValue = map[IndexObjectMetaValueID]*RichObjectMetaValue{}
		var list = d.ObjectMetaValueTable.FilterByObjectID(row.ID)
		for _, item := range list {
			item.Object = richRow
			richRow.ObjectMetaValue[IndexObjectMetaValueID{ObjectID: item.ObjectID, MetaID: item.MetaID}] = item
		}
	}

	return true
}

func (d *Database) UpdateObject(row *model.Object) (ok bool) {
	if row.DeleteTime != nil {
		return d.DeleteObject(row)
	}
	var index = IndexObjectID{row.ID}
	origin, ok := d.ObjectTable.ID[index]
	if !ok {
		return d.InsertObject(row)
	}
	origin.Object = *row
	return true
}

func (d *Database) DeleteObject(row *model.Object) (ok bool) {
	var index = IndexObjectID{row.ID}
	richRow, ok := d.ObjectTable.ID[index]
	if !ok {
		return false
	}

	{
		var index = IndexObjectID{row.ID}
		delete(d.ObjectTable.ID, index)

	}

	{
		for _, item := range richRow.ObjectMetaValue {
			item.Object = nil
		}
	}

	return true
}
func (d *Database) InsertObjectMetaValue(row *model.ObjectMetaValue) (ok bool) {
	if row.DeleteTime != nil {
		return false
	}
	var richRow = &RichObjectMetaValue{
		ObjectMetaValue: *row,
	}

	{
		var index = IndexObjectMetaValueID{row.ObjectID, row.MetaID}
		_, ok := d.ObjectMetaValueTable.ID[index]
		if ok {
			return false
		}
		d.ObjectMetaValueTable.ID[index] = richRow

	}

	{
		var index = IndexObjectMetaValueObjectID{row.ObjectID}

		list := d.ObjectMetaValueTable.ObjectID[index]
		list = append(list, richRow)
		d.ObjectMetaValueTable.ObjectID[index] = list
	}

	{ //belongs_to
		richRow.Object, _ = d.ObjectTable.GetByID(row.ObjectID)
		if richRow.Object != nil {
			richRow.Object.ObjectMetaValue[IndexObjectMetaValueID{row.ObjectID, row.MetaID}] = richRow
		}
	}

	return true
}

func (d *Database) UpdateObjectMetaValue(row *model.ObjectMetaValue) (ok bool) {
	if row.DeleteTime != nil {
		return d.DeleteObjectMetaValue(row)
	}
	var index = IndexObjectMetaValueID{row.ObjectID, row.MetaID}
	origin, ok := d.ObjectMetaValueTable.ID[index]
	if !ok {
		return d.InsertObjectMetaValue(row)
	}
	origin.ObjectMetaValue = *row
	return true
}

func (d *Database) DeleteObjectMetaValue(row *model.ObjectMetaValue) (ok bool) {
	var index = IndexObjectMetaValueID{row.ObjectID, row.MetaID}
	richRow, ok := d.ObjectMetaValueTable.ID[index]
	if !ok {
		return false
	}

	{
		var index = IndexObjectMetaValueID{row.ObjectID, row.MetaID}
		delete(d.ObjectMetaValueTable.ID, index)

	}

	{
		var index = IndexObjectMetaValueObjectID{row.ObjectID}

		list := d.ObjectMetaValueTable.ObjectID[index]
		var newList = make([]*RichObjectMetaValue, 0, len(list)-1)
		for _, item := range list {
			if item.ObjectMetaValue.ObjectID != row.ObjectID || item.ObjectMetaValue.MetaID != row.MetaID {
				newList = append(newList, item)
			}
		}
		d.ObjectMetaValueTable.ObjectID[index] = newList
	}

	{
		if richRow.Object != nil {
			delete(richRow.Object.ObjectMetaValue, IndexObjectMetaValueID{row.ObjectID, row.MetaID})
		}
	}

	return true
}

func (d *Database) OnEvents(transaction []cdc.Event) {
	for _, event := range transaction {
		switch row := event.Row.(type) {
		case *model.Object:
			switch event.Type {
			case cdc.Create:
				d.InsertObject(row)
			case cdc.Update:
				d.UpdateObject(row)
			case cdc.Delete:
				d.DeleteObject(row)
			}
		case *model.ObjectMetaValue:
			switch event.Type {
			case cdc.Create:
				d.InsertObjectMetaValue(row)
			case cdc.Update:
				d.UpdateObjectMetaValue(row)
			case cdc.Delete:
				d.DeleteObjectMetaValue(row)
			}
		}
	}
}
