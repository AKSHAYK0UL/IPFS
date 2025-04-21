## Interplanetary File System

A GitHub-backed transaction service with IPFS-like pooling and blockchain-style validation written in Go. This service stores transactions in pools on a GitHub repository, validates chains via a smart contract approach, and includes a planned snapshot-based rollback mechanism to recover from invalid transaction pools.

---

## Features

- **GitHub Storage**: Transactions are stored as JSON files within a GitHub repo, organized in sequentially numbered pools for scalability and traceability.
- **Blockchain-Style Validation**: Each pool functions as a mini-blockchain, leveraging SHA-256 hashing to ensure data integrity and order.
- **Smart Contract Validation**: A Go-based smart contract helper scans and verifies each pool, detecting any tampering or inconsistency.
- **Snapshot-Based Rollback (Planned)**: A rollback feature is under development that will snapshot repository state before validation and revert on failure.
- **RESTful API with Gin**: A simple, secure HTTP API for adding, querying, and validating transactions.

---

## Technologies

- **Language**: Go 1.18+
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **GitHub Client**: google/go-github
- **Auth**: OAuth2 (golang.org/x/oauth2), dotenv for local env vars
- **Middleware**: Rate limiting and header-based auth

---

## Prerequisites

1. Go 1.18 or later installed (`go version`)
2. A GitHub repository (public or private) to store pools
3. Personal Access Token (PAT) with `repo` scope

---

## Installation & Setup

1. **Clone** the repo:
   ```bash
   git clone https://github.com/<username>/<repo>.git
   cd <repo>
   ```
2. **Dependencies**:
   ```bash
   go mod download
   ```
3. **Build**:
   ```bash
   go build -o tx-service ./...
   ```
4. **Environment**: Create a `.env` file in project root:
   ```dotenv
   ACCESSTOKEN=ghp_<your_token>
   AUTHHEADERKEY=<your_api_key>
   REPO_OWNER=<github_user_or_org>
   REPO_NAME=<repository_name>
   DIR_NAME=pools
   FILE_NAME=transactions.json
   MAX_TXN=5
   ```

---

## Running the Service

```bash
./tx-service
```

By default, the server listens on port `8080`. Customize via environment variables or flags as needed.

---

## API Reference

### Authentication

All modifying endpoints require an `Authorization` header:

```
Authorization: <AUTHHEADERKEY>
```

Missing or invalid headers return `401 Unauthorized`.

### Endpoints

| Method | Path                   | Description                                   |
|--------|------------------------|-----------------------------------------------|
| POST   | `/add-transaction`     | Add a new transaction to the current pool.    |
| GET    | `/transactions`        | Retrieve all transactions across all pools.   |
| GET    | `/transaction/:id`     | Retrieve a specific transaction by its ID.    |
| GET    | `/smart-contract`      | Run validation and report invalid entries.    |

#### Add Transaction

- **URL**: `/add-transaction`
- **Body** (application/json):
  ```json
  {
    "txn_id": "tx123",
    "to_id": "userB",
    "from_id": "userA",
    "amount": 100.5,
    "nonce": 1,
    "time": "2025-04-21T10:00:00Z"
  }
  ```
- **Responses**:
  - `200 OK` – transaction accepted
  - `400 Bad Request` – malformed body
  - `401 Unauthorized` – missing/invalid header
  - `500 Internal Server Error` – GitHub or internal error

#### Get All Transactions

- **URL**: `/transactions`
- **Response**: `200 OK` with JSON array of `IPFSTransaction` objects.

#### Get Transaction by ID

- **URL**: `/transaction/:id`
- **Response**:
  - `200 OK` with single object when found
  - `404 Not Found` if the ID does not exist

#### Smart Contract Validation

- **URL**: `/smart-contract`
- **Response**:
  - `200 OK` with `{ "message": "chain is valid" }` if no issues
  - `400 Bad Request` with `{ "invalid_transactions": [...] }` listing detected invalid pools

---

## Snapshot-Based Rollback (Under Development)

The rollback mechanism ensures that any detected corruption or invalid transaction pool can be reverted automatically:

1. **Snapshot Creation**: Before validation, the service will create a Git tag (e.g., `snapshot-20250421-143000`) capturing the current `main` branch reference.
2. **Validation Phase**: Pools are scanned; any inconsistency triggers rollback logic.
3. **Rollback Execution**:
   - Use GitHub Tags & Refs API to reset `main` branch to the snapshot tag commit SHA.
   - Discard all commits made after the snapshot.
4. **Error Reporting**: The API response will include the invalid transaction details and a rollback status.

### Development TODOs

- [ ] Implement `CreateSnapshotTag(ctx)` in `githubdb`:
  - Calls `Repositories.CreateRef` with `refs/tags/<snapshot>` pointing at current HEAD.
- [ ] Implement `RollbackToTag(ctx, tag)`:
  - Fetch tag SHA via `Repositories.GetRef`, then `Repositories.UpdateRef` on `refs/heads/main` with `force=true`.
- [ ] Add unit tests for happy-path and failure rollback.
- [ ] Extend API to expose snapshot and rollback endpoints for manual control.

---
