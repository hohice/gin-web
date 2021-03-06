// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2018-09-18 14:27:49.043027077 +0800 CST m=+0.079683128

package docs

import (
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "Gin Web API server starter.",
        "title": "ginS",
        "contact": {
            "name": "hohice",
            "url": "https://github.com/hohice",
            "email": "hohice@163.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0.0"
    },
    "basePath": "/api/v1",
    "paths": {
        "/application/build/name/{name}/version/{version}": {
            "post": {
                "description": "Modify Application Config",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "application"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of the config",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "version of the config",
                        "name": "version",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid Name supplied!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "404": {
                        "description": "Instance not found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "405": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    }
                }
            }
        },
        "/config": {
            "put": {
                "description": "Modify Application Config",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "parameters": [
                    {
                        "description": "data of the config",
                        "name": "config",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid Name supplied!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "404": {
                        "description": "Instance not found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "405": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "Modify Application Config",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "parameters": [
                    {
                        "description": "data of the config",
                        "name": "config",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid Name supplied!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "404": {
                        "description": "Instance not found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "405": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    }
                }
            }
        },
        "/config/name/{name}/version/{version}": {
            "get": {
                "description": "Get Application Config",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of the config",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "version of the config",
                        "name": "version",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/config.ConfigType"
                        }
                    },
                    "400": {
                        "description": "Invalid Name supplied!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "404": {
                        "description": "Instance not found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "405": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete Application Config",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of the config",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "version of the config",
                        "name": "version",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid Name supplied!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "404": {
                        "description": "Instance not found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "405": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/ex.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "config.ConfigType": {
            "type": "object",
            "properties": {
                "context": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "ex.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

type s struct{}

func (s *s) ReadDoc() string {
	return doc
}
func init() {
	swag.Register(swag.Name, &s{})
}
