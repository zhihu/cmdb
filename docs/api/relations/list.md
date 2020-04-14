# List object relations

Get the details of all object relations

**URL** : `/v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The from object type name. Wildcard `-` can be used to list metas of all from object types.
  * `object_name=[string]` The from object name. Wildcard `-` can be used to list metas of all from objects.
  * `to_type_name=[string]` The to object type name. Wildcard `-` can be used to list metas of all to object types.
  * `relation=[string]` The relation type name. Wildcard `-` can be used to list metas of all relation types.

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
    "kind": "cmdb#relationList",
    "next_page_token": "xxx",
    "relations": [
        {
            "kind": "cmdb#relation",
            "from_type_name": "SERVER_NODE",
            "from_object_name": "xxx.xxx",
            "to_type_name": "BUSINESS_LINE",
            "to_object_name": "xxx",
            "relation": "BELONGS",
            "create_time": "2017-01-15T01:30:15.01Z"
        },
        {
            "kind": "cmdb#relation",
            "from_type_name": "SERVER_NODE",
            "from_object_name": "yyy.yyy",
            "to_type_name": "BUSINESS_LINE",
            "to_object_name": "yyy",
            "relation": "BELONGS",
            "create_time": "2017-01-15T01:30:15.01Z",
            "delete_time": "2017-01-15T01:30:15.01Z"
        }
    ]
}
```

## Notes

* Deleted object states will not be included unless explicitly specified