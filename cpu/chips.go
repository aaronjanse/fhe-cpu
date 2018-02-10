package cpu

import (
	"fmt"
)

var builtinChips = make(map[string]builtinChip)

type builtinChip struct {
	name string
	chipIO
}

func init() {
	builtinChips["nand"] = builtinChip{name: "nand", chipIO: chipIOShorthand([]string{"a", "b"}, []string{"out"})}
	builtinChips["and"] = builtinChip{name: "and", chipIO: chipIOShorthand([]string{"a", "b"}, []string{"out"})}
	builtinChips["not"] = builtinChip{name: "not", chipIO: chipIOShorthand([]string{"in"}, []string{"out"})}
	builtinChips["bind"] = builtinChip{name: "bind", chipIO: chipIOShorthand([]string{"in"}, []string{"out"})}
	builtinChips["or"] = builtinChip{name: "or", chipIO: chipIOShorthand([]string{"a", "b"}, []string{"out"})}
	builtinChips["xor"] = builtinChip{name: "xor", chipIO: chipIOShorthand([]string{"a", "b"}, []string{"out"})}
	builtinChips["mux"] = builtinChip{name: "mux", chipIO: chipIOShorthand([]string{"a", "b", "sel"}, []string{"out"})}
}

type evaluatableChip interface {
	evaluate(map[string]FheType, string) map[string]FheType
}

type chipIO struct {
	inputs  map[string]bool
	outputs map[string]bool
}

func chipIOShorthand(inputsList, outputsList []string) chipIO {
	inputs := make(map[string]bool)
	for _, name := range inputsList {
		inputs[name] = true
	}
	outputs := make(map[string]bool)
	for _, name := range outputsList {
		outputs[name] = true
	}

	return chipIO{inputs: inputs, outputs: outputs}
}

type hdlChip struct {
	chipIO
	name     string
	stateful bool
	wires    map[string]FheType
	subChips []subChip
}

type subChip struct {
	chipIO
	evaluatableChip
	name                string
	stateful            bool
	dead                bool
	childParentWiresIn  map[string]string // inputs
	childParentWiresOut map[string]string // outputs
}

func (chip builtinChip) evaluate(inputs map[string]FheType, stack string) map[string]FheType {
	if verbose {
		fmt.Printf("\r%-5s", chip.name)
	}

	switch chip.name {
	case "nand":
		return map[string]FheType{"out": fheNand(inputs["a"], inputs["b"])}
	case "and":
		return map[string]FheType{"out": fheAnd(inputs["a"], inputs["b"])}
	case "not":
		return map[string]FheType{"out": fheNot(inputs["in"])}
	case "or":
		return map[string]FheType{"out": fheOr(inputs["a"], inputs["b"])}
	case "xor":
		return map[string]FheType{"out": fheXor(inputs["a"], inputs["b"])}
	case "mux":
		return map[string]FheType{"out": fheMux(inputs["sel"], inputs["b"], inputs["a"])}
	case "bind":
		return map[string]FheType{"out": fheCopy(inputs["in"])}
	default:
		return nil
	}
}

type registerChip struct {
	chipIO
	state    FheType
	stateful bool
}

func newRegisterChip() *registerChip {
	return &registerChip{chipIO: chipIOShorthand([]string{"should_set", "value"}, []string{"out"}), stateful: true, state: NewConstant(0)}
}

func (chip *registerChip) evaluate(inputs map[string]FheType, stack string) map[string]FheType {
	oldState := chip.state
	if !inputs["should_set"].IsNil && !inputs["value"].IsNil {
		chip.state = fheMux(inputs["should_set"], inputs["value"], chip.state)
	}
	return map[string]FheType{"out": oldState}
}
