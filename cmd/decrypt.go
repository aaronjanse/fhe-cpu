package cmd

/*
#cgo LDFLAGS:  -ltfhe-spqlios-fma
#cgo CFLAGS: -I/usr/local/include
#include <stdlib.h>
#include <tfhe/tfhe.h>
*/
import "C"

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		secretKeyFile := C.fopen(C.CString(secretKeyPath), C.CString("rb"))
		secretKey := C.new_tfheGateBootstrappingSecretKeySet_fromFile(secretKeyFile)
		params := secretKey.cloud.params

		outputCount, _ := strconv.Atoi(readLastLine(outputCiphertextPath))
		outputsFile := C.fopen(C.CString(outputCiphertextPath), C.CString("rb"))
		for i := 0; i < outputCount; i++ {
			values := make([]int, 0)
			for idx := 0; idx < 2; idx++ {
				ciphertext := C.new_gate_bootstrapping_ciphertext(params)
				C.import_gate_bootstrapping_ciphertext_fromFile(outputsFile, ciphertext, params)
				value := int(C.bootsSymDecrypt(ciphertext, secretKey))
				values = append(values, value)
			}
			fmt.Println(values)
		}
		C.fclose(outputsFile)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringVarP(&outputCiphertextPath, "input", "i", "./cloud.out", "path to output from cpu")
	decryptCmd.Flags().StringVarP(&secretKeyPath, "secret", "", "./secret.key", "path to secret key")
}
