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
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"

	"github.com/aaronduino/fhe-cpu/cpu"
	logging "github.com/op/go-logging"

	"github.com/spf13/cobra"
)

var log = logging.MustGetLogger("")

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the cpu",
	Long:  `The instructions should be encrypted via the sibling 'encrypt' command`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := "profile_results.txt"
		if profile != "" {
			f, err := os.Create(profile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}

		cloudKeyPath, _ = filepath.Abs(cloudKeyPath)
		cloudDataPath, _ = filepath.Abs(cloudDataPath)
		// load cloud key and params
		cpu.LoadCloudKeyFromFile(cloudKeyPath)

		instructionCount, _ := strconv.Atoi(readLastLine(cloudDataPath))
		instructions := cpu.LoadInstructionsFromFile(instructionCount, cloudDataPath)
		cpu.InitTapes(instructions)

		b, err := ioutil.ReadFile("hdl/computer.hdl")
		if err != nil {
			fmt.Print(err)
		}

		computerHdl := string(b)
		computer, e := cpu.ParseHdl(computerHdl)

		if e != nil {
			fmt.Println(e)
		}

		os.Remove(outputCiphertextPath)

		log.Debug("Running program...")
		runtimeInputs := map[string]cpu.FheType{"input_bit": cpu.NewConstant(0)}
		for {
			outputsCiphertext := computer.Evaluate(runtimeInputs, "#")
			cpu.SaveCiphertext(outputCiphertextPath, outputsCiphertext["code0"])
			cpu.SaveCiphertext(outputCiphertextPath, outputsCiphertext["code1"])
			log.Info(outputsCiphertext)

			outputCount++

			if maxOutputCount != -1 && outputCount >= maxOutputCount {
				saveOutputCount()
				os.Exit(0)
			}
		}
	},
}

var verbose bool
var outputCiphertextPath string
var outputCount int
var maxOutputCount int

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&cloudKeyPath, "cloud", "k", "./cloud.key", "path to cloud key")
	runCmd.Flags().StringVarP(&cloudDataPath, "instructions", "i", "./cloud.data", "path to instruction data")
	runCmd.Flags().StringVarP(&outputCiphertextPath, "output", "o", "./cloud.out", "path to output data")
	runCmd.Flags().IntVarP(&maxOutputCount, "cycles", "c", -1, "Max number of cycles to run CPU")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Be verbose")

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan

		saveOutputCount()
		os.Exit(0)
	}()
}

func saveOutputCount() {
	// append the output count to the file
	f, err := os.OpenFile(outputCiphertextPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("\n" + strconv.Itoa(outputCount) + "\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func readLastLine(fname string) string {
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	buf := make([]byte, 32)
	n, err := file.ReadAt(buf, fi.Size()-int64(len(buf)))
	if err != nil {
		fmt.Println(err)
	}

	return strings.Split(string(buf[:n]), "\n")[1]
}
