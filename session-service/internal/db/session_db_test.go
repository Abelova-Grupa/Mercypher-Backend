package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDBUrl(t *testing.T) {
	dbUrl := GetDBUrl()
	require.NotEqual(t, dbUrl, "")
}

func TestConnect(t *testing.T) {
	db := Connect(GetDBUrl())
	require.NotEmpty(t, db)
}
