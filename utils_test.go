package crud

import (
	"fmt"
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

func TestDefaultIdGenerator(t *testing.T) {
	idGen := CreateNewIdGenerator[string]()

	rando := idGen.GetNewId()
	rando2 := idGen.GetNewId()
	assert.NotEqual(t, rando, rando2)

	idGenInt := CreateNewIdGenerator[int32]()

	randoInt := idGenInt.GetNewId()
	randoInt2 := idGenInt.GetNewId()
	fmt.Printf("%d, %d", randoInt, randoInt2)
	assert.NotEqual(t, randoInt, randoInt2)
}

func TestParse(t *testing.T) {
	parsedInt32 := Parse[int32]("34")
	assert.Equal(t, int32(34), parsedInt32)

	parsedInt64 := Parse[int64]("34")
	assert.Equal(t, int64(34), parsedInt64)

	parsedBool := Parse[bool]("true")
	assert.Equal(t, true, parsedBool)

	parsedStr := Parse[string]("str")
	assert.Equal(t, "str", parsedStr)

	parsedInvalid := Parse[int32]("invalid_number")
	assert.Equal(t, int32(0), parsedInvalid)
}