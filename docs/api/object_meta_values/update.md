# Update object meta value

**URL** : `/v1/api/objects/:type_name/:object_name/metas/:meta_name`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name.
  * `object_name=[string]` The object name.
  * `meta_name=[string]` The meta name

**Method** : `PUT`

**Auth required** : YES

**Permissions required** : None

**Data constraints** : None

**Header constraints** : None

**Data examples**

Partial data is allowed.

```json
{
    "version": 123456,
    "value": "XXX"
}
```

## Success Responses

**Condition** : Data provided is valid

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#objectMetaValue",
    "type_name": "SERVER_NODE",
    "object_name": "xxx.xxx.xxx",
    "meta_name": "CPU_DESCRIPTION",
    "version": 123457,
    "value": "XXX"
}
```

## Error Response

**Condition** : If provided data is invalid, e.g. a name field is too long, conflicted with existing object meta or simply missing.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
    "name": "EXISTING_OBJECT_META"
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
