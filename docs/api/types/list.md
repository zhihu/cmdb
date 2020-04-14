# List object types

List the details of all types

**URL** : `/v1/api/types`

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
  "kind": "cmdb#typeList",
  "next_page_token": "xxx",
  "types": [
      {
          "kind": "cmdb#type",
          "name": "ANY",
          "description": "Special object type that represents anything",
          "create_time": "2017-01-15T01:30:15.01Z"
      },
      {
          "kind": "cmdb#type",
          "name": "SERVER_NODE",
          "description": "Server Node",
          "create_time": "2017-01-15T01:30:15.01Z",
          "delete_time": "2017-01-15T01:30:15.01Z"
      }
  ]
}
```

## Notes

* Deleted object types will not be included unless explicitly specified
