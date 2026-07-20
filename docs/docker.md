# Docker

This project includes a complete Docker development and production environment.

## Features

- Multi-stage production image
- Docker Compose
- PostgreSQL 17
- pgAdmin 4
- Air hot reload (development)
- Automatic database migrations
- Automatic database seeding
- Persistent PostgreSQL volumes
- Development and production Compose configurations

---

# Project Structure

```text
deployments/
└── docker/
    ├── Dockerfile
    ├── Dockerfile.dev
    ├── docker-compose.yml
    ├── docker-compose.dev.yml
    ├── docker-compose.prod.yml
    ├── entrypoint.sh
    └── pgadmin/
```

---

# Requirements

- Docker
- Docker Compose
- Task (recommended)

Verify installation:

```bash
docker --version
docker compose version
task --version
```

---

# Environment

Create a local environment file.

```bash
cp .env.example .env
```

Update values as needed.

---

# Development

Start the development environment.

```bash
task docker:dev
```

or rebuild it.

```bash
task docker:dev:build
```

The development environment includes:

- PostgreSQL
- pgAdmin
- Go API
- Air hot reload

Whenever a Go file changes, the application is rebuilt and restarted automatically.

---

# Production

Run the production stack.

```bash
task docker:prod
```

Build a fresh production image.

```bash
task docker:prod:build
```

The production image:

- Uses a multi-stage Docker build
- Runs database migrations
- Runs database seeders
- Starts the API

---

# pgAdmin

Open:

http://localhost:5050

Default credentials:

Email:

```
admin@example.com
```

Password:

```
admin
```

Once logged in, connect to PostgreSQL using:

Host

```
postgres
```

Port

```
5432
```

Database

```
go_api_db
```

Username

```
postgres
```

Password

```
postgres
```

---

# API

The API is available at

http://localhost:8080

---

# Useful Commands

Start development

```bash
task docker:dev
```

Rebuild development

```bash
task docker:dev:build
```

Start production

```bash
task docker:prod
```

Rebuild production

```bash
task docker:prod:build
```

View logs

```bash
task docker:logs
```

Stop containers

```bash
task docker:down
```

List running containers

```bash
task docker:ps
```

Reset PostgreSQL volume

```bash
task docker:clean
```

---

# Database

The PostgreSQL data directory is stored in a Docker volume.

Deleting containers does not remove the database.

To completely reset the database:

```bash
task docker:clean
```

---

# Hot Reload

Development mode uses Air.

Saving any Go file automatically:

1. Rebuilds the application
2. Restarts the API
3. Keeps PostgreSQL running

No manual rebuild is required.

---

# Migrations

Production automatically runs Goose migrations before starting the API.

Development expects migrations to already exist and focuses on rapid iteration.

---

# Seeders

Production automatically executes project seeders after migrations.

Seeders are idempotent and can safely run multiple times.

---

# Troubleshooting

## Port already in use

Example:

```
bind: address already in use
```

Another service is already using the port.

Either stop that service or change the exposed port in Docker Compose.

---

## PostgreSQL hostname cannot be resolved

Inside Docker, always use:

```
DB_HOST=postgres
```

Never use:

```
localhost
```

---

## Reset everything

```bash
task docker:clean
task docker:dev:build
```

---

# Future Improvements

The Docker infrastructure is designed to support future additions including:

- Redis
- Mailpit
- MinIO
- Nginx
- Prometheus
- Grafana
- Kubernetes deployment