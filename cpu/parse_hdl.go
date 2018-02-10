package cpu

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var chipRegex = regexp.MustCompile(`(?ms)(?P<stateful>(?:STATEFUL )?)CHIP (?P<name>.+?) {.+?IN +(?P<inputs>.+).+?OUT +(?P<outputs>.+)PARTS:(?P<parts>.+)}`)
var inputsRegex = regexp.MustCompile(`(?ms)\n? *([^,;]+)[,;]`)
var outputsRegex = regexp.MustCompile(`(?ms)\n? *([^,;]+)[,;]`)
var partLineRegex = regexp.MustCompile(`(?ms)([^(]+\([^)]+\))[^\n]*`)
var partNameRegex = regexp.MustCompile(`(?ms).*\n\s*([^(]+)`)
var partChipWireRegex = regexp.MustCompile(`(?ms)[(,]\s*([^=]+)=`)
var partThisWireRegex = regexp.MustCompile(`(?ms)=([^),]+)[),]`)

func ParseHdl(hdl string) (hdlChip, error) {
	chipMatches := chipRegex.FindStringSubmatch(hdl)

	stateful := chipMatches[1] == "STATEFUL "
	name := chipMatches[2]

	wires := make(map[string]FheType)

	inputsList := findMatches(inputsRegex, chipMatches[3])
	inputs := make(map[string]bool)
	for _, key := range inputsList {
		inputs[key] = true
		wires[key] = FheType{IsNil: true}
	}

	outputsList := findMatches(outputsRegex, chipMatches[4])
	outputs := make(map[string]bool)
	for _, key := range outputsList {
		outputs[key] = true
		wires[key] = FheType{IsNil: true}
	}

	thisChipIO := chipIO{inputs: inputs, outputs: outputs}

	partsLines := findMatches(partLineRegex, chipMatches[5])

	subChips := make([]subChip, 0)
	for _, partStr := range partsLines {
		partName := findMatches(partNameRegex, partStr)[0]
		chipWires := findMatches(partChipWireRegex, partStr)
		thisWires := findMatches(partThisWireRegex, partStr)

		var subChipEvaluatable evaluatableChip
		var subChipIO chipIO
		subChipStateful := false
		subChipName := "__builtin__"
		if partName == "register_false" {
			subChipInstance := newRegisterChip()
			subChipEvaluatable = subChipInstance
			subChipName = "register"
			subChipIO = subChipInstance.chipIO
			subChipStateful = true
		} else if strings.HasPrefix(partName, "program_tape") {
			tapeNum, _ := strconv.Atoi(string(partName[len(partName)-1]))
			subChipInstance := programTapes[tapeNum]
			subChipEvaluatable = subChipInstance
			subChipName = "program_tape"
			subChipIO = subChipInstance.chipIO
			subChipStateful = true
		} else if builtinChip, isBuiltinChip := builtinChips[partName]; isBuiltinChip {
			subChipEvaluatable = builtinChip
			subChipName = builtinChip.name
			subChipIO = builtinChip.chipIO
			subChipStateful = false
		} else {
			b, err := ioutil.ReadFile("./hdl/" + partName + ".hdl")
			if err != nil {
				fmt.Print(err)
			}

			subChipHdl := string(b)
			parsedHdl, e := ParseHdl(subChipHdl)
			if e != nil {
				return hdlChip{}, e
			}
			subChipEvaluatable = parsedHdl
			subChipIO = parsedHdl.chipIO
			subChipStateful = parsedHdl.stateful
			subChipName = parsedHdl.name
		}

		var childParentWiresIn = make(map[string]string)  // inputs
		var childParentWiresOut = make(map[string]string) // outputs
		for idx := range chipWires {
			thisWire := thisWires[idx]
			chipWire := chipWires[idx]
			wires[thisWire] = FheType{IsNil: true}
			if _, isInput := subChipIO.inputs[chipWire]; isInput {
				childParentWiresIn[chipWire] = thisWire
			} else if _, isOutput := subChipIO.outputs[chipWire]; isOutput {
				childParentWiresOut[chipWire] = thisWire
			} else {
				return hdlChip{chipIO: thisChipIO}, fmt.Errorf("Cannot resolve: %s=%s", chipWire, thisWire)
			}
		}

		subChips = append(subChips, subChip{
			chipIO:              subChipIO,
			evaluatableChip:     subChipEvaluatable,
			name:                subChipName,
			stateful:            subChipStateful,
			dead:                false,
			childParentWiresIn:  childParentWiresIn,
			childParentWiresOut: childParentWiresOut})
	}

	return hdlChip{chipIO: thisChipIO, name: name, stateful: stateful, subChips: subChips, wires: wires}, nil
}

func findMatches(regexObj *regexp.Regexp, str string) []string {
	matchesDirty := regexObj.FindAllStringSubmatch(str, -1)
	matchesClean := make([]string, 0)
	for _, arr := range matchesDirty {
		matchesClean = append(matchesClean, arr[len(arr)-1])
	}

	return matchesClean
}
