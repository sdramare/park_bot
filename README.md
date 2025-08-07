# park_bot

A Go-based project for automating tasks related to parking management, including integrations with Telegram and Twilio.

## Features

- Telegram bot integration (see `tg.http`)
- Twilio VoIP integration (see `voip/twilio.go`)
- AWS Lambda deployment scripts (see `scripts/`)

## Getting Started

### Prerequisites

- Go 1.24 or newer

### Build

```
go build -v ./...
```

### Test

```
go test -v ./...
```

## Deployment

See scripts in the `scripts/` directory for AWS Lambda deployment.

## License

MIT
