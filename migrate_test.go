package migrate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDriver is a simple in-memory driver for testing purposes.
type mockDriver struct {
	version    int
	dirty      bool
	applied    []int
	forceErr   error
	runErr     error
}

func (m *mockDriver) Open(url string) (Driver, error) {
	return m, nil
}

func (m *mockDriver) Close() error {
	return nil
}

func (m *mockDriver) Lock() error {
	return nil
}

func (m *mockDriver) Unlock() error {
	return nil
}

func (m *mockDriver) Run(migration io.Reader) error {
	if m.runErr != nil {
		return m.runErr
	}
	return nil
}

func (m *mockDriver) SetVersion(version int, dirty bool) error {
	if m.forceErr != nil {
		return m.forceErr
	}
	m.version = version
	m.dirty = dirty
	m.applied = append(m.applied, version)
	return nil
}

func (m *mockDriver) Version() (int, bool, error) {
	return m.version, m.dirty, nil
}

func (m *mockDriver) Drop() error {
	m.version = NilVersion
	m.dirty = false
	m.applied = nil
	return nil
}

func TestNew(t *testing.T) {
	t.Run("valid source and database", func(t *testing.T) {
		mig, err := New("file://testdata/migrations", "stub://")
		require.NoError(t, err)
		assert.NotNil(t, mig)
	})

	t.Run("invalid source URL", func(t *testing.T) {
		_, err := New("invalid://", "stub://")
		assert.Error(t, err)
	})

	t.Run("invalid database URL", func(t *testing.T) {
		_, err := New("file://testdata/migrations", "invalid://")
		assert.Error(t, err)
	})
}

func TestMigrate_Up(t *testing.T) {
	t.Run("applies all pending migrations", func(t *testing.T) {
		mig, err := New("file://testdata/migrations", "stub://")
		require.NoError(t, err)

		err = mig.Up()
		assert.NoError(t, err)
	})

	t.Run("no pending migrations returns ErrNoChange", func(t *testing.T) {
		mig, err := New("file://testdata/migrations", "stub://")
		require.NoError(t, err)

		// Apply all migrations first
		_ = mig.Up()

		// Running Up again should yield ErrNoChange
		err = mig.Up()
		assert.True(t, errors.Is(err, ErrNoChange))
	})
}

func TestMigrate_Down(t *testing.T) {
	t.Run("rolls back all applied migrations", func(t *testing.T) {
		mig, err := New("file://testdata/migrations", "stub://")
		require.NoError(t, err)

		require.NoError(t, mig.Up())

		err = mig.Down()
		assert.NoError(t, err)
	})
}

func TestMigrate_Steps(t *testing.T) {
	t.Run("applies N steps up", func(t *testing.T) {
		mig, err := New("file://testdata/migrations", "stub://")
		require.NoError(t, err)

		err = mig.Steps(1)
		assert.NoError(t, err)
	})

	t.Run("rolls back N steps down", func(t *testing.T) {
		mig, err := New("file://testdata/migrations", "stub://")
		require.NoError(t, err)

		require.NoError(t, mig.Up())

		err = mig.Steps(-1)
		assert.NoError(t, err)
	})

	t.Run("zero steps returns ErrNoChange", func(t *testing.T) {
		mig, err := New("file://testdata/migrations", "stub://")
		require.NoError(t, err)

		err = mig.Steps(0)
		assert.True(t, errors.Is(err, ErrNoChange))
	})
}

func TestMigrate_Version(t *testing.T) {
	t.Run("returns NilVersion when no migrations applied", func(t *testing.T) {
		mig, err := New("file://testdata/migrations", "stub://")
		require.NoError(t, err)

		version, dirty, err := mig.Version()
		require.NoError(t, err)
		assert.Equal(t, uint(NilVersion), version)
		assert.False(t, dirty)
	})
}
