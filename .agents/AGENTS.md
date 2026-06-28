# Expense Tracker Project Rules

This document outlines the architecture, coding standards, and guidelines for the Expense Tracker project. All agents and assistants must adhere to these rules when modifying or adding code.

## Tech Stack & Architecture

- **Backend / Primary Language**: Go (Golang)
- **Frontend**: Server-Side Rendering (SSR) using Go standard library `html/template`.
- **Interactivity**: HTMX for dynamic, AJAX-like interactions without full page reloads.
- **Styling**: Bootstrap 5 (CSS framework).
- **Database**: Supabase (PostgreSQL).

---

## Coding Guidelines

### 1. Go (Backend)
- Follow standard Go project layout conventions, organized by feature (screaming architecture/feature-based grouping).
- Place business logic in the `internal/` directory to prevent external exposure.
- Use explicit error handling (`if err != nil`). Do not panic or ignore errors.
- Ensure database connections use Supabase-compatible PostgreSQL drivers.

### 2. Frontend (HTML + HTMX + Bootstrap 5)
- Keep UI components clean and structured using Go's standard templates.
- Use HTMX attributes (e.g., `hx-get`, `hx-post`, `hx-target`, `hx-swap`) for dynamic updates instead of custom JavaScript where possible.
- Style components strictly using Bootstrap 5 utility classes and components. Avoid adding custom inline styles.

### 3. Database (Supabase)
- Write clean SQL queries or use structured libraries.
- Respect Supabase Row Level Security (RLS) policies if interacting directly with the client client-side, or route transactions securely through the backend.
