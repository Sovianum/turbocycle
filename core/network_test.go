package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEmptyRoots(t *testing.T) {
	var hasEmptyRoots = linkTableType{
		"1": rowType{"2": true},
		"2": rowType{},
	}

	var emptyRoots = getEmptyRootIds(hasEmptyRoots)
	assert.Equal(t, 1, len(emptyRoots))
	assert.Equal(t, "2", emptyRoots[0])

	var noEmptyRoots = linkTableType{
		"1": {"2": true},
	}
	emptyRoots = getEmptyRootIds(noEmptyRoots)

	assert.Equal(t, 0, len(emptyRoots))
}

func TestGetCallOrder_CyclicDependency(t *testing.T) {
	var requireTree = linkTableType{
		"inlet_source":  {},
		"outlet_source": {},
		"compressor":    {"inlet_source": true},
		"regenerator":   {"compressor": true, "free_turbine": true},
		"burner":        {"regenerator": true},
		"turbine":       {"compressor": true, "burner": true},
		"pressure_loss": {"turbine": true},
		"free_turbine":  {"pressure_loss": true, "outflow": true},
		"outflow":       {"outlet_source": true},
	}

	var updateTree = linkTableType{
		"inlet_source":  {"compressor": true},
		"outlet_source": {"outflow": true},
		"compressor":    {"regenerator": true, "turbine": true},
		"regenerator":   {"burner": true, "outflow": true},
		"burner":        {"turbine": true},
		"turbine":       {"pressure_loss": true},
		"pressure_loss": {"free_turbine": true},
		"free_turbine":  {},
		"outflow":       {"free_turbine": true},
	}

	var _, err = getCallOrder(requireTree, updateTree)
	assert.NotNil(t, err)
}

func TestGetCallOrder_OK(t *testing.T) {
	var requireTree = linkTableType{
		"inlet_source":  {},
		"outlet_source": {},
		"compressor":    {"inlet_source": true},
		"regenerator":   {"compressor": true, "cycle_breaker": true},
		"burner":        {"regenerator": true},
		"turbine":       {"compressor": true, "burner": true},
		"pressure_loss": {"turbine": true},
		"free_turbine":  {"pressure_loss": true, "cycle_breaker": true},
		"outflow":       {"outlet_source": true},
		"cycle_breaker": {},
	}

	var updateTree = linkTableType{
		"inlet_source":  {"compressor": true},
		"outlet_source": {"outflow": true},
		"compressor":    {"regenerator": true, "turbine": true},
		"regenerator":   {"burner": true, "outflow": true},
		"burner":        {"turbine": true},
		"turbine":       {"pressure_loss": true},
		"pressure_loss": {"free_turbine": true},
		"free_turbine":  {},
		"outflow":       {"free_turbine": true},
		"cycle_breaker": {"free_turbine": true, "regenerator": true},
	}

	var items = map[string]bool{
		"inlet_source":  true,
		"outlet_source": true,
		"compressor":    true,
		"regenerator":   true,
		"burner":        true,
		"pressure_loss": true,
		"free_turbine":  true,
		"outflow":       true,
	}

	var callOrder, err = getCallOrder(requireTree, updateTree)
	assert.Nil(t, err)

	for _, item := range callOrder {
		delete(items, item)
	}

	assert.Equal(t, 0, len(items), fmt.Sprintf("Nodes %v has not been called", items))
}
