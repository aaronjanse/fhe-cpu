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

	"github.com/spf13/cobra"
)

// genKeypairCmd represents the genKeypair command
var genKeysCmd = &cobra.Command{
	Use:   "gen-keys",
	Short: "Generate a fhe keypair",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating fhe parameters...")

		minimumLambda := C.int32_t(110)
		params := C.new_default_gate_bootstrapping_parameters(minimumLambda)

		fmt.Println("Generating fhe keypair...")

		// seed := []int{314, 1592, 657}
		// buf := new(bytes.Buffer)
		// binary.Write(buf, binary.LittleEndian, seed)
		// b := buf.Bytes()
		// C.tfhe_random_generator_setSeed((*C.uint32_t)(unsafe.Pointer(&b[0])), C.int32_t(3)) // buf.Len()

		secretKey := C.new_random_gate_bootstrapping_secret_keyset(params)
		cloudKey := &secretKey.cloud

		secretKeyFile := C.fopen(C.CString(secretKeyPath), C.CString("wb"))
		C.export_tfheGateBootstrappingSecretKeySet_toFile(secretKeyFile, secretKey)
		C.fclose(secretKeyFile)

		cloudKeyFile := C.fopen(C.CString(cloudKeyPath), C.CString("wb"))
		C.export_tfheGateBootstrappingCloudKeySet_toFile(cloudKeyFile, cloudKey)
		C.fclose(cloudKeyFile)
	},
}

var secretKeyPath string
var cloudKeyPath string

func init() {
	rootCmd.AddCommand(genKeysCmd)

	genKeysCmd.Flags().StringVarP(&secretKeyPath, "secret", "", "./secret.key", "path where secret key will be stored")
	genKeysCmd.Flags().StringVarP(&cloudKeyPath, "cloud", "", "./cloud.key", "path where cloud key will be stored")
}
