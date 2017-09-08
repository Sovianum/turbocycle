package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEmptyRoots(t *testing.T) {
	var hasEmptyRoots = linkTableType{
		1: rowType{2: true},
		2: rowType{},
	}

	var emptyRoots = getEmptyRootIds(hasEmptyRoots)
	assert.Equal(t, 1, len(emptyRoots))
	assert.Equal(t, 2, emptyRoots[0])

	var noEmptyRoots = linkTableType{
		1: rowType{2: true},
	}
	emptyRoots = getEmptyRootIds(noEmptyRoots)

	assert.Equal(t, 0, len(emptyRoots))
}

func TestGetCallOrder_CyclicDependency(t *testing.T) {
	var requireTree = linkTableType{
		0: {},
		1: {0: true},
		2: {1: true, 5: true, 6: true},
		3: {2: true},
		4: {1: true, 3: true},
		5: {2: true, 4: true},
		6: {},
	}

	var updateTree = linkTableType{
		0: {1: true},
		1: {2: true, 4: true},
		2: {3: true, 5: true},
		3: {4: true},
		4: {5: true},
		5: {2: true},
		6: {2: true},
	}

	var _, err = getCallOrder(requireTree, updateTree)
	assert.NotNil(t, err)

	fmt.Println(err)
}

func TestGetCallOrder_OK(t *testing.T) {
	var requireTree = linkTableType{
		0: {},
		1: {0: true},
		2: {1: true, 6: true, 7: true},
		3: {2: true},
		4: {1: true, 3: true},
		5: {6: true, 4: true},
		6: {},
		7: {},
	}

	var updateTree = linkTableType{
		0: {1: true},
		1: {2: true, 4: true},
		2: {3: true, 6: true},
		3: {4: true},
		4: {5: true},
		5: {6: true},
		6: {2: true, 5: true},
		7: {2: true},
	}

	var order, err = getCallOrder(requireTree, updateTree)
	assert.Nil(t, err)

	fmt.Println(order)
}
