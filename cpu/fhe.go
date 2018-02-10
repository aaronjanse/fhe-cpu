package cpu

/*
#cgo LDFLAGS:  -ltfhe-spqlios-fma
#cgo CFLAGS: -I/usr/local/include
#include <stdlib.h>
#include <tfhe/tfhe.h>
*/
import "C"

var Params *C.TFheGateBootstrappingParameterSet
var CloudKey *C.TFheGateBootstrappingCloudKeySet

func LoadCloudKeyFromFile(path string) {
	cloudKeyFile := C.fopen(C.CString(path), C.CString("rb"))
	CloudKey = C.new_tfheGateBootstrappingCloudKeySet_fromFile(cloudKeyFile)
	Params = CloudKey.params
}

// func EncryptMap(plaintext map[string]int) map[string]FheType {
// 	output := make(map[string]FheType)
// 	for key, value := range plaintext {
// 		output[key] = encrypt(value)
// 	}
// 	return output
// }

// func encrypt(bit int) FheType {
// 	ciphertext := C.new_gate_bootstrapping_ciphertext(params)
// 	C.bootsSymEncrypt(ciphertext, C.int32_t(bit), key)
// 	return FheType{isNil: false, value: ciphertext}
// }

// func decryptMap(ciphertext map[string]FheType) map[string]int {
// 	output := make(map[string]int)
// 	for key, value := range ciphertext {
// 		output[key] = decrypt(value)
// 	}
// 	return output
// }

// func decrypt(ciphertext FheType) int {
// 	return int(C.bootsSymDecrypt(ciphertext.value, key))
// }

func LoadInstructionsFromFile(instructionCount int, path string) [3][]FheType {
	var instructions [3][]FheType

	cloudDataFile := C.fopen(C.CString(path), C.CString("rb"))
	for i := 0; i < instructionCount; i++ {
		for tapeIdx := 0; tapeIdx < 3; tapeIdx++ {
			value := NewCiphertext()
			C.import_gate_bootstrapping_ciphertext_fromFile(cloudDataFile, value, Params)
			instructions[tapeIdx] = append(instructions[tapeIdx], FheType{false, value})
		}
	}
	C.fclose(cloudDataFile)

	return instructions
}

func SaveCiphertext(path string, data FheType) {
	ciphertextFile := C.fopen(C.CString(path), C.CString("ab"))
	C.export_gate_bootstrapping_ciphertext_toFile(ciphertextFile, data.Value, Params)
	C.fclose(ciphertextFile)
}

func NewCiphertext() *C.LweSample {
	return C.new_gate_bootstrapping_ciphertext(Params)
}

func NewConstant(value int) FheType {
	result := NewCiphertext()
	C.bootsCONSTANT(result, C.int32_t(value), CloudKey)
	return FheType{false, result}
}

var fheNand = decorator(func(result, a, b *C.LweSample) {
	C.bootsNAND(result, a, b, CloudKey)
})

var fheAnd = decorator(func(result, a, b *C.LweSample) {
	C.bootsAND(result, a, b, CloudKey)
})

var fheOr = decorator(func(result, a, b *C.LweSample) {
	C.bootsOR(result, a, b, CloudKey)
})

var fheXor = decorator(func(result, a, b *C.LweSample) {
	C.bootsXOR(result, a, b, CloudKey)
})

// sel?b:a
func fheMux(sel, b, a FheType) FheType {
	result := NewCiphertext()
	C.bootsMUX(result, sel.Value, b.Value, a.Value, CloudKey)
	return FheType{false, result}
}

func fheNot(a FheType) FheType {
	result := NewCiphertext()
	C.bootsNOT(result, a.Value, CloudKey)
	return FheType{false, result}
}

func fheCopy(a FheType) FheType {
	result := NewCiphertext()
	C.bootsCOPY(result, a.Value, CloudKey)
	return FheType{false, result}
}

func decorator(f func(result, a, b *C.LweSample)) func(a, b FheType) FheType {
	return func(a, b FheType) FheType {
		result := NewCiphertext()
		f(result, a.Value, b.Value)
		return FheType{false, result}
	}
}

type FheType struct {
	IsNil bool
	Value *C.LweSample
}
