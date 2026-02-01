# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**数字文旅一体化小程序平台** (Digital Cultural Tourism Mini-Program Platform) - A Go-based backend API service built with Gin and Tencent CloudBase (TCB) HTTP API, providing RESTful APIs for cultural tourism content discovery, user-generated content (UGC), and location-based services (LBS).

**Core Value Chain**: Content Discovery (LBS) → User Participation (UGC) → Offline Photo Pickup (Travel Photo Booth) → Interactive Reviews

## Build & Development Commands

### Backend Development

```bash
# Navigate to backend directory
cd cultural-tourism-backend

# Install dependencies
go mod download

# Run development server with hot-reload (requires Air)
air

# Run without hot-reload
go run main.go

# Build for production
go build -o main .

# Run tests
go test ./...
```

### Environment Setup

Create a `.env` file in `cultural-tourism-backend/`:

```env
CLOUDBASE_ENV_ID=your-env-id
CLOUDBASE_ACCESS_TOKEN=your-access-token
```

### API Documentation

- **Swagger UI**: Available at `http://localhost:8080/swagger/index.html` when server is running
- Regenerate Swagger docs: `swag init` (from backend directory)

## Architecture Overview

### Three-Layer Architecture

1. **Routes Layer** (`routes/router.go`)
   - RESTful route definitions
   - Groups API endpoints by resource type

2. **Controllers Layer** (`controllers/`)
   - Request validation and binding
   - HTTP response formatting
   - Delegates business logic to services

3. **Services Layer** (`services/`)
   - Core business logic
   - Data validation and transformation
   - Directly calls TCB client for database operations

### TCB Client Abstraction

**Critical**: All database operations MUST use the generic `tcb/client.go` HTTP API wrapper.

- **Location**: `tcb/client.go`
- **Purpose**: Unified CloudBase HTTP API client
- **Key Constraint**: Bottom layer MUST remain generic - no hardcoded business logic
- **Filter Pattern**: Controllers/Services construct complete TCB query objects with `where`, `orderBy`, etc.

**Example Filter Construction**:
```go
filter := map[string]interface{}{
    "where": map[string]interface{}{
        "region_id": map[string]interface{}{"$eq": regionID},
        "status": map[string]interface{}{"$eq": 1},
    },
    "orderBy": []interface{}{
        map[string]interface{}{"field": "created_at", "order": "desc"},
    },
}
```

### Data Models

Models are defined in two places:
- **Go Structs**: `models/` directory (for type safety and Swagger annotations)
- **JSON Schemas**: `model-json/` directory (for CloudBase data model definitions)

**Collections**: `regions`, `pois`, `theme`, `photo`, `comment`, `product`

## Core Development Rules

### 1. Database Operations Protocol

**RED LINES (Do NOT violate)**:
- ❌ NEVER use MongoDB native driver
- ❌ NEVER hardcode operators (e.g., `$eq`) in `tcb/client.go`
- ❌ NEVER bypass the TCB client abstraction
- ✅ ALWAYS construct filters at Controller/Service layer
- ✅ ALWAYS use `PUT` + `/update` + proper filter structure for updates
- ✅ ALWAYS use `POST` + `/delete` + proper filter for deletions

### 2. RESTful API Standards

**Operations**:
- **Create**: `POST /api/{resource}`
- **List**: `GET /api/{resource}` (with query params)
- **Detail**: `GET /api/{resource}/:id`
- **Update**: `PUT /api/{resource}/:id`
- **Delete**: `DELETE /api/{resource}/:id`

**Pagination**: Default support with `page` and `size` query parameters.

### 3. Security & Data Integrity

- **System Fields**: Auto-managed fields like `_id`, `_openid`, `created_at`, `updated_at` must be stripped from user input
- **Status Defaults**: UGC content (photos, comments) defaults to `status=0` (pending review)
- **Validation**: Always validate enum values (e.g., POI type: `scenic`, `food`, `hotel`, `booth`)

### 4. Business Boundaries

**NOT in Scope** (Do NOT implement):
- Payment processing (only mini-program jump redirects)
- Online image editing / AI composition
- Cross-platform jumps (WeChat mini-program only)

## Project Structure

```
cultural-tourism-backend/
├── config/           # Configuration management
├── controllers/      # HTTP request handlers
├── database/         # Database initialization
├── docs/            # Swagger generated files
├── model-json/      # CloudBase JSON data models
├── models/          # Go struct definitions
├── routes/          # Route registration
├── services/        # Business logic layer
├── tcb/             # CloudBase HTTP client
├── tests/           # Test files
├── .env             # Environment variables (not in git)
├── go.mod           # Go dependencies
├── main.go          # Application entry point
└── PROJECT_CONTEXT.md  # Detailed project context (Chinese)
```

## Key Implementation Patterns

### LBS Distance Calculation

POI distance queries use MongoDB geospatial operators:

```go
where := map[string]interface{}{
    "location": map[string]interface{}{
        "$near": map[string]interface{}{
            "$geometry": map[string]interface{}{
                "type": "Point",
                "coordinates": []float64{longitude, latitude},
            },
            "$maxDistance": radiusKM * 1000,
        },
    },
}
```

### Service Layer Pattern

Services handle:
- Field initialization (timestamps, status defaults)
- Security enforcement (strip system fields)
- Query filter construction
- Calling TCB client methods

**Example**:
```go
func CreatePOI(poi *models.POI) (map[string]interface{}, error) {
    poi.ID = ""  // Security: prevent ID injection
    poi.Status = 1
    poi.CreatedAt = time.Now().Format(time.RFC3339)
    poi.UpdatedAt = time.Now().Format(time.RFC3339)

    return tcb.Client.CreateData(CollectionPOI, poi)
}
```

## Important References

- **PROJECT_CONTEXT.md**: Chinese documentation with detailed roadmap, data schemas, and development protocols
- **PRD Document**: `数字文旅一体化小程序平台___产品需求规格说明书__PRD_ (1).pdf`
- **Code Review Guide**: `code_reviwe.md` - High-standard audit checklist for 10k+ users
- **CloudBase AI Rules**: `.codebuddy/rules/tcb/CODEBUDDY.md` - CloudBase development best practices

## Current Implementation Status

**Completed (✅)**:
- Phase 1: Infrastructure (TCB client, Swagger, error handling)
- Phase 2: Region management
- Phase 3: POI management with LBS
- Phase 4: UGC travel photo community (themes + photos)
- Phase 5: Interactive features (comments + product navigation)

**Pending**:
- Phase 6: Travel photo booth integration, album, favorites
- Phase 7: Admin management and content moderation APIs

## Development Workflow

1. Define/update data model in `model-json/` and `models/`
2. Create service functions in `services/`
3. Implement controller handlers in `controllers/`
4. Register routes in `routes/router.go`
5. Add Swagger annotations to controller methods
6. Run `swag init` to regenerate documentation
7. Test endpoints via Swagger UI or API client

## Notes for AI Assistants

- Always read `PROJECT_CONTEXT.md` for detailed Chinese specifications
- Follow the three-layer architecture pattern strictly
- Maintain the generic nature of `tcb/client.go`
- Construct query filters at service/controller level
- Use proper error handling and validation
- All Swagger annotations should be in English
