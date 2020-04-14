# Create object relation meta

**URL** : `/v1/api/relation_types/:from_type_name/:to_type_name/:relation/metas`

**URL Parameters**

* **Required:**
  * `from_type_name=[string]` The from object type name.
  * `to_type_name=[string]` The to object type name.
  * `relation=[string]` The relation type name.

**Method** : `POST`

**Auth required** : YES

**Permissions required** : None

**Data constraints** : None

**Header constraints** : None

**Data examples**

Partial data is allowed.

```json
{
   "from_type_name": "SERVER_NODE",
   "to_type_name": "BUSINESS_LINE",
   "relation": "BELONGS",
   "meta_name": "PURCHASE_TIME",
   "value_type": "String",
   "description": "When a server was purchased by a business line",
}
```

## Success Responses

**Condition** : Data provided is valid

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#relationMeta",
    "from_type_name": "SERVER_NODE",
    "to_type_name": "BUSINESS_LINE",
    "relation": "BELONGS",
    "meta_name": "PURCHASE_TIME",
    "value_type": "String",
    "description": "When a server was purchased by a business line",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```

## Error Response

**Condition** : If provided data is invalid, e.g. a name field is too long, conflicted with existing object relation meta or simply missing.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
    "meta_name": "EXISTING_OBJECT_RELATION_META"
}
```

## Notes

* Endpoint will ignore irrelevant and read-only data such as parameters that
  don't exist, or fields that are not editable like `create_time`.
