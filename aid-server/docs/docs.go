// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Leon Lin",
            "url": "github.com/leon123858",
            "email": "a0970785699@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/ask": {
            "post": {
                "description": "Ask",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Ask",
                "parameters": [
                    {
                        "description": "Ask Request",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.AskRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "aid string",
                        "schema": {
                            "$ref": "#/definitions/server.Response"
                        }
                    }
                }
            }
        },
        "/api/login": {
            "post": {
                "description": "Login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Login Request",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT Token",
                        "schema": {
                            "$ref": "#/definitions/server.Response"
                        }
                    }
                }
            }
        },
        "/api/logout": {
            "post": {
                "description": "Logout",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Logout",
                "parameters": [
                    {
                        "description": "Logout Request",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.LogoutRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "empty string",
                        "schema": {
                            "$ref": "#/definitions/server.Response"
                        }
                    }
                }
            }
        },
        "/api/register": {
            "post": {
                "description": "Register",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register",
                "parameters": [
                    {
                        "description": "Register Request",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT Token",
                        "schema": {
                            "$ref": "#/definitions/server.Response"
                        }
                    }
                }
            }
        },
        "/api/trigger": {
            "post": {
                "description": "Trigger",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Trigger",
                "parameters": [
                    {
                        "description": "Trigger Request",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.TriggerRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "aid string",
                        "schema": {
                            "$ref": "#/definitions/server.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "server.AskRequest": {
            "type": "object",
            "properties": {
                "browser": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                }
            }
        },
        "server.LoginRequest": {
            "type": "object",
            "properties": {
                "aid": {
                    "type": "string"
                },
                "browser": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "sign": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        },
        "server.LogoutRequest": {
            "type": "object",
            "properties": {
                "aid": {
                    "type": "string"
                },
                "browser": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                }
            }
        },
        "server.RegisterRequest": {
            "type": "object",
            "properties": {
                "aid": {
                    "type": "string"
                },
                "browser": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "publicKey": {
                    "type": "string"
                }
            }
        },
        "server.Response": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "result": {
                    "type": "boolean"
                }
            }
        },
        "server.TriggerRequest": {
            "type": "object",
            "properties": {
                "browser": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "AID API Server",
	Description:      "This is a AID server implementation for my paper.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
