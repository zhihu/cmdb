# Create object relation type

**URL** : `/v1/api/relation_types/:from_type_name/:to_type_name`

**URL Parameters**

* **Required:**
  * `from_type_name=[string]` The from object type name.
  * `to_type_name=[string]` The to object type name.

**Method** : `POST`

**Auth required** : YES

**Permissions required** : None

**Data constraints** : None

**Header constraints** : None

**Data examples**

Partial data is allowed.

```json
{
    "relation": "BUDGETS",
    "description": "Budget of business line",
}
```

## Success Responses

**Condition** : Data provided is valid

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#relationType",
    "from_type_name": "BUSINESS_LINE",
    "to_type_name": "BUDGET",
    "relation": "BUDGETS",
    "description": "Budget of business line",
}
```

## Error Response

**Condition** : If provided data is invalid, e.g. a name field is too long, conflicted with existing object relation type or simply missing.

**Code** : `400 BAD REQUEST`

**Content example** :

```json
{
    "relation": "EXISTING_OBJECT_RELATION_TYPE"
}
```

## Notes

* Endpoint will ignore irrelevant and read-only data such as parameters that
  don't exist, or fields that are not editable like `name` or `create_time`.