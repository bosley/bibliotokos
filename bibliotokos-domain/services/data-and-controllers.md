# Data and Controllers

## Overview

This project uses a two-layer pattern for all persistent data:

1. `pkg/datastore` — owns the SQLite file, connection pools, and database-level configuration. It knows nothing about application domain.
2. `pkg/<domain>` — owns its own schema, migrations, and business logic. It receives a `*datastore.DB` at construction time and uses the shared connections.

This means a single `.db` file can back any number of domain controllers without each one re-opening the file or fighting over connection settings.

---

## `pkg/datastore`

`datastore.Open(location string) (*DB, error)` is the single entry point. It configures:

- **WAL journal mode** — allows one concurrent writer and many concurrent readers without blocking each other.
- **Busy timeout (5000ms)** — retries automatically instead of returning `SQLITE_BUSY` immediately.
- **Foreign keys ON** — enforced at the writer connection.
- **Connection pool limits** — writer capped at 1 open connection (SQLite only supports one concurrent writer), reader pool capped at 8.

The returned `*DB` exposes `Writer() *sql.DB` and `Reader() *sql.DB`. Domain packages call these to get their connections. `Close()` shuts both pools down.

```
ds, err := datastore.Open("/var/data/app.db")
defer ds.Close()
```

`datastore` is intentionally schema-agnostic. It never runs `CREATE TABLE`. That responsibility belongs exclusively to each domain package.

---

## Domain controller pattern

### File layout

```
pkg/<domain>/
    record.go       — plain data types (structs, enums, constants)
    errors.go       — sentinel error variables, one per failure mode
    controller.go   — Config struct, public interface, FromDB constructor
    sqlite.go       — unexported struct implementing the interface
    controller_test.go
```

### `record.go`

Plain Go structs and any associated types (roles, statuses, etc.). No methods, no logic. These are the values the controller returns to callers.

### `errors.go`

One exported `var` per distinct failure mode. Callers use `errors.Is` to handle them. Never return raw SQLite errors or `fmt.Errorf` strings to callers — map everything to a sentinel.

```go
var (
    ErrThingNotFound  = errors.New("thing not found")
    ErrThingDuplicate = errors.New("thing already exists")
)
```

### `controller.go`

Three things only:

1. A `Config` struct holding `*datastore.DB` and any domain-specific options (secrets, timeouts, etc.).
2. A public interface declaring every operation the domain exposes.
3. A `FromDB(cfg Config) (TheInterface, error)` function that constructs the implementation.

```go
type Config struct {
    DB      *datastore.DB
    // domain-specific fields
}

type ThingController interface {
    Create(...) (string, error)
    Delete(id string) error
    // ...
}

func FromDB(cfg Config) (ThingController, error) {
    return newSQLiteController(cfg)
}
```

`controller.go` imports `datastore` but nothing SQLite-specific. The interface itself has no knowledge of storage.

### `sqlite.go`

The unexported implementation. Structure:

```go
type sqliteController struct {
    writer *sql.DB
    reader *sql.DB
    cfg    Config
}
```

`newSQLiteController(cfg Config)` does exactly three things in order:

1. Pull `cfg.DB.Writer()` and `cfg.DB.Reader()` into local fields.
2. Run `CREATE TABLE IF NOT EXISTS` for every table this domain owns.
3. Run any additive `ALTER TABLE ... ADD COLUMN` migrations, ignoring "duplicate column name" errors so the same binary can be deployed against an existing database.

All `SELECT` queries go through `c.reader`. All `INSERT`, `UPDATE`, and `DELETE` go through `c.writer`.

SQLite constraint errors (e.g. `UNIQUE constraint failed`) are detected by substring-matching the error string and mapped to the appropriate sentinel before returning to the caller.

`RowsAffected() == 0` after a mutation means the record was not found — return the corresponding sentinel.

### `controller_test.go`

- `TestMain` creates a temporary `.db` file, opens a `datastore.DB`, constructs the controller, runs all tests, then deletes every row whose natural identifier ends with `_test` before calling `os.Exit`.
- Test identifiers always carry the `_test` suffix (e.g. `alice_test`, `delete2_test`) so cleanup is unambiguous and cannot touch real data if a future test somehow runs against a non-ephemeral database.
- One `Test*` function per interface method, covering both the happy path and every defined error sentinel.

---

## Wiring multiple controllers

```go
ds, err := datastore.Open(cfg.DBPath)
if err != nil { ... }
defer ds.Close()

users, err := user.FromDB(user.Config{
    DB:          ds,
    JWTSecret:   cfg.JWTSecret,
    TokenExpiry: 30 * time.Minute,
})

// future domain:
// posts, err := post.FromDB(post.Config{DB: ds})
```

Each controller runs its own schema setup independently on the shared connection. Because `CREATE TABLE IF NOT EXISTS` is idempotent, order does not matter and re-opening an existing database is always safe.

---

## What a new controller must NOT do

- Open its own SQLite connections or reference a file path directly. That is `datastore`'s job.
- Return raw `database/sql` errors or driver error strings to callers. Map to sentinels.
- Mix reads and writes on the same connection. Reads always use `c.reader`, writes always use `c.writer`.
- Perform destructive schema changes (DROP, ALTER COLUMN, RENAME). Migrations are additive only (ADD COLUMN).
