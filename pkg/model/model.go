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

package model

import (
	"strconv"
	"time"
)

const (
	STRING  = 0
	INTEGER = 1
	DOUBLE  = 2
	BOOLEAN = 3
)

func BooleanValue(b bool) string {
	return strconv.FormatBool(b)
}

type ObjectType struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type ObjectStatus struct {
	ID          int    `db:"id"`
	TypeID      int    `db:"type_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	DeleteTime *time.Time `db:"delete_time"`
	States     map[string]*ObjectState
}

type ObjectState struct {
	ID          int    `db:"id"`
	StatusID    int    `db:"status_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type Object struct {
	ID              int    `db:"id"`
	TypeID          int    `db:"type_id"`
	Name            string `db:"name"`
	Version         int    `db:"version"`
	RelationVersion int    `db:"relation_version"`
	Description     string `db:"description"`
	StatusID        int    `db:"status_id"`
	StateID         int    `db:"state_id"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type DeletedObject struct {
	ID              int    `db:"id"`
	TypeID          int    `db:"type_id"`
	Name            string `db:"name"`
	Version         int    `db:"version"`
	RelationVersion int    `db:"relation_version"`
	Description     string `db:"description"`
	StatusID        int    `db:"status_id"`
	StateID         int    `db:"state_id"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type ObjectMeta struct {
	ID     int    `db:"id"`
	TypeID int    `db:"type_id"`
	Name   string `db:"name"`
	// Comment: 1: STRING 2: INTEGER, 3: DOUBLE, 4: BOOLEAN
	// Default: 1
	ValueType   int    `db:"value_type"`
	Description string `db:"description"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type ObjectMetaValue struct {
	ObjectID int    `db:"object_id"`
	MetaID   int    `db:"meta_id"`
	Value    string `db:"value"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type DeletedObjectMetaValue struct {
	ObjectID int    `db:"object_id"`
	MetaID   int    `db:"meta_id"`
	Value    string `db:"value"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type ObjectLog struct {
	ID       int `db:"id"`
	ObjectID int `db:"object_id"`
	// Comment: 0: EMERGENCY 1: ALERT 2: CRITICAL 3: ERROR 4: WARNING 5: NOTICE 6: INFORMATIONAL 7: DEBUG 8: NOTE
	Level int `db:"level"`
	// Comment: 0: text/plain 1: application/json
	Format int `db:"format"`
	// Comment: 0: INTERNAL 1: API 2: USER 3: SYSTEM
	Source   int    `db:"source"`
	Message  string `db:"message"`
	CreateBy string `db:"create_by"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type DeletedObjectLog struct {
	ID       int `db:"id"`
	ObjectID int `db:"object_id"`
	// Comment: 0: EMERGENCY 1: ALERT 2: CRITICAL 3: ERROR 4: WARNING 5: NOTICE 6: INFORMATIONAL 7: DEBUG 8: NOTE
	Level int `db:"level"`
	// Comment: 0: text/plain 1: application/json
	Format int `db:"format"`
	// Comment: 0: INTERNAL 1: API 2: USER 3: SYSTEM
	Source   int    `db:"source"`
	Message  string `db:"message"`
	CreateBy string `db:"create_by"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type ObjectRelationType struct {
	ID          int    `db:"id"`
	FromTypeID  int    `db:"from_type_id"`
	ToTypeID    int    `db:"to_type_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type DeletedObjectRelationType struct {
	ID          int    `db:"id"`
	FromTypeID  int    `db:"from_type_id"`
	ToTypeID    int    `db:"to_type_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type ObjectRelation struct {
	FromObjectID   int `db:"from_object_id"`
	RelationTypeID int `db:"relation_type_id"`
	ToObjectID     int `db:"to_object_id"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time `db:"create_time"`
	// Default: CURRENT_TIMESTAMP
	UpdateTime time.Time  `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type DeletedObjectRelation struct {
	FromObjectID   int `db:"from_object_id"`
	RelationTypeID int `db:"relation_type_id"`
	ToObjectID     int `db:"to_object_id"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time `db:"create_time"`
	// Default: CURRENT_TIMESTAMP
	UpdateTime time.Time  `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type ObjectRelationMeta struct {
	ID             int    `db:"id"`
	RelationTypeID int    `db:"relation_type_id"`
	Name           string `db:"name"`
	// Comment: 1: STRING 2: INTEGER, 3: DOUBLE, 4: BOOLEAN
	// Default: 1
	ValueType   int    `db:"value_type"`
	Description string `db:"description"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type ObjectRelationMetaValue struct {
	FromObjectID   int    `db:"from_object_id"`
	RelationTypeID int    `db:"relation_type_id"`
	ToObjectID     int    `db:"to_object_id"`
	MetaID         int    `db:"meta_id"`
	Value          string `db:"value"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}

type DeletedObjectRelationMetaValue struct {
	FromObjectID   int    `db:"from_object_id"`
	RelationTypeID int    `db:"relation_type_id"`
	ToObjectID     int    `db:"to_object_id"`
	MetaID         int    `db:"meta_id"`
	Value          string `db:"value"`
	// Default: CURRENT_TIMESTAMP
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
	DeleteTime *time.Time `db:"delete_time"`
}
