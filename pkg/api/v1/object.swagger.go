package v1
var Swagger=`
{
  "swagger": "2.0",
  "info": {
    "title": "pkg/api/v1/object.proto",
    "version": "version not set"
  },
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/v1/api/objects/{type}": {
      "get": {
        "operationId": "Objects_List",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/v1ObjectListResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "parameters": [
          {
            "name": "type",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "view",
            "in": "query",
            "required": false,
            "type": "string",
            "enum": [
              "BASIC",
              "NORMAL",
              "RICH"
            ],
            "default": "BASIC"
          },
          {
            "name": "query",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "page_token",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "page_size",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int32"
          },
          {
            "name": "show_deleted",
            "in": "query",
            "required": false,
            "type": "boolean",
            "format": "boolean"
          }
        ],
        "tags": [
          "Objects"
        ]
      }
    }
  },
  "definitions": {
    "protobufAny": {
      "type": "object",
      "properties": {
        "type_url": {
          "type": "string"
        },
        "value": {
          "type": "string",
          "format": "byte"
        }
      }
    },
    "runtimeError": {
      "type": "object",
      "properties": {
        "error": {
          "type": "string"
        },
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        },
        "details": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/protobufAny"
          }
        }
      }
    },
    "v1Object": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "status": {
          "type": "string"
        },
        "state": {
          "type": "string"
        },
        "create_time": {
          "type": "string",
          "format": "date-time"
        },
        "metas": {
          "type": "object",
          "additionalProperties": {
            "$ref": "#/definitions/v1ObjectMeta"
          }
        }
      }
    },
    "v1ObjectListResponse": {
      "type": "object",
      "properties": {
        "kind": {
          "type": "string"
        },
        "objects": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1Object"
          }
        },
        "next_page_token": {
          "type": "string"
        }
      }
    },
    "v1ObjectMeta": {
      "type": "object",
      "properties": {
        "value_type": {
          "$ref": "#/definitions/v1ValueType"
        },
        "value": {
          "type": "string"
        }
      }
    },
    "v1ObjectView": {
      "type": "string",
      "enum": [
        "BASIC",
        "NORMAL",
        "RICH"
      ],
      "default": "BASIC"
    },
    "v1ValueType": {
      "type": "string",
      "enum": [
        "STRING",
        "INTEGER",
        "DOUBLE",
        "BOOLEAN"
      ],
      "default": "STRING"
    }
  }
}
`
