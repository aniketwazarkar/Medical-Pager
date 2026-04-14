# Medical Pager

A production-grade, multi-tenant, white-label medical communication platform designed for doctors, hospitals, and healthcare staff.

## Tech Stack
- **Frontend**: React, TypeScript, Vite, Zustand, DOM Router, CSS Variables (for Tenant White-Labeling).
- **Backend**: Golang, Fiber REST APIs, Fiber WebSockets, Clean Architecture.
- **Database / Infrastructure**: MongoDB Atlas, Redis (Pub/Sub for scaled socket sync).

## Features
- **Multi-Tenant Isolation**: Strict data access limits defined at query-level and schema level using MongoDB indexes.
- **E2E Encryption Support**: Provides symmetric Go AES-GCM encryption utilities on the backend to enforce compliance while supporting future native client-side E2E.
- **Real-Time Delivery**: Native WebSockets with scalable cross-node message distribution powered by Redis.
- **Medical Specifics**: Ready for Patient context linking and Audio/Video communication channels via Agora placeholders.
- **Role-Based Access Control**: Granular auth middleware managing super admins, tenant admins, doctors, nurses, and staff.

## Setup Instructions

### Pre-requisites
- Go 1.22+
- Node.js 20+
- MongoDB instance (Local or Atlas)
- Redis instance

### 1. Environment Configuration
Copy the environment template in the root directory:
```bash
cp .env.example backend/.env
```
Ensure you set up `MONGODB_URI` and proper `ENCRYPTION_KEY` (must be 32 bytes).

### 2. Start Backend Layer
Navigate to backend:
```bash
cd backend
go mod tidy
go run ./cmd/server
```
The server will start on `localhost:5000` and initialize all necessary MongoDB indexes automatically.

### 3. Start Frontend Layer
Navigate to frontend:
```bash
cd frontend
npm install
npm run dev
```
The client application will start on `localhost:5173`.

## Architecture Details
- Data models isolate tightly using `TenantID`.
- Audit logs track critical actions like login patterns and decryption events.
- An `internal/integrations/fhir` stub exists to begin connecting patient data seamlessly to the EMR.
