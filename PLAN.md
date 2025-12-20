# Development Plan

This document outlines the development roadmap for Personal Excalidraw.

## Phase 1: Frontend Infrastructure ✅

**Status**: Complete

- [x] Project initialization
- [x] SvelteKit setup with TypeScript
- [x] Tailwind CSS + DaisyUI integration
- [x] ExcalidrawWrapper component with React integration
- [x] Svelte stores for state management
- [x] Home page with drawings list
- [x] Drawing editor with dynamic routing
- [x] Mock data with 8 sample drawings

**Outcome**: Functional frontend with mock data and basic navigation.

## Phase 2: Local Storage & Real Data ✅

**Status**: Complete

### Completed

- [x] Connect drawing editor to load/save specific drawings
- [x] LocalStorage persistence with per-drawing storage (`excalidraw-drawing-{id}`)
- [x] Auto-save with 1-second debounce for performance
- [x] Implement actual delete functionality (removes both metadata and content)
- [x] ID type abstraction for consistent type safety
- [x] Drawing metadata timestamp synchronization
- [x] Data validation and error handling (corrupted data, quota exceeded)
- [x] Add drawing name editing with inline editing UI
- [x] Replace all mock data with real localStorage-backed data
- [x] Implement full CRUD operations for drawings

**Outcome**: Fully functional local-first application with complete drawing management capabilities.

## Phase 3: Backend Integration ✅

**Status**: Complete

### Backend Development ✅

- [x] Set up Go project structure with Clean Architecture
- [x] Implement PostgreSQL database schema with migrations
- [x] Create RESTful API endpoints:
  - [x] `GET /api/drawings` - List all drawings
  - [x] `GET /api/drawings/:id` - Get specific drawing
  - [x] `POST /api/drawings` - Create new drawing
  - [x] `PUT /api/drawings/:id` - Update drawing (with partial updates support)
  - [x] `DELETE /api/drawings/:id` - Delete drawing
- [x] Add slug-based URL support for SEO-friendly URLs
- [x] Database migrations system with reversible migrations
- [x] Comprehensive error handling and validation
- [x] Request/Response DTOs with proper validation
- [x] Comprehensive test coverage (unit + integration tests)

### Frontend Integration ✅

- [x] Create API client service with TypeScript types
- [x] Replace localStorage calls with API calls
- [x] Add TanStack Query for data fetching and caching
- [x] Implement loading states and error handling
- [x] Real-time data synchronization with optimistic updates
- [x] Inline name editing with API integration

**Outcome**: Full-stack application with backend API, database persistence, and seamless frontend-backend integration. As for personel purpose, this will not include authentication feature.

### Goals

- Transition from local-only to cloud-backed storage
- Enable cross-device access

## Phase 4: Production Ready ✅

**Status**: Complete

Focuses on making the application stable, production-ready, and ready for daily use.

### Authentication & Security ✅

- [x] Backend authentication middleware with access key validation
- [x] Frontend auth store with localStorage persistence
- [x] AuthModal component for key input
- [x] API endpoint protection with Bearer token authentication
- [x] Auth validation endpoint (`GET /api/auth/validate`)
- [x] Configuration support for enabling/disabling auth
- [x] Constant-time comparison for access key validation (security)
- [x] Cross-tab authentication synchronization

### Code Quality & Refactoring ✅

- [x] Remove duplicate `respondJSON` function
- [x] Create shared HTTP utilities package
- [x] Consolidate response handling across handlers

### Production Deployment ✅

- [x] Auto-migration on backend startup
- [x] Embedded migration files with Go embed
- [x] Production-ready Dockerfile with multi-stage build
- [x] Docker Compose production configuration
- [x] Health check endpoints and Docker healthchecks
- [x] Non-root user for security
- [x] Optimized build flags (-ldflags="-w -s")
- [x] Production environment template (.env.production.example)

### Documentation ✅

- [x] Comprehensive deployment guide (DEPLOYMENT.md)
- [x] Production setup instructions
- [x] SSL/HTTPS configuration with Nginx
- [x] Backup and maintenance procedures
- [x] Troubleshooting guide
- [x] Security best practices
- [x] Updated README with deployment links

**Outcome**: Production-ready application with complete deployment documentation, automated migrations, and containerized deployment. Ready to deploy to any VPS.

## Phase 5: Public Share URL

**Status**: In Progress (Backend Complete ✅)

Enable sharing drawings with public URLs that allow view-only access with temporary editing capabilities.

### Backend Development ✅

- [x] Database schema updates:
  - [x] Add `share_token` column to drawings table (nullable, unique, indexed)
  - [x] Migration for new column with unique index (003_add_share_token_to_drawings)
  - [x] Drawing is public if `share_token` is not null
- [x] New API endpoints:
  - [x] `POST /api/drawings/:id/share` - Generate/regenerate public share token (UUID-based)
  - [x] `DELETE /api/drawings/:id/share` - Revoke public sharing (set share_token to null)
  - [x] `GET /api/public/:share_token` - Get drawing by share token (no auth required)
- [x] Public endpoint security:
  - [x] No authentication required for public endpoints
  - [x] Return only drawing content and metadata (no sensitive data)
  - [x] Authentication middleware updated to skip `/public/*` paths
  - [ ] Rate limiting for public endpoints (optional, not implemented)
- [x] Comprehensive test coverage:
  - [x] Service layer tests (14 test cases)
  - [x] Repository tests (FindByShareToken method)
  - [x] Handler layer tests (mock repository updated)
  - [x] All tests passing

### Frontend Development (~85% Complete)

- [x] Share button in drawing editor:
  - [x] UI component for "Share Public" button (ShareButton.svelte)
  - [x] ShareDialog component with modal UI
  - [x] Generate share URL on click
  - [x] Copy-to-clipboard functionality with visual feedback
  - [x] Success notification ("Link copied to clipboard!")
  - [x] Option to revoke sharing
- [x] Public view page (`/public/[token]`):
  - [x] New route for public share URLs
  - [x] Load drawing using share token (no auth)
  - [x] Full Excalidraw functionality enabled
  - [x] Save/update operations disabled (client-side enforcement)
  - [x] All changes are temporary (session-only)
  - [x] Clear indication of "View-Only" mode with banner
  - [x] Link to landing page (personal-excalidraw.lukenguyen.me)
  - [x] No navigation to other pages (isolated view)
  - [x] API endpoint for fetching public drawings (getByShareToken)
  - [x] PublicExcalidrawWrapper component for view-only mode
  - [x] Error handling for invalid/expired tokens
- [ ] Drawing list page updates:
  - [ ] Visual indicator for publicly shared drawings
  - [ ] Quick copy share URL from list view

### User Experience

- [x] Toast notifications:
  - [x] "Public link copied to clipboard"
  - [x] "Link shared successfully"
  - [x] "Sharing disabled"
- [x] View-only mode banner:
  - [x] Clear messaging about temporary edits
  - [x] Reload warning (changes will be lost)
  - [x] Link to project landing page
- [x] URL structure:
  - [x] Clean public URLs: `/public/[share_token]`
  - [x] Short, shareable tokens (UUID-based)

### Goals

- Enable easy sharing of drawings with friends/colleagues
- Maintain security by keeping drawings view-only
- Allow temporary collaboration/discussion without affecting original
- Promote the open-source project through landing page link

**Outcome**: Users can generate public share URLs for drawings that allow anyone to view and temporarily interact with the drawing without authentication or ability to save changes.

## Future Considerations

## Milestones

| Phase                            | Target         | Status           |
| -------------------------------- | -------------- | ---------------- |
| Phase 1: Frontend Infrastructure | ✅ Complete    | Done             |
| Phase 2: Local Storage           | ✅ Complete    | Done             |
| Phase 3: Backend Integration     | ✅ Complete    | Done             |
| Phase 4: Production Ready        | ✅ Complete    | Done             |
| Phase 5: Public Share URL        | 🚧 In Progress | Backend Complete |
