package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFromEnvDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("NEXT_TAKINGS_MIN", "")

	cfg, err := NewFromEnv()
	require.NoError(t, err)
	require.Equal(t, "8080", cfg.Port)
	require.Equal(t, "medicine.db", cfg.DBPath)
	require.Equal(t, 60, cfg.NextTakingsMin)
}

func TestNewFromEnv(t *testing.T) {
	t.Setenv("PORT", "3030")
	t.Setenv("DB_PATH", "other.db")
	t.Setenv("NEXT_TAKINGS_MIN", "120")

	cfg, err := NewFromEnv()
	require.NoError(t, err)
	require.Equal(t, "3030", cfg.Port)
	require.Equal(t, "other.db", cfg.DBPath)
	require.Equal(t, 120, cfg.NextTakingsMin)
}

func TestNewFromEnvInvalidValue(t *testing.T) {
	t.Setenv("NEXT_TAKINGS_MIN", "abc")

	_, err := NewFromEnv()
	require.Error(t, err)
}

func TestNewFromEnvNegativeValue(t *testing.T) {
	t.Setenv("NEXT_TAKINGS_MIN", "-10")

	_, err := NewFromEnv()
	require.Error(t, err)
}

func TestNewFromEnvZeroValue(t *testing.T) {
	t.Setenv("NEXT_TAKINGS_MIN", "0")

	_, err := NewFromEnv()
	require.Error(t, err)
}
