## Go Fortune

### Overview
A small Go service that fetches the latest GitHub release for a given repository and stores it in PostgreSQL. On each run, it will:
- Create the `packages` table if it does not exist
- Fetch the latest release via the GitHub API
- Insert the repo/version if missing
- Update the stored version only if a newer semantic version is available

The default example tracks `snowflakedb/snowpark-python` (see `server.go`).

### Tech Stack
- **Language**: Go 1.22
- **Database**: PostgreSQL 16
- **Deps**: `github.com/joho/godotenv`, `github.com/lib/pq`, `github.com/hashicorp/go-version`

### Project Structure
- `server.go`: Program entry. Loads env, connects DB, creates table, fetches release, inserts/updates version
- `config/`: Environment and Postgres config loaders
- `database/`: DB connection and queries (create table, insert, get, update)
- `models/`: Data models for responses and updates
- `requests/`: GitHub API client for fetching latest release
- `main.go`: Legacy/local sample code (not used by `server.go`)
- `docker-compose.yml`: Dev stack for Postgres and running the app

### Prerequisites
- Go 1.22+
- Docker (if using Docker Compose)

### Environment Variables
Create a `.env` file in the project root with:

```
DB_HOST=postgres             # or "db" when using docker-compose
DB_PORT=5432                 
DB_PASSWORD=mysecretpassword
DB_NAME=mydb
```

Notes:
- When running with the provided `docker-compose.yml`, `DB_HOST` should be `db` (the service name) since the app connects to `postgres://<host>:<password>@db:<port>/<name>?sslmode=disable`.
- When running locally against a local Postgres, set `DB_HOST=localhost`.

### Running with Docker Compose (recommended for dev)
1) Create `.env` as above.
2) Start the stack:

```bash
docker compose up --build
```

The app container will run `go run server.go`, connect to Postgres, and log the inserted/updated version.

### Running Locally (without Docker)
1) Ensure Postgres is running and accessible. Create `.env` with local values (e.g., `DB_HOST=localhost`).
2) Run the program:

```bash
go run server.go
```

### Changing the Tracked Repository
Edit the request in `server.go`:

```go
request := requests.Request{Owner: "owner", Repo: "repository"}
```

The program will fetch the latest release for that repository and upsert the version accordingly.

### Database Schema
The table is created automatically if missing:

```sql
CREATE TABLE IF NOT EXISTS packages (
  id SERIAL PRIMARY KEY,
  owner VARCHAR(100) NOT NULL,
  repo VARCHAR(100) NOT NULL,
  version VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
```

### Troubleshooting
- Ensure env vars are loaded; missing values will cause startup errors
- If running via Docker, verify both `db` and `app` services are up
- GitHub API rate limits may apply for unauthenticated requests



