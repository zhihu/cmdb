# Create object state

**URL** : `/v1/api/objects/:type_name/:object_name/logs`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name.
  * `object_name=[string]` The object name.

**Method** : `POST`

**Auth required** : YES

**Permissions required** : None

**Data constraints** : None

**Header constraints** : None

**Data examples**

Partial data is allowed.

```json
{
    "level": "EMERGENCY",
    "format": "TEXT",
    "source": "USER",
    "message": "log",
    "created_by": "administrator"
}
```

## Success Responses

**Condition** : Data provided is valid

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#objectLog",
    "type_name": "SERVER_NODE",
    "object_name": "xxx.xxx.xxx",
    "level": "EMERGENCY",
    "format": "TEXT",
    "source": "USER",
    "message": "log",
    "created_by": "administrator",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```

## Error Response

**Condition** : If provided data is invalid, e.g. invalid log level

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
    "level": "INVALID_LOG_LEVEL"
}
```

## Notes

* Endpoint will ignore irrelevant and read-only data such as parameters that
  don't exist, or fields that are not editable like `create_time` or `delete_time`.