package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": [
        "http"
    ],
    "swagger": "2.0",
    "info": {
        "description": "User Service API for microservice architecture",
        "title": "User Service API",
        "contact": {},
        "version": "1.0"
    },
    "host": "localhost:8081",
    "basePath": "/",
    "paths": {
        "/users": {
            "get": {
                "description": "Get list of all users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "Success response with user list",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "string",
                                    "example": "success"
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Lấy danh sách users thành công"
                                },
                                "data": {
                                    "type": "array",
                                    "items": {
                                        "$ref": "#/definitions/User"
                                    }
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new user with the given data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User created successfully",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "string",
                                    "example": "success"
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Tạo user thành công"
                                },
                                "data": {
                                    "$ref": "#/definitions/User"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input data",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Get a user by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User found successfully",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "string",
                                    "example": "success"
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Lấy thông tin user thành công"
                                },
                                "data": {
                                    "$ref": "#/definitions/User"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid ID",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update user data by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update a user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated user data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User updated successfully",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "string",
                                    "example": "success"
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Cập nhật user thành công"
                                },
                                "data": {
                                    "$ref": "#/definitions/User"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input data",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete user by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Delete a user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User deleted successfully",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "string",
                                    "example": "success"
                                },
                                "message": {
                                    "type": "string",
                                    "example": "Xóa user thành công"
                                },
                                "data": {
                                    "type": "object",
                                    "nullable": true
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid ID",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Check User Service health status",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "Service is healthy",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {
                                    "type": "string",
                                    "example": "ok"
                                },
                                "service": {
                                    "type": "string",
                                    "example": "user-service"
                                },
                                "time": {
                                    "type": "string",
                                    "format": "date-time"
                                }
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Nguyen Van A"
                },
                "email": {
                    "type": "string",
                    "example": "a@example.com"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time"
                }
            }
        },
        "ErrorResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "error"
                },
                "message": {
                    "type": "string",
                    "example": "Something went wrong"
                },
                "error": {
                    "type": "string",
                    "example": "Detailed error message"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "User Service API",
	Description:      "User Service API for microservice architecture",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
