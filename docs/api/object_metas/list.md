# List object metas

Get the details of all object metas

**URL** : `/v1/api/types/:type_name/metas`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name. Wildcard `-` can be used to list metas of all types.

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
    "kind": "cmdb#stateList",
    "next_page_token": "xxx",
    "metas": [
        {
            "kind": "cmdb#objectMeta",
            "type_name": "BUDGET",
            "meta_name": "TIME_PERIOD",
            "value_type": "STRING",
            "description": "Budget time period",
            "create_time": "2017-01-15T01:30:15.01Z"
        },
        {
            "kind": "cmdb#objectMeta",
            "type_name": "SERVER_NODE",
            "meta_name": "CPU_CORES",
            "value_type": "INTEGER",
            "description": "Number of cores per physical CPU",
            "create_time": "2017-01-15T01:30:15.01Z",
            "delete_time": "2017-01-15T01:30:15.01Z"
        }
    ]
}
```

## Notes

* Deleted object metas will not be included unless explicitly specified
