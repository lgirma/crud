package crud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomStr(t *testing.T) {
	rando := GetRandomStr(4)
	rando2 := GetRandomStr(4)
	assert.Equal(t, 4, len(rando))
	assert.Equal(t, 4, len(rando2))
	assert.NotEqual(t, rando, rando2)
}