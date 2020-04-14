# Create object

**URL** : `/v1/api/objects/:type_name`

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
    "object_name": "abc.def.com",
    "description": "server",
    "status": "NEW",
    "state": "STARTING"
}
```

## Success Responses

**Condition** : Data provided is valid

**Code** : `200 OK`

**Content examples**

{
    "kind": "cmdb#object",
    "type_name": "SERVER_NODE",
    "object_name": "abc.def.com",
    "description": "server",
    "status": "NEW",
    "state": "STARTING", 
    "create_time": "2017-01-15T01:30:15.01Z"
},

## Error Response

**Condition** : If provided data is invalid, e.g. a name field is too long, conflicted with existing object state or simply missing.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
    "name": "EXISTING_OBJECT"
}
```

## Notes

* Endpoint will ignore irrelevant and read-only data such as parameters that
  don't exist, or fields that are not editable like `name` or `create_time`.
