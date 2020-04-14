# Get object relation meta

Get the details of a given object relation meta

**URL** : `/v1/api/relation_types/:from_type_name/:to_type_name/:relation/metas/:meta_name`

**URL Parameters**

* **Required:**
  * `from_type_name=[string]` The from object type name.
  * `to_type_name=[string]` The to object type name.
  * `relation=[string]` The to object type name.
  * `meta_name=[string]` The relation meta name.

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#relationMeta",
    "from_type_name": "SERVER_NODE",
    "to_type_name": "BUSINESS_LINE",
    "relation": "BELONGS",
    "meta_name": "PURCHASE_TIME",
    "value_type": "String",
    "description": "When a server was purchased by a business line",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```
