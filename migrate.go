// Package migrate provides database migration functionality.
// It is a fork of golang-migrate/migrate with additional features and fixes.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// ErrNoChange is returned when no migration is needed.
var ErrNoChange = errors.New("no change")

// ErrNilVersion is returned when the version is nil.
var ErrNilVersion = errors.New("nil version")

// ErrLocked is returned when the migration lock cannot be acquired.
var ErrLocked = errors.New("database locked")

// ErrLockTimeout is returned when the lock acquisition times out.
var ErrLockTimeout = errors.New("timeout: could not acquire database lock")

// DefaultPrefetchMigrations is the number of migrations to prefetch.
const DefaultPrefetchMigrations = 10

// DefaultLockTimeout is the default timeout for acquiring a lock in seconds.
const DefaultLockTimeout = 15

// Migrate is the main struct that holds the migration state.
type Migrate struct {
	// sourceName is the name of the source driver.
	sourceName string
	// sourceDrv is the source driver instance.
	sourceDrv Source

	// databaseName is the name of the database driver.
	databaseName string
	// databaseDrv is the database driver instance.
	databaseDrv Database

	// Log is an optional logger.
	Log Logger

	// GracefulStop is a channel to signal graceful stop.
	GracefulStop chan bool
	gracefulStopErr error

	isGracefulStop bool
	isLockedMu     sync.Mutex
	isLocked       bool

	// PrefetchMigrations is the number of migrations to prefetch.
	PrefetchMigrations uint

	// LockTimeout is the timeout for acquiring a lock.
	LockTimeout uint
}

// Logger is the interface for logging migration progress.
type Logger interface {
	Printf(format string, v ...interface{})
	Verbose() bool
}

// New creates a new Migrate instance from source and database URLs.
func New(sourceURL, databaseURL string) (*Migrate, error) {
	m := &Migrate{
		GracefulStop:       make(chan bool, 1),
		PrefetchMigrations: DefaultPrefetchMigrations,
		LockTimeout:        DefaultLockTimeout,
	}

	sourceDrv, err := Open(sourceURL)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	m.sourceDrv = sourceDrv

	databaseDrv, err := OpenDatabase(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	m.databaseDrv = databaseDrv

	return m, nil
}

// Close closes the source and database connections.
func (m *Migrate) Close() (source error, database error) {
	dbErr := m.databaseDrv.Close()
	srcErr := m.sourceDrv.Close()
	return srcErr, dbErr
}

// Migrate applies all pending migrations.
func (m *Migrate) Migrate(version uint) error {
	if version == 0 {
		return ErrNilVersion
	}
	return m.migrate(int(version))
}

// Up applies all available up migrations.
func (m *Migrate) Up() error {
	return m.migrate(-1)
}

// Down reverts all applied migrations.
func (m *Migrate) Down() error {
	return m.migrate(-2)
}

// Steps applies n migrations. Positive n moves up, negative n moves down.
func (m *Migrate) Steps(n int) error {
	if n == 0 {
		return ErrNoChange
	}
	return m.migrate(n)
}

// Version returns the currently active migration version.
// If no migration has been applied, it returns ErrNilVersion.
func (m *Migrate) Version() (version uint, dirty bool, err error) {
	v, d, err := m.databaseDrv.Version()
	if err != nil {
		return 0, false, err
	}
	if v == -1 {
		return 0, false, ErrNilVersion
	}
	return uint(v), d, nil
}

// migrate is the internal migration runner.
func (m *Migrate) migrate(limit int) error {
	if !m.lock() {
		return ErrLocked
	}
	defer m.unlock()

	_ = limit
	// TODO: implement migration execution logic
	return nil
}

func (m *Migrate) lock() bool {
	m.isLockedMu.Lock()
	defer m.isLockedMu.Unlock()
	if m.isLocked {
		return false
	}
	m.isLocked = true
	return true
}

func (m *Migrate) unlock() {
	m.isLockedMu.Lock()
	defer m.isLockedMu.Unlock()
	m.isLocked = false
}

// logVerbosePrintf logs a message if verbose logging is enabled.
func (m *Migrate) logVerbosePrintf(format string, v ...interface{}) {
	if m.Log != nil && m.Log.Verbose() {
		m.Log.Printf(format, v...)
	}
}

// logPrintf logs a message unconditionally.
func (m *Migrate) logPrintf(format string, v ...interface{}) {
	if m.Log != nil {
		m.Log.Printf(format, v...)
	} else {
		fmt.Fprintf(os.Stderr, format, v...)
	}
}
