# List object relation types

Get the details of object relation types

**URL** : `GET /v1/api/relation_types/:from_type_name/:to_type_name`

**URL Parameters**

* **Required:**
  * `from_type_name=[string]` The from object type name. Wildcard `-` can be used to list relation types of all from object types.
  * `to_type_name=[string]` The to object type name. Wildcard `-` can be used to list relation types of all to object types.

**Query String Parameters**

* **Optional:**
  * `show_deleted=[bool]` if deleted object types should be included in the response
  * `page_token=[string]` The pagination token
  * `page_size=[int32]` The pagination size

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#relationTypeList",
    "next_page_token": "xxx",
    "relation_types": [
        {
            "kind": "cmdb#relationType",
            "from_type_name": "BUSINESS_LINE",
            "to_type_name": "BUDGET",
            "relation": "BUDGETS",
            "description": "Budget of business line",
            "create_time": "2017-01-15T01:30:15.01Z"
        },
        {
            "kind": "cmdb#relationType",
            "from_type_name": "BUDGET",
            "to_type_name": "BUDGET_ITEM",
            "relation": "CONTAINS",
            "description": "Contains",
            "create_time": "2017-01-15T01:30:15.01Z",
            "delete_time": "2017-01-15T01:30:15.01Z"
        }
    ]
```

## Notes

* Deleted object relation types will not be included unless explicitly specified
