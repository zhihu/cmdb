# Create object meta

**URL** : `/v1/api/types/:type_name/metas`

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
    "meta_name": "CPU_DESCRIPTION",
    "value_type": "STRING",
    "description": "CPU description, vendor labels"
}
```

## Success Responses

**Condition** : Data provided is valid

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#objectMeta",
    "type_name": "SERVER_NODE",
    "meta_name": "CPU_DESCRIPTION",
    "value_type": "STRING",
    "description": "CPU description, vendor labels",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```

## Error Response

**Condition** : If provided data is invalid, e.g. a name field is too long, conflicted with existing object meta or simply missing.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
    "meta_name": "EXISTING_OBJECT_META"
}
```

## Notes

* Endpoint will ignore irrelevant and read-only data such as parameters that
  don't exist, or fields that are not editable like `create_time`.
