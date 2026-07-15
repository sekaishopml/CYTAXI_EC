package openapi

const BaseSpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "CYTAXI API",
    "description": "Conversation-First Mobility Platform",
    "version": "1.0.0",
    "contact": { "name": "CYTAXI Team" }
  },
  "servers": [{ "url": "/api/v1", "description": "API Gateway" }],
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      },
      "apiKey": {
        "type": "apiKey",
        "in": "header",
        "name": "X-API-Key"
      }
    },
    "schemas": {
      "Error": {
        "type": "object",
        "properties": {
          "error": { "type": "string" },
          "details": { "type": "string" }
        }
      }
    }
  },
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check",
        "operationId": "healthCheck",
        "responses": {
          "200": {
            "description": "Service healthy",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": { "type": "string" },
                    "service": { "type": "string" },
                    "version": { "type": "string" }
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "tags": [
    { "name": "customer", "description": "Customer operations" },
    { "name": "driver", "description": "Driver operations" },
    { "name": "trip", "description": "Trip lifecycle" },
    { "name": "pricing", "description": "Fare and pricing" },
    { "name": "payment", "description": "Payment operations" },
    { "name": "notification", "description": "Notifications" },
    { "name": "admin", "description": "Administration" },
    { "name": "analytics", "description": "Business intelligence" },
    { "name": "matching", "description": "Driver matching" }
  ]
}`

type SpecLoader struct{}

func NewSpecLoader() *SpecLoader { return &SpecLoader{} }

func (sl *SpecLoader) Load() string { return BaseSpec }
