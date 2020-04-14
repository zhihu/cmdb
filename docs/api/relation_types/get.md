# Get related objects of a given object

Get the details of related objects for a given object

**URL** : `/v1/api/relation_types/:from_type_name/:to_type_name/:relation`

**URL Parameters**

* **Required:**
  * `from_type_name=[string]` The from object type name.
  * `to_type_name=[string]` The to object type name.
  * `relation=[string]` The relation type name.

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#relationType",
    "from_type_name": "FROM_OBJECT_NAME",
    "to_type_name": "TO_OBJECT_NAME",
    "relation": "BUDGETS",
    "description": "Budget of business line",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```
