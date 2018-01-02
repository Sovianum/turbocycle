package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBiMap_Add_OK(t *testing.T) {
	var biMap = NewBiMap()
	var err = biMap.Add(0, "a")
	assert.Nil(t, err)

	var val, valOk = biMap.GetByKey(0)
	assert.True(t, valOk)
	assert.Equal(t, val, "a")

	var key, keyOk = biMap.GetByVal("a")
	assert.True(t, keyOk)
	assert.Equal(t, key, 0)

	assert.Equal(t, 1, biMap.Length())
}

func TestBiMap_Add_Conflict(t *testing.T) {
	var biMap = NewBiMap()
	var err = biMap.Add(0, "a")
	assert.Nil(t, err)

	err = biMap.Add(0, "a")
	assert.NotNil(t, err)

	var val, valOk = biMap.GetByKey(0)
	assert.True(t, valOk)
	assert.Equal(t, val, "a")

	var key, keyOk = biMap.GetByVal("a")
	assert.True(t, keyOk)
	assert.Equal(t, key, 0)

	assert.Equal(t, 1, biMap.Length())
}

func TestBiMap_ContainsKey(t *testing.T) {
	var biMap = NewBiMap()
	biMap.Add(0, "a")

	assert.True(t, biMap.ContainsKey(0))
	assert.False(t, biMap.ContainsKey(1))
}

func TestBiMap_ContainsVal(t *testing.T) {
	var biMap = NewBiMap()
	biMap.Add(0, "a")

	assert.True(t, biMap.ContainsVal("a"))
	assert.False(t, biMap.ContainsVal("n"))
}

func TestBiMap_DeleteByKey(t *testing.T) {
	var biMap = NewBiMap()
	biMap.Add(0, "a")
	biMap.Add(1, "b")
	biMap.DeleteByKey(0)

	assert.True(t, biMap.ContainsKey(1))
	assert.True(t, biMap.ContainsVal("b"))
	assert.False(t, biMap.ContainsKey(0))
	assert.False(t, biMap.ContainsVal("a"))
	assert.Equal(t, 1, biMap.Length())
}

func TestBiMap_DeleteByVal(t *testing.T) {
	var biMap = NewBiMap()
	biMap.Add(0, "a")
	biMap.Add(1, "b")
	biMap.DeleteByVal("a")

	assert.True(t, biMap.ContainsKey(1))
	assert.True(t, biMap.ContainsVal("b"))
	assert.False(t, biMap.ContainsKey(0))
	assert.False(t, biMap.ContainsVal("a"))
	assert.Equal(t, 1, biMap.Length())
}

func TestBiMap_Iterate(t *testing.T) {
	var biMap = NewBiMap()
	biMap.Add(0, "a")
	biMap.Add(1, "b")
	var checkMap = map[int]interface{}{
		0: "a",
		1: "b",
	}

	for pair := range biMap.Iterate() {
		var val, ok = checkMap[pair.Key]
		assert.True(t, ok)
		assert.Equal(t, val, pair.Val)
	}
}
