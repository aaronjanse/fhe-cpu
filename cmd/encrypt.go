package cmd

/*
#cgo LDFLAGS:  -ltfhe-spqlios-fma
#cgo CFLAGS: -I/usr/local/include
#include <tfhe/tfhe.h>
#include <tfhe/tfhe_io.h>
#include <stdlib.h>
#include <stdio.h>
*/
import "C"

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var programSymbolParser = map[byte][]int{
	' ': []int{0, 0, 0},
	'+': []int{0, 0, 1},
	'<': []int{0, 1, 0},
	'>': []int{0, 1, 1},
	'[': []int{1, 0, 0},
	']': []int{1, 0, 1},
	';': []int{1, 1, 0},
	',': []int{1, 1, 1}}

var brainToBoolMap = map[byte]string{
	'+': ">[>]+<[+<]>>>>>>>>>[+]<<<<<<<<<",
	'-': ">>>>>>>>>+<<<<<<<<+[>+]<[<]>>>>>>>>>[+]<<<<<<<<<",
	'<': "<<<<<<<<<",
	'>': ">>>>>>>>>",
	',': ">,>,>,>,>,>,>,>,<<<<<<<<",
	'.': ">;>;>;>;>;>;>;>;<<<<<<<<",
	'[': ">>>>>>>>>+<<<<<<<<+[>+]<[<]>>>>>>>>>[+<<<<<<<<[>]+<[+<]",
	']': ">>>>>>>>>+<<<<<<<<+[>+]<[<]>>>>>>>>>]<[+<]",
}

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load fhe key and params
		secretKeyFile := C.fopen(C.CString(secretKeyPath), C.CString("rb"))
		secretKey := C.new_tfheGateBootstrappingSecretKeySet_fromFile(secretKeyFile)
		C.fclose(secretKeyFile)

		params := secretKey.params

		// encrypt util inputs
		// fmt.Printf("%+v\n", *moveTapeInputsCiphertext["clock"].a)

		// read program from stdin
		stdin := os.Stdin
		stdinMessage := make([]byte, 3)
		_, err := io.ReadFull(stdin, stdinMessage)
		programSrc := ""
		for err == nil {
			programSrc += string(stdinMessage)
			_, err = io.ReadFull(stdin, stdinMessage)
		}

		// optionally translate brainf* to boolf*
		if brain {
			tmp := ""
			for _, char := range []byte(programSrc) {
				tmp += brainToBoolMap[char]
			}

			programSrc = tmp
		}

		fmt.Println(len(programSrc))

		// transform source code to valid cpu inputs
		programSrc = strings.TrimSpace(programSrc)
		programSrc = strings.Replace(programSrc, "[", "[ ", -1)
		programSrc = strings.Replace(programSrc, "]", " ]", -1)

		programBytes := []byte(programSrc)
		programBytes = reverse(programBytes)

		// if _, err := tmpfile.Write(content); err != nil {
		// 	log.Fatal(err)
		// }
		// if err := tmpfile.Close(); err != nil {
		// 	log.Fatal(err)
		// }

		cloudDataFile := C.fopen(C.CString(cloudDataPath), C.CString("wb"))
		for _, char := range programBytes {
			machineCode, isValidChar := programSymbolParser[char]
			if !isValidChar {
				continue
			}

			for _, value := range machineCode {
				ciphertext := encrypt(value, secretKey, params)
				C.export_gate_bootstrapping_ciphertext_toFile(cloudDataFile, &ciphertext, params)
			}
		}
		C.fclose(cloudDataFile)

		// append the instruction count to the file
		f, err := os.OpenFile(cloudDataPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte("\n" + strconv.Itoa(len(programBytes)) + "\n")); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	},
}

var cloudDataPath string
var brain bool

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVarP(&secretKeyPath, "secret", "", "./secret.key", "path to secret key")
	encryptCmd.Flags().StringVarP(&cloudDataPath, "data", "", "./cloud.data", "path to output data")
	encryptCmd.Flags().BoolVarP(&brain, "brain", "b", false, "Interpret the input as brainf*ck (as opposed to boolf*ck)")
}

// func encryptMap(plaintext map[string]int, secretKey *C.TFheGateBootstrappingSecretKeySet, params *C.TFheGateBootstrappingParameterSet) map[string]C.LweSample {
// 	output := make(map[string]C.LweSample)
// 	for key, value := range plaintext {
// 		output[key] = encrypt(value, secretKey, params)
// 	}
// 	return output
// }

func encrypt(bit int, secretKey *C.TFheGateBootstrappingSecretKeySet, params *C.TFheGateBootstrappingParameterSet) C.LweSample {
	ciphertext := C.new_gate_bootstrapping_ciphertext(params)
	C.bootsSymEncrypt(ciphertext, C.int32_t(bit), secretKey)
	return *ciphertext
}

func reverse(numbers []byte) []byte {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}
