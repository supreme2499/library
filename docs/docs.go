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
        "/songs": {
            "get": {
                "description": "Метод возвращает список песен, с фильтрацией и пагинацией(количество песен на странице определено в servise.SearchSongs и является константным).\nНомер страницы является параметром запроса. Если не указать страницу, то вернётся первые 10 песен подходящие по фильтру,\nесли же не указать ни одного фильтра, то метод вернёт первые 10 песен.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Получить песни с фильтрацией по всем полям и пагинацией",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Название песни для фильтра",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название группы для фильтра",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "yyyy-mm-dd",
                        "description": "Год для фильтра",
                        "name": "year",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Номер страницы",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "В случае успеха сервер вернёт массив JSON",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.RequestSong"
                            }
                        }
                    },
                    "400": {
                        "description": "{Status: \"ERROR\", Error: msg}",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "Изменяет данные песни. В качестве индикатора использует название песни и группу",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Изменить данные песни",
                "parameters": [
                    {
                        "description": "Данные необходимые для изменения данных",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RequestSong"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ответ в случае успешного изменения {Status: \"OK\"}",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    },
                    "400": {
                        "description": "Ответ в случае возникновения ошибки {Status: \"ERROR\", Error: msg}",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "После получения запроса сервис делает API запрос, если API даёт ответ с дополнительной информацией о песне,\nто мы сохраняем данные в Postgres. Если же API возвращает ошибку, то мы ничего не сохраняем и возвращаем ошибку.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Добавить песню в библиотеку",
                "parameters": [
                    {
                        "description": "Данные необходимые для добавления.",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ответ сервера, при успешном сохранении {Status: \"OK\"}",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    },
                    "400": {
                        "description": "{Status: \"ERROR\", Error: msg}",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаляет песню из библиотеки по названию и группе",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Удалить песню из библиотеки",
                "parameters": [
                    {
                        "description": "Данные необходимые для удаления песни",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ответ в случае успешного удаления {Status: \"OK\"}",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    },
                    "400": {
                        "description": "Ответ в случае возникновения ошибки {Status: \"ERROR\", Error: msg}",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    }
                }
            }
        },
        "/songs/verse": {
            "get": {
                "description": "Метод вернёт текс песни c пагинацией по куплетам строкой.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Получить текст песни с пагинацией по куплетам",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Название песни для её поиска",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Название группы для её поиска",
                        "name": "group",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Количество строк которе необходимо вернуть (по умолчанию 1)",
                        "name": "verse",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "С какой строки возвращать (по умолчанию 1)",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "В случае успеха сервер венёт строку",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{Status: \"ERROR\", Error: msg}",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.Request": {
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "type": "string",
                    "example": "Кишлак"
                },
                "song": {
                    "type": "string",
                    "example": "Первый миллион"
                }
            }
        },
        "handlers.RequestSong": {
            "type": "object",
            "required": [
                "group",
                "link",
                "release_date",
                "song",
                "text"
            ],
            "properties": {
                "group": {
                    "type": "string",
                    "example": "Кишлак"
                },
                "link": {
                    "type": "string",
                    "example": "https://youtu.be/SdcNXIPP9UY?si=1QIYeuOSlNsDtZaS"
                },
                "release_date": {
                    "type": "string",
                    "example": "2024-11-29T10:00:00+03:00"
                },
                "song": {
                    "type": "string",
                    "example": "Первый миллион"
                },
                "text": {
                    "type": "string",
                    "example": "текст песни"
                }
            }
        },
        "handlers.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
