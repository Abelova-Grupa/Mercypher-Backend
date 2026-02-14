package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	db := Connect()
	require.NotEmpty(t, db)
}
