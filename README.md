# Grant Tool

A grant management application for freelance grant writers.

## Tech Stack
- **Backend**: Go, Chi, sqlc, PostgreSQL
- **Frontend**: Next.js, TypeScript, Tailwind, shadcn/ui
- **Infrastructure**: AWS

## Development Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker

### Backend
```bash
docker compose up -d
make migrate-up
make run
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## Status
Active development. Backend API complete for core entities (auth, clients, grants, applications) and deployed to AWS. Frontend basic views complete and deployed to Vercel.