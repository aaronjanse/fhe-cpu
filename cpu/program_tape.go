package cpu

var programTapes []*programTapeChip

func InitTapes(program [3][]FheType) {
	for i := 0; i < 3; i++ {
		newChip := &programTapeChip{chipIO: chipIOShorthand([]string{"direction"}, []string{"out"}), stateful: true, state: program[i]}
		programTapes = append(programTapes, newChip)
	}
}

type programTapeChip struct {
	chipIO
	state    []FheType
	stateful bool
}

func (chip *programTapeChip) evaluate(inputs map[string]FheType, stack string) map[string]FheType {
	output := chip.state[0]

	if !inputs["direction"].IsNil {
		newState := make([]FheType, 0)
		for idx := range chip.state {
			var leftBit FheType
			if idx == 0 {
				leftBit = chip.state[len(chip.state)-1]
			} else {
				leftBit = chip.state[idx-1]
			}
			var rightBit FheType
			if idx == len(chip.state)-1 {
				rightBit = chip.state[0]
			} else {
				rightBit = chip.state[idx+1]
			}

			thisBit := fheMux(inputs["direction"], rightBit, leftBit)
			newState = append(newState, thisBit)
		}

		chip.state = newState
	}

	return map[string]FheType{"out": output}
}
