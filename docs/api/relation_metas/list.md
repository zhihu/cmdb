# List object metas

Get the details of relation metas

**URL** : `/v1/api/relation_types/:from_type_name/:to_type_name/:relation/metas`

**URL Parameters**

* **Required:**
  * `from_type_name=[string]` The from object type name. Wildcard `-` can be used to list metas of all from object types.
  * `to_type_name=[string]` The to object type name. Wildcard `-` can be used to list metas of all to object types.
  * `relation=[string]` The to object type name. Wildcard `-` can be used to list metas of all relation types.

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
    "kind": "cmdb#relationMetaList",
    "next_page_token": "xxx",
    "relation_metas": [
        {
            "kind": "cmdb#relationMeta",
            "from_type_name": "SERVER_NODE",
            "to_type_name": "BUSINESS_LINE",
            "relation": "BELONGS",
            "meta_name": "PURCHASE_TIME",
            "value_type": "String",
            "description": "When a server was purchased by a business line",
            "create_time": "2017-01-15T01:30:15.01Z"
        },
        {
            "kind": "cmdb#relationMeta",
            "from_type_name": "SERVER_NODE",
            "to_type_name": "BUSINESS_LINE",
            "relation": "BELONGS",
            "meta_name": "TERMINATE_TIME",
            "value_type": "Boolean",
            "description": "When a server was terminated by a business line",
            "create_time": "2017-01-15T01:30:15.01Z",
            "delete_time": "2017-01-15T01:30:15.01Z"
        }
    ]
}
```

## Notes

* Deleted object metas will not be included unless explicitly specified
