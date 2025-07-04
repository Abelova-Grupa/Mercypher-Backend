package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests don't work because of the change of LoadEnv directory
func TestGetDBUrl(t *testing.T) {
	dbUrl := GetDBUrl()
	require.NotEqual(t, dbUrl, "")
}

func TestConnect(t *testing.T) {
	db := Connect()
	require.NotEmpty(t, db)
}
