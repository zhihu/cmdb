# Update object meta value

**URL** : `/v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation/:to_object_name/metas/:meta_name`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The from object type name.
  * `object_name=[string]` The from object name.
  * `to_type_name=[string]` The to object type name.
  * `to_object_name=[string]` The to object name.
  * `relation=[string]` The relation type name.
  * `meta_name=[string]` The meta name.

**Method** : `PUT`

**Auth required** : YES

**Permissions required** : None

**Data constraints** : None

**Header constraints** : None

**Data examples**

Partial data is allowed.

```json
{
    "version": 123457,
    "value": "XXX"
}
```

## Success Responses

**Condition** : Data provided is valid

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#relationMetaValue",
    "from_type_name": "SERVER_NODE",
    "from_object_name": "xxx.xxx.xxx",
    "relation": "OWNER",
    "to_type_name": "BUSINESS_LINE",
    "to_object_name": "xxx",
    "version": 123457,
    "meta_name": "PURCHASE_TIME",
    "value": "XXX",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```

## Error Response

**Condition** : If provided data is invalid, e.g. a name field is too long.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
    "meta_name": "OBJECT_META_NAME_THAT_IS_TOO_LONG"
}
```

**Condition** : If provided version does not match the current version.

**Code** : `409 CONFLICT`

**Content example** :

```json
{
    "version": -1
}
```

## Notes

* Endpoint will ignore irrelevant and read-only data such as parameters that
  don't exist, or fields that are not editable like `create_time`.