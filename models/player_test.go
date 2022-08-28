package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayer(t *testing.T) {
	db, cleanup := ConnectWithTestDB()
	defer cleanup()

	player, _ := NewPlayer(db, "Mikey")
	assert.True(t, player.ID > 0)
	assert.Equal(t, player.Name, "Mikey")
	assert.NotNil(t, player.ServerEvents)
}
