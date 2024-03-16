// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/assets": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves a list of all assets, including their associated user details.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Get all assets",
                "responses": {
                    "200": {
                        "description": "Success: Returns all assets.",
                        "schema": {
                            "$ref": "#/definitions/api.getAllAssetsResponse"
                        }
                    },
                    "500": {
                        "description": "Error: Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Creates a new asset based on the title, privacy setting, asset URL, and asset type provided in the request.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Create new asset",
                "parameters": [
                    {
                        "description": "Create Asset Request",
                        "name": "CreateAssetRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CreateAssetRequest"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Asset creation successful, returns created asset details along with a success message.",
                        "schema": {
                            "$ref": "#/definitions/api.CreateAssetsResponse"
                        }
                    }
                }
            }
        },
        "/assets/gaussian/{id}": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Updates the URL for a specific gaussian asset based on the provided ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Update point cloud URL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Asset ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update Gaussian URL Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UpdateGaussianUrlRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "URL updated successfully",
                        "schema": {
                            "$ref": "#/definitions/api.UpdateGaussianUrlResponse"
                        }
                    }
                }
            }
        },
        "/assets/like/{id}": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Marks an asset as liked by the current user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Like an asset",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Asset ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Asset liked successfully",
                        "schema": {
                            "$ref": "#/definitions/api.LikeAssetResponse"
                        }
                    }
                }
            }
        },
        "/assets/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves a list of my assets",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Get my assets",
                "responses": {
                    "200": {
                        "description": "Success: Returns all assets.",
                        "schema": {
                            "$ref": "#/definitions/api.getMyAssetsResponse"
                        }
                    },
                    "500": {
                        "description": "Error: Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/assets/pointcloud/{id}": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Updates the URL for a specific point cloud asset based on the provided ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Update point cloud URL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Asset ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update Point Cloud URL Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UpdatePointCloudUrlRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "URL updated successfully",
                        "schema": {
                            "$ref": "#/definitions/api.UpdatePointCloudUrlResponse"
                        }
                    }
                }
            }
        },
        "/assets/unlike/{id}": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Marks an asset as unliked by the current user, removing the like.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Unlike an asset",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Asset ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Asset unliked successfully",
                        "schema": {
                            "$ref": "#/definitions/api.UnlikeAssetResponse"
                        }
                    }
                }
            }
        },
        "/assets/{id}": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Remove my asset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Remove my asset",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Asset ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Asset removed successfully",
                        "schema": {
                            "$ref": "#/definitions/api.removeAssetResponse"
                        }
                    }
                }
            }
        },
        "/auth/google": {
            "post": {
                "description": "Authenticate user with Google OAuth token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Google Auth",
                "parameters": [
                    {
                        "description": "Google OAuth token",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.googleRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User info from Google",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/signin": {
            "post": {
                "description": "Login user with the provided credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "User login details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.loginUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User login successful",
                        "schema": {
                            "$ref": "#/definitions/api.loginUserResponse"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Register a new user with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User registration details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.registerUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User registration successful",
                        "schema": {
                            "$ref": "#/definitions/api.registerUserResponse"
                        }
                    }
                }
            }
        },
        "/tags/search": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves a list of tags that match the given search keyword",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tags"
                ],
                "summary": "Get tags by search keyword",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search keyword",
                        "name": "keyword",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/api.GetTagBySearchKeywordResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieve user information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user data",
                "responses": {
                    "200": {
                        "description": "User information retrieved successfully",
                        "schema": {
                            "$ref": "#/definitions/api.UserResponse"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update user information based on the provided user ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update user information",
                "parameters": [
                    {
                        "description": "User update details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.updateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User information updated successfully",
                        "schema": {
                            "$ref": "#/definitions/api.updateUserResponse"
                        }
                    }
                }
            }
        },
        "/users/password": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Change user password based on the provided user ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Change user password",
                "parameters": [
                    {
                        "description": "User password update details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.changeUserPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User password updated successfully",
                        "schema": {
                            "$ref": "#/definitions/api.changeUserPasswordResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.AssetResponse": {
            "type": "object",
            "properties": {
                "assetType": {
                    "type": "string"
                },
                "assetUrl": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "gaussianUrl": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "isLikedByMe": {
                    "type": "boolean"
                },
                "isPrivate": {
                    "type": "boolean"
                },
                "likes": {
                    "type": "integer"
                },
                "pointCloudUrl": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "thumbnailUrl": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/api.UserResponse"
                }
            }
        },
        "api.CreateAssetRequest": {
            "type": "object",
            "required": [
                "assetType",
                "assetUrl",
                "isPrivate",
                "title"
            ],
            "properties": {
                "assetType": {
                    "type": "string",
                    "enum": [
                        "images",
                        "video"
                    ]
                },
                "assetUrl": {
                    "type": "string"
                },
                "isPrivate": {
                    "type": "boolean"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "api.CreateAssetsResponse": {
            "type": "object",
            "properties": {
                "asset": {
                    "$ref": "#/definitions/api.AssetResponse"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "api.GetTagBySearchKeywordResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/db.Tags"
                    }
                }
            }
        },
        "api.LikeAssetResponse": {
            "type": "object",
            "properties": {
                "asset": {
                    "$ref": "#/definitions/api.AssetResponse"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.UnlikeAssetResponse": {
            "type": "object",
            "properties": {
                "asset": {
                    "$ref": "#/definitions/api.AssetResponse"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.UpdateGaussianUrlRequest": {
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "api.UpdateGaussianUrlResponse": {
            "type": "object",
            "properties": {
                "asset": {
                    "$ref": "#/definitions/api.AssetResponse"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.UpdatePointCloudUrlRequest": {
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "api.UpdatePointCloudUrlResponse": {
            "type": "object",
            "properties": {
                "asset": {
                    "$ref": "#/definitions/api.AssetResponse"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.UserResponse": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "passwordChangedAt": {
                    "type": "string"
                },
                "provider": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "api.changeUserPasswordRequest": {
            "type": "object",
            "properties": {
                "newPassword": {
                    "type": "string"
                },
                "oldPassword": {
                    "type": "string"
                }
            }
        },
        "api.changeUserPasswordResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/api.UserResponse"
                }
            }
        },
        "api.getAllAssetsResponse": {
            "type": "object",
            "properties": {
                "assets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.AssetResponse"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.getMyAssetsResponse": {
            "type": "object",
            "properties": {
                "assets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.AssetResponse"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.googleRequest": {
            "type": "object",
            "required": [
                "token"
            ],
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "api.loginUserRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                }
            }
        },
        "api.loginUserResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/api.UserResponse"
                }
            }
        },
        "api.registerUserRequest": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                }
            }
        },
        "api.registerUserResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/api.UserResponse"
                }
            }
        },
        "api.removeAssetResponse": {
            "type": "object",
            "properties": {
                "asset": {
                    "$ref": "#/definitions/api.AssetResponse"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.updateUserRequest": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "api.updateUserResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/api.UserResponse"
                }
            }
        },
        "db.Tags": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Segment3d App API Documentation",
	Description:      "This is a documentation for Segment3d App API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
