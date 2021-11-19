package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run command: %v\n", err)
		os.Exit(1)
	}
}

func FindTinygo() string {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(exe)
	return filepath.Join(exePath, "tinygo")
}

func main() {
	tinygo := FindTinygo()
	args := []string{"build", "-gc=leaking", "-target", "eosio", "-wasm-abi=generic", "-scheduler=none", "-opt", "z"}
	if len(os.Args) >= 2 && os.Args[1] == "build" {
		args = append(args, os.Args[2:]...)
		fmt.Println(tinygo, strings.Join(args, " "))
		runCommand(tinygo, args...)
	} else {
		runCommand(tinygo, os.Args[1:]...)
	}
}
