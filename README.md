# Ciplock
Ciplock is a robust backend service built with Go and the Gin framework, designed to provide secure, scalable authentication and project management for multi-tenant SaaS applications.
# Ciplock

**Ciplock** is a secure, scalable API platform designed to manage user authentication, project management, and customer onboarding workflows for multi-tenant SaaS applications.

---

## Table of Contents
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Authentication & Authorization](#authentication--authorization)
- [Environment Variables](#environment-variables)
- [Deployment](#deployment)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- Admin user signup and authentication with bcrypt password hashing
- Project creation and management per admin user
- Customer signup and validation scoped to projects
- JWT-based authentication with middleware protection
- Structured error handling with meaningful HTTP status codes
- Modular architecture for extensibility

---

## Technology Stack

- **Language:** Go (Golang)
- **Web Framework:** Gin-Gonic
- **Authentication:** JWT tokens, bcrypt for password hashing
- **Database:** PostgreSQL (assumed; adjust as per your actual DB)
- **UUID:** github.com/google/uuid for unique identifiers

---

## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL or preferred SQL DB
- `make` (optional, for build automation)
- Git

### Installation

```bash
git clone https://github.com/yourusername/ciplock.git
cd ciplock
go mod download
