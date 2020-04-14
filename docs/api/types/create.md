# Create object type

**URL** : `/v1/api/types`

**Method** : `POST`

**Auth required** : YES

**Permissions required** : None

**Data constraints** : None

**Header constraints** : None

**Data examples**

Partial data is allowed.

```json
{
    "name": "DATA_CENTER",
    "description": "Data Center"
}
```

## Success Responses

**Condition** : Data provided is valid

**Code** : `200 OK`

**Content example**

```json
{
    "kind": "cmdb#type",
    "name": "DATA_CENTER",
    "description": "Data Center",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```

## Error Response

**Condition** : If provided data is invalid, e.g. a name field is too long, conflicted with existing object type or simply missing.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
    "name": "EXISTING_OBJECT_TYPE"
}
```

## Notes

* Endpoint will ignore irrelevant and read-only data such as parameters that
  don't exist, or fields that are not editable like `delete_time` or `create_time`.
