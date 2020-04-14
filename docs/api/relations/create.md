# Create objects relation

**URL** : `/v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation/:to_object_name`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The from object type name.
  * `object_name=[string]` The from object name.
  * `to_type_name=[string]` The to object type name.
  * `relation=[string]` The relation type name.
  * `to_object_name=[string]` The to object name.

**Method** : `POST`

**Auth required** : YES

**Permissions required** : None

**Data constraints** : None

**Header constraints** : None

**Data examples** : `{}`

## Success Responses

**Condition** : Data provided is valid

**Code** : `204 NO CONTENT`

**Content examples**

```json
{
    "kind": "cmdb#relation",
    "from_type_name": "SERVER_NODE",
    "from_object_name": "xxx.xxx",
    "to_type_name": "BUSINESS_LINE",
    "to_object_name": "xxx",
    "relation": "BELONGS",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```

## Error Response

**Condition** : If provided data is invalid, e.g. an invalid from type name.

**Code** : `400 BAD REQUEST`

## Notes

* Endpoint will ignore irrelevant and read-only data such as parameters that
  don't exist, or fields that are not editable like `create_time`.