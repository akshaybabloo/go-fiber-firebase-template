# Fiber with Firebase Template

A template for Go's [Fiber v3](https://gofiber.io/) framework with Firebase
(Auth + Firestore).

Look for `TODO` markers — those are the spots you need to fill in.

## Features

- Firebase Auth ID-token verification middleware (always enabled)
- Firestore client, initialised once and reused
- RFC 9457 (`application/problem+json`) error responses
- `helmet`, `cors`, `recover`, `requestid`, and structured `zap` logging
- Graceful shutdown on `SIGINT` / `SIGTERM`

## Configuration

`config.yaml` holds non-secret settings:

| Key              | Description                                          |
| ---------------- | ---------------------------------------------------- |
| `debug`          | Development mode + uses the local Firestore emulator |
| `firebaseApiKey` | Your Firebase web API key (if needed client-side)    |
| `version`        | App version                                          |

Environment variables:

| Variable               | Description                         |
| ---------------------- | ----------------------------------- |
| `PORT`                 | Port to listen on (default `3000`)  |
| `FIREBASE_CREDENTIALS` | Path to a service-account JSON file |

## Credentials

Service-account credentials are **not** committed or embedded. Provide them at
runtime via either:

1. `FIREBASE_CREDENTIALS=/path/to/service-account.json`, or
2. Application Default Credentials (e.g. `GOOGLE_APPLICATION_CREDENTIALS`, or the
   metadata server when running on Google Cloud).

See [config/firebase_config.example.json](config/firebase_config.example.json)
for the expected shape. Keep real credentials out of version control — they are
covered by `.gitignore`.

## Running

In debug mode (`debug: true`) the app expects the Firestore emulator to be
running at `localhost:8080`. Start it first, e.g.:

```bash
firebase emulators:start --only firestore
```

If you are not using the emulator, set `debug: false` in `config.yaml` and
provide real credentials (see [Credentials](#credentials)).

```bash
# debug mode (requires the Firestore emulator on localhost:8080)
go run .

# tests
go test ./...
```
