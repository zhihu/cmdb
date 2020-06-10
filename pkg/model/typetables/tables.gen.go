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

package typetables

import (
	cdc "github.com/zhihu/cmdb/pkg/storage/cdc"

	model "github.com/zhihu/cmdb/pkg/model"
)

type RichObjectMeta struct {
	model.ObjectMeta
	ObjectType *RichObjectType
}

type ObjectMetaTable struct {
	ID       map[IndexObjectMetaID]*RichObjectMeta
	TypeID   map[IndexObjectMetaTypeID][]*RichObjectMeta
	TypeName map[IndexObjectMetaTypeName]*RichObjectMeta
}

func (t *ObjectMetaTable) Init() {
	t.ID = map[IndexObjectMetaID]*RichObjectMeta{}
	t.TypeID = map[IndexObjectMetaTypeID][]*RichObjectMeta{}
	t.TypeName = map[IndexObjectMetaTypeName]*RichObjectMeta{}

}

type IndexObjectMetaID struct {
	ID int
}

type IndexObjectMetaTypeID struct {
	TypeID int
}

type IndexObjectMetaTypeName struct {
	TypeID int
	Name   string
}

func (t *ObjectMetaTable) GetByID(ID int) (row *RichObjectMeta, ok bool) {
	row, ok = t.ID[IndexObjectMetaID{ID: ID}]
	return
}
func (t *ObjectMetaTable) FilterByTypeID(TypeID int) (rows []*RichObjectMeta) {
	rows, _ = t.TypeID[IndexObjectMetaTypeID{TypeID: TypeID}]
	return
}

func (t *ObjectMetaTable) GetByTypeName(TypeID int, Name string) (row *RichObjectMeta, ok bool) {
	row, ok = t.TypeName[IndexObjectMetaTypeName{TypeID: TypeID, Name: Name}]
	return
}

type RichObjectState struct {
	model.ObjectState
	ObjectStatus *RichObjectStatus
}

type ObjectStateTable struct {
	ID           map[IndexObjectStateID]*RichObjectState
	StatusID     map[IndexObjectStateStatusID][]*RichObjectState
	StatusIDName map[IndexObjectStateStatusIDName]*RichObjectState
}

func (t *ObjectStateTable) Init() {
	t.ID = map[IndexObjectStateID]*RichObjectState{}
	t.StatusID = map[IndexObjectStateStatusID][]*RichObjectState{}
	t.StatusIDName = map[IndexObjectStateStatusIDName]*RichObjectState{}

}

type IndexObjectStateID struct {
	ID int
}

type IndexObjectStateStatusID struct {
	StatusID int
}

type IndexObjectStateStatusIDName struct {
	StatusID int
	Name     string
}

func (t *ObjectStateTable) GetByID(ID int) (row *RichObjectState, ok bool) {
	row, ok = t.ID[IndexObjectStateID{ID: ID}]
	return
}
func (t *ObjectStateTable) FilterByStatusID(StatusID int) (rows []*RichObjectState) {
	rows, _ = t.StatusID[IndexObjectStateStatusID{StatusID: StatusID}]
	return
}

func (t *ObjectStateTable) GetByStatusIDName(StatusID int, Name string) (row *RichObjectState, ok bool) {
	row, ok = t.StatusIDName[IndexObjectStateStatusIDName{StatusID: StatusID, Name: Name}]
	return
}

type RichObjectStatus struct {
	model.ObjectStatus
	ObjectType *RichObjectType
	// has_many
	ObjectState map[IndexObjectStateID]*RichObjectState
}

type ObjectStatusTable struct {
	ID         map[IndexObjectStatusID]*RichObjectStatus
	TypeID     map[IndexObjectStatusTypeID][]*RichObjectStatus
	TypeIDName map[IndexObjectStatusTypeIDName]*RichObjectStatus
}

func (t *ObjectStatusTable) Init() {
	t.ID = map[IndexObjectStatusID]*RichObjectStatus{}
	t.TypeID = map[IndexObjectStatusTypeID][]*RichObjectStatus{}
	t.TypeIDName = map[IndexObjectStatusTypeIDName]*RichObjectStatus{}

}

type IndexObjectStatusID struct {
	ID int
}

type IndexObjectStatusTypeID struct {
	TypeID int
}

type IndexObjectStatusTypeIDName struct {
	TypeID int
	Name   string
}

func (t *ObjectStatusTable) GetByID(ID int) (row *RichObjectStatus, ok bool) {
	row, ok = t.ID[IndexObjectStatusID{ID: ID}]
	return
}
func (t *ObjectStatusTable) FilterByTypeID(TypeID int) (rows []*RichObjectStatus) {
	rows, _ = t.TypeID[IndexObjectStatusTypeID{TypeID: TypeID}]
	return
}

func (t *ObjectStatusTable) GetByTypeIDName(TypeID int, Name string) (row *RichObjectStatus, ok bool) {
	row, ok = t.TypeIDName[IndexObjectStatusTypeIDName{TypeID: TypeID, Name: Name}]
	return
}

type RichObjectType struct {
	model.ObjectType
	// has_many
	ObjectStatus map[IndexObjectStatusID]*RichObjectStatus
	// has_many
	ObjectMeta map[IndexObjectMetaID]*RichObjectMeta
}

type ObjectTypeTable struct {
	ID   map[IndexObjectTypeID]*RichObjectType
	Name map[IndexObjectTypeName]*RichObjectType
}

func (t *ObjectTypeTable) Init() {
	t.ID = map[IndexObjectTypeID]*RichObjectType{}
	t.Name = map[IndexObjectTypeName]*RichObjectType{}

}

type IndexObjectTypeID struct {
	ID int
}

type IndexObjectTypeName struct {
	Name string
}

func (t *ObjectTypeTable) GetByID(ID int) (row *RichObjectType, ok bool) {
	row, ok = t.ID[IndexObjectTypeID{ID: ID}]
	return
}
func (t *ObjectTypeTable) GetByName(Name string) (row *RichObjectType, ok bool) {
	row, ok = t.Name[IndexObjectTypeName{Name: Name}]
	return
}

type Database struct {
	ObjectMetaTable ObjectMetaTable

	ObjectStateTable ObjectStateTable

	ObjectStatusTable ObjectStatusTable

	ObjectTypeTable ObjectTypeTable
}

func (d *Database) Init() {
	d.ObjectMetaTable.Init()
	d.ObjectStateTable.Init()
	d.ObjectStatusTable.Init()
	d.ObjectTypeTable.Init()
}

func (d *Database) InsertObjectMeta(row *model.ObjectMeta) (ok bool) {
	if row.DeleteTime != nil {
		return false
	}
	var richRow = &RichObjectMeta{
		ObjectMeta: *row,
	}

	{
		var index = IndexObjectMetaID{row.ID}
		_, ok := d.ObjectMetaTable.ID[index]
		if ok {
			return false
		}
		d.ObjectMetaTable.ID[index] = richRow

	}

	{
		var index = IndexObjectMetaTypeID{row.TypeID}

		list := d.ObjectMetaTable.TypeID[index]
		list = append(list, richRow)
		d.ObjectMetaTable.TypeID[index] = list
	}

	{
		var index = IndexObjectMetaTypeName{row.TypeID, row.Name}
		_, ok := d.ObjectMetaTable.TypeName[index]
		if ok {
			return false
		}
		d.ObjectMetaTable.TypeName[index] = richRow

	}

	{ //belongs_to
		richRow.ObjectType, _ = d.ObjectTypeTable.GetByID(row.TypeID)
		if richRow.ObjectType != nil {
			richRow.ObjectType.ObjectMeta[IndexObjectMetaID{row.ID}] = richRow
		}
	}

	return true
}

func (d *Database) UpdateObjectMeta(row *model.ObjectMeta) (ok bool) {
	if row.DeleteTime != nil {
		return d.DeleteObjectMeta(row)
	}
	var index = IndexObjectMetaID{row.ID}
	origin, ok := d.ObjectMetaTable.ID[index]
	if !ok {
		return d.InsertObjectMeta(row)
	}
	origin.ObjectMeta = *row
	return true
}

func (d *Database) DeleteObjectMeta(row *model.ObjectMeta) (ok bool) {
	var index = IndexObjectMetaID{row.ID}
	richRow, ok := d.ObjectMetaTable.ID[index]
	if !ok {
		return false
	}

	{
		var index = IndexObjectMetaID{row.ID}
		delete(d.ObjectMetaTable.ID, index)

	}

	{
		var index = IndexObjectMetaTypeID{row.TypeID}

		list := d.ObjectMetaTable.TypeID[index]
		var newList = make([]*RichObjectMeta, 0, len(list)-1)
		for _, item := range list {
			if item.ObjectMeta.ID != row.ID {
				newList = append(newList, item)
			}
		}
		d.ObjectMetaTable.TypeID[index] = newList
	}

	{
		var index = IndexObjectMetaTypeName{row.TypeID, row.Name}
		delete(d.ObjectMetaTable.TypeName, index)

	}

	{
		if richRow.ObjectType != nil {
			delete(richRow.ObjectType.ObjectMeta, IndexObjectMetaID{row.ID})
		}
	}

	return true
}
func (d *Database) InsertObjectState(row *model.ObjectState) (ok bool) {
	if row.DeleteTime != nil {
		return false
	}
	var richRow = &RichObjectState{
		ObjectState: *row,
	}

	{
		var index = IndexObjectStateID{row.ID}
		_, ok := d.ObjectStateTable.ID[index]
		if ok {
			return false
		}
		d.ObjectStateTable.ID[index] = richRow

	}

	{
		var index = IndexObjectStateStatusID{row.StatusID}

		list := d.ObjectStateTable.StatusID[index]
		list = append(list, richRow)
		d.ObjectStateTable.StatusID[index] = list
	}

	{
		var index = IndexObjectStateStatusIDName{row.StatusID, row.Name}
		_, ok := d.ObjectStateTable.StatusIDName[index]
		if ok {
			return false
		}
		d.ObjectStateTable.StatusIDName[index] = richRow

	}

	{ //belongs_to
		richRow.ObjectStatus, _ = d.ObjectStatusTable.GetByID(row.StatusID)
		if richRow.ObjectStatus != nil {
			richRow.ObjectStatus.ObjectState[IndexObjectStateID{row.ID}] = richRow
		}
	}

	return true
}

func (d *Database) UpdateObjectState(row *model.ObjectState) (ok bool) {
	if row.DeleteTime != nil {
		return d.DeleteObjectState(row)
	}
	var index = IndexObjectStateID{row.ID}
	origin, ok := d.ObjectStateTable.ID[index]
	if !ok {
		return d.InsertObjectState(row)
	}
	origin.ObjectState = *row
	return true
}

func (d *Database) DeleteObjectState(row *model.ObjectState) (ok bool) {
	var index = IndexObjectStateID{row.ID}
	richRow, ok := d.ObjectStateTable.ID[index]
	if !ok {
		return false
	}

	{
		var index = IndexObjectStateID{row.ID}
		delete(d.ObjectStateTable.ID, index)

	}

	{
		var index = IndexObjectStateStatusID{row.StatusID}

		list := d.ObjectStateTable.StatusID[index]
		var newList = make([]*RichObjectState, 0, len(list)-1)
		for _, item := range list {
			if item.ObjectState.ID != row.ID {
				newList = append(newList, item)
			}
		}
		d.ObjectStateTable.StatusID[index] = newList
	}

	{
		var index = IndexObjectStateStatusIDName{row.StatusID, row.Name}
		delete(d.ObjectStateTable.StatusIDName, index)

	}

	{
		if richRow.ObjectStatus != nil {
			delete(richRow.ObjectStatus.ObjectState, IndexObjectStateID{row.ID})
		}
	}

	return true
}
func (d *Database) InsertObjectStatus(row *model.ObjectStatus) (ok bool) {
	if row.DeleteTime != nil {
		return false
	}
	var richRow = &RichObjectStatus{
		ObjectStatus: *row,
		ObjectState:  map[IndexObjectStateID]*RichObjectState{},
	}

	{
		var index = IndexObjectStatusID{row.ID}
		_, ok := d.ObjectStatusTable.ID[index]
		if ok {
			return false
		}
		d.ObjectStatusTable.ID[index] = richRow

	}

	{
		var index = IndexObjectStatusTypeID{row.TypeID}

		list := d.ObjectStatusTable.TypeID[index]
		list = append(list, richRow)
		d.ObjectStatusTable.TypeID[index] = list
	}

	{
		var index = IndexObjectStatusTypeIDName{row.TypeID, row.Name}
		_, ok := d.ObjectStatusTable.TypeIDName[index]
		if ok {
			return false
		}
		d.ObjectStatusTable.TypeIDName[index] = richRow

	}

	{ //belongs_to
		richRow.ObjectType, _ = d.ObjectTypeTable.GetByID(row.TypeID)
		if richRow.ObjectType != nil {
			richRow.ObjectType.ObjectStatus[IndexObjectStatusID{row.ID}] = richRow
		}
	}

	{ // has_many
		richRow.ObjectState = map[IndexObjectStateID]*RichObjectState{}
		var list = d.ObjectStateTable.FilterByStatusID(row.ID)
		for _, item := range list {
			item.ObjectStatus = richRow
			richRow.ObjectState[IndexObjectStateID{ID: item.ID}] = item
		}
	}

	return true
}

func (d *Database) UpdateObjectStatus(row *model.ObjectStatus) (ok bool) {
	if row.DeleteTime != nil {
		return d.DeleteObjectStatus(row)
	}
	var index = IndexObjectStatusID{row.ID}
	origin, ok := d.ObjectStatusTable.ID[index]
	if !ok {
		return d.InsertObjectStatus(row)
	}
	origin.ObjectStatus = *row
	return true
}

func (d *Database) DeleteObjectStatus(row *model.ObjectStatus) (ok bool) {
	var index = IndexObjectStatusID{row.ID}
	richRow, ok := d.ObjectStatusTable.ID[index]
	if !ok {
		return false
	}

	{
		var index = IndexObjectStatusID{row.ID}
		delete(d.ObjectStatusTable.ID, index)

	}

	{
		var index = IndexObjectStatusTypeID{row.TypeID}

		list := d.ObjectStatusTable.TypeID[index]
		var newList = make([]*RichObjectStatus, 0, len(list)-1)
		for _, item := range list {
			if item.ObjectStatus.ID != row.ID {
				newList = append(newList, item)
			}
		}
		d.ObjectStatusTable.TypeID[index] = newList
	}

	{
		var index = IndexObjectStatusTypeIDName{row.TypeID, row.Name}
		delete(d.ObjectStatusTable.TypeIDName, index)

	}

	{
		if richRow.ObjectType != nil {
			delete(richRow.ObjectType.ObjectStatus, IndexObjectStatusID{row.ID})
		}
	}

	{
		for _, item := range richRow.ObjectState {
			item.ObjectStatus = nil
		}
	}

	return true
}
func (d *Database) InsertObjectType(row *model.ObjectType) (ok bool) {
	if row.DeleteTime != nil {
		return false
	}
	var richRow = &RichObjectType{
		ObjectType:   *row,
		ObjectStatus: map[IndexObjectStatusID]*RichObjectStatus{},
		ObjectMeta:   map[IndexObjectMetaID]*RichObjectMeta{},
	}

	{
		var index = IndexObjectTypeID{row.ID}
		_, ok := d.ObjectTypeTable.ID[index]
		if ok {
			return false
		}
		d.ObjectTypeTable.ID[index] = richRow

	}

	{
		var index = IndexObjectTypeName{row.Name}
		_, ok := d.ObjectTypeTable.Name[index]
		if ok {
			return false
		}
		d.ObjectTypeTable.Name[index] = richRow

	}

	{ // has_many
		richRow.ObjectStatus = map[IndexObjectStatusID]*RichObjectStatus{}
		var list = d.ObjectStatusTable.FilterByTypeID(row.ID)
		for _, item := range list {
			item.ObjectType = richRow
			richRow.ObjectStatus[IndexObjectStatusID{ID: item.ID}] = item
		}
	}

	{ // has_many
		richRow.ObjectMeta = map[IndexObjectMetaID]*RichObjectMeta{}
		var list = d.ObjectMetaTable.FilterByTypeID(row.ID)
		for _, item := range list {
			item.ObjectType = richRow
			richRow.ObjectMeta[IndexObjectMetaID{ID: item.ID}] = item
		}
	}

	return true
}

func (d *Database) UpdateObjectType(row *model.ObjectType) (ok bool) {
	if row.DeleteTime != nil {
		return d.DeleteObjectType(row)
	}
	var index = IndexObjectTypeID{row.ID}
	origin, ok := d.ObjectTypeTable.ID[index]
	if !ok {
		return d.InsertObjectType(row)
	}
	origin.ObjectType = *row
	return true
}

func (d *Database) DeleteObjectType(row *model.ObjectType) (ok bool) {
	var index = IndexObjectTypeID{row.ID}
	richRow, ok := d.ObjectTypeTable.ID[index]
	if !ok {
		return false
	}

	{
		var index = IndexObjectTypeID{row.ID}
		delete(d.ObjectTypeTable.ID, index)

	}

	{
		var index = IndexObjectTypeName{row.Name}
		delete(d.ObjectTypeTable.Name, index)

	}

	{
		for _, item := range richRow.ObjectStatus {
			item.ObjectType = nil
		}
	}

	{
		for _, item := range richRow.ObjectMeta {
			item.ObjectType = nil
		}
	}

	return true
}

func (d *Database) OnEvents(transaction []cdc.Event) {
	for _, event := range transaction {
		switch row := event.Row.(type) {
		case *model.ObjectMeta:
			switch event.Type {
			case cdc.Create:
				d.InsertObjectMeta(row)
			case cdc.Update:
				d.UpdateObjectMeta(row)
			case cdc.Delete:
				d.DeleteObjectMeta(row)
			}
		case *model.ObjectState:
			switch event.Type {
			case cdc.Create:
				d.InsertObjectState(row)
			case cdc.Update:
				d.UpdateObjectState(row)
			case cdc.Delete:
				d.DeleteObjectState(row)
			}
		case *model.ObjectStatus:
			switch event.Type {
			case cdc.Create:
				d.InsertObjectStatus(row)
			case cdc.Update:
				d.UpdateObjectStatus(row)
			case cdc.Delete:
				d.DeleteObjectStatus(row)
			}
		case *model.ObjectType:
			switch event.Type {
			case cdc.Create:
				d.InsertObjectType(row)
			case cdc.Update:
				d.UpdateObjectType(row)
			case cdc.Delete:
				d.DeleteObjectType(row)
			}
		}
	}
}
