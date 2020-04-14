# Create object status

**URL** : `/v1/api/types/:type_name/statuses`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name. 

**Method** : `POST`

**Auth required** : YES

**Permissions required** : None

**Data constraints** : None

**Header constraints** : None

**Data examples**

Partial data is allowed.

```json
{
    "type_name": "BUDGET",
    "status": "REVIEWED",
    "description": "Budget is reviewed"
}
```

## Success Responses

**Condition** : Data provided is valid

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#status",
    "type_name": "BUDGET",
    "status": "REVIEWED",
    "description": "Budget is reviewed",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```

## Error Response

**Condition** : If provided data is invalid, e.g. a name field is too long, conflicted with existing object status or simply missing.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
    "status": "EXISTING_OBJECT_STATUS"
}
```

## Notes

* Endpoint will ignore irrelevant and read-only data such as parameters that
  don't exist, or fields that are not editable like `name` or `create_time`.
