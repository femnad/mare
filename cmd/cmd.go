package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

var defaultShell = "sh"

// Input represents configuration for a command to be executed.
type Input struct {
	Command string
	Pwd     string
	Shell   bool
	Sudo    bool
}

// Output represents the response from an executed command.
type Output struct {
	Code   int
	Stdout string
	Stderr string
}

func getCmd(in Input) *exec.Cmd {
	var cmdSlice []string
	if in.Shell {
		cmdSlice = append([]string{defaultShell, "-c"}, in.Command)
	} else {
		cmdSlice = strings.Split(in.Command, " ")
	}
	if in.Sudo {
		cmdSlice = append([]string{"sudo"}, cmdSlice...)
	}

	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	if in.Pwd != "" {
		cmd.Dir = in.Pwd
	}

	return cmd
}

// Run runs a command based on the CmdIn input and returns a CmdOut.
func Run(in Input) (Output, error) {
	cmd := getCmd(in)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return Output{Stdout: stdout.String(), Stderr: stderr.String(), Code: cmd.ProcessState.ExitCode()}, err
}

// RunFormatError runs a command and on any errors returns an error which contains information about the execution.
func RunFormatError(in Input) (Output, error) {
	out, err := Run(in)
	if err == nil {
		return out, nil
	}

	stdout := strings.TrimSpace(out.Stdout)
	stderr := strings.TrimSpace(out.Stderr)
	outStr := fmt.Sprintf("error running command %s", in.Command)

	if stdout != "" {
		outStr += fmt.Sprintf(", stdout: %s", stdout)
	}
	if stderr != "" {
		outStr += fmt.Sprintf(", stderr: %s", stderr)
	}
	outStr += fmt.Sprintf(", error: %v", err)

	return out, fmt.Errorf(outStr)
}

// RunNoOutput runs a command and discards its output.
func RunNoOutput(in Input) error {
	cmd := getCmd(in)
	return cmd.Run()
}
