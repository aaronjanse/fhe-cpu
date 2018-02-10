package cpu

import (
	"fmt"
	"strings"
)

func (thisChip hdlChip) Evaluate(inputs map[string]FheType, stack string) map[string]FheType {
	return thisChip.evaluate(inputs, stack)
}

func (thisChip hdlChip) evaluate(inputs map[string]FheType, stack string) map[string]FheType {
	newStack := stack + ">" + thisChip.name
	if verbose {
		fmt.Printf("\r     %-64s", newStack)
	}

	for key := range thisChip.wires {
		thisChip.wires[key] = FheType{IsNil: true}
	}

	for key, value := range inputs {
		thisChip.wires[key] = value
	}

	for idx := range thisChip.subChips {
		thisChip.subChips[idx].dead = false
	}

	for _, subChip := range thisChip.subChips {
		if !subChip.stateful {
			continue
		}

		subInputs := make(map[string]FheType)
		// for key := range subChip.chipIO.inputs {
		// 	subInputs[key] = -1
		// }
		for childWire, parentWire := range subChip.childParentWiresIn {
			subInputs[childWire] = thisChip.wires[parentWire]
		}

		subOutputs := subChip.evaluatableChip.evaluate(subInputs, newStack)

		for childWire, parentWire := range subChip.childParentWiresOut {
			thisChip.wires[parentWire] = subOutputs[childWire]
		}
	}

	for {
		for idx, subChip := range thisChip.subChips {
			if subChip.dead {
				continue
			}

			subInputs := make(map[string]FheType)

			ready := true
			for childWire, parentWire := range subChip.childParentWiresIn {
				parentValue := thisChip.wires[parentWire]
				if !parentValue.IsNil {
					subInputs[childWire] = parentValue
				} else {
					ready = false
					break
				}
			}
			if !ready {
				continue
			}

			subOutputs := subChip.evaluatableChip.evaluate(subInputs, newStack)

			for childWire, parentWire := range subChip.childParentWiresOut {
				thisChip.wires[parentWire] = subOutputs[childWire]
			}

			thisChip.subChips[idx].dead = true
		}

		done := true
		if thisChip.stateful {
			for key := range thisChip.chipIO.outputs {
				if thisChip.wires[key].IsNil {
					done = false
					break
				}
			}
		} else {
			for _, value := range thisChip.wires {
				if value.IsNil {
					done = false
					break
				}
			}
		}

		if done {
			break
		}
	}

	outputs := make(map[string]FheType)

	for key := range thisChip.chipIO.outputs {
		outputs[key] = thisChip.wires[key]
	}

	if stack == "#" && verbose {
		fmt.Print("\r" + strings.Repeat(" ", 64) + "\r")
	}

	return outputs
}
