# Task Platform Worker

Task Platform Worker is a Go service responsible for consuming task events from NATS and processing them asynchronously.

The worker updates task status in PostgreSQL and executes task-specific business logic.

## Tech Stack

* Go 1.25+
* PostgreSQL
* NATS

## Features

* Subscribe to task events
* Parse task payloads
* Process asynchronous jobs
* Update task lifecycle status
* PostgreSQL integration

## Architecture

```text
Task Platform API
      ↓
     NATS
      ↓
Task Platform Worker
      ↓
 PostgreSQL
```

## Message Format

```json
{
  "task_id": "019ea58c-9235-70d1-8dd3-3158ba789400",
  "type": "email.send",
  "payload": {
    "to": "fahmi@example.com"
  }
}
```

## Current Workflow

```text
Task Created
      ↓
pending
      ↓
Worker Receives Message
      ↓
processing
      ↓
completed
```

## Project Structure

```text
cmd/
└── worker/

internal/
├── repository/
├── task/
└── worker/
```

## Local Development

### Install Dependencies

```bash
go mod tidy
```

### Run Worker

```bash
go run ./cmd/worker
```

## Example Output

```text
task_id=019ea58c-9235-70d1-8dd3-3158ba789400 type=email.send
task_id=019ea58c-9235-70d1-8dd3-3158ba789400 processing
task_id=019ea58c-9235-70d1-8dd3-3158ba789400 completed
```

## Roadmap

### Phase 1

* [x] NATS subscriber
* [x] Task processing
* [x] PostgreSQL integration

### Phase 2

* [ ] JetStream
* [ ] Durable Consumer
* [ ] Acknowledgements

### Phase 3

* [ ] Retry mechanism
* [ ] Dead Letter Queue (DLQ)

### Phase 4

* [ ] Worker Pool
* [ ] Concurrency control
* [ ] Metrics & monitoring
