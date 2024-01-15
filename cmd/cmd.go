package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"

	"github.com/femnad/mare"
)

var (
	defaultShell     = "sh"
	pathEnvKey       = "PATH"
	pathEnvSeparator = ":"
)

// Input represents configuration for a command to be executed.
type Input struct {
	Command         string
	Env             map[string]string
	Pwd             string
	Shell           bool
	ShellCmd        string
	Sudo            bool
	SudoPreserveEnv bool
}

// Output represents the response from an executed command.
type Output struct {
	Code   int
	Stdout string
	Stderr string
}

func mergePaths(curPath, newPath string) string {
	var newPaths []string
	envPaths := strings.Split(curPath, pathEnvSeparator)
	for _, path := range strings.Split(newPath, pathEnvSeparator) {
		if !mare.Contains(envPaths, path) {
			newPaths = append(newPaths, path)
		}
	}

	var allPaths string
	if len(newPaths) > 0 {
		allPaths = fmt.Sprintf("%s%s%s", curPath, pathEnvSeparator, strings.Join(newPaths, pathEnvSeparator))
	} else {
		allPaths = curPath
	}

	return allPaths
}

func getEnv(in Input, curEnv []string) ([]string, error) {
	if len(in.Env) == 0 {
		return curEnv, nil
	}

	var cmdEnv []string
	desiredEnv := in.Env

	for _, envVal := range curEnv {
		keyAndValue := strings.SplitN(envVal, "=", 2)
		if len(keyAndValue) != 2 {
			return nil, fmt.Errorf("unexpected key and value in environment: %s", keyAndValue)
		}
		key := keyAndValue[0]
		value := keyAndValue[1]

		if key == pathEnvKey {
			addPaths, hasPath := desiredEnv[pathEnvKey]
			if !hasPath {
				cmdEnv = append(cmdEnv, envVal)
				continue
			}

			merged := mergePaths(value, addPaths)
			cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", pathEnvKey, merged))

			delete(desiredEnv, pathEnvKey)
			continue
		}

		newValue, ok := desiredEnv[key]
		if ok {
			cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", key, newValue))

			delete(desiredEnv, key)
			continue
		}

		cmdEnv = append(cmdEnv, envVal)
	}

	for k, v := range desiredEnv {
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, v))
	}

	return cmdEnv, nil
}

func getCmd(in Input) (*exec.Cmd, error) {
	var cmdSlice []string
	var err error

	shell := in.ShellCmd
	if shell == "" {
		shell = defaultShell
	}

	if in.Shell {
		cmdSlice = append([]string{shell, "-c"}, in.Command)
	} else {
		cmdSlice, err = shlex.Split(in.Command)
		if err != nil {
			return &exec.Cmd{}, err
		}
	}
	if in.Sudo {
		sudoSlice := []string{"sudo"}
		if in.SudoPreserveEnv {
			sudoSlice = append(sudoSlice, "-E")
		}
		cmdSlice = append(sudoSlice, cmdSlice...)
	}

	newPath, ok := in.Env[pathEnvKey]
	if ok {
		var updatedPath string

		path := os.Getenv(pathEnvKey)
		if path == "" {
			updatedPath = path
		} else {
			updatedPath = mergePaths(path, newPath)
		}

		err = os.Setenv(pathEnvKey, updatedPath)
		if err != nil {
			return &exec.Cmd{}, fmt.Errorf("error setting path to %s: %v", updatedPath, err)
		}
	}
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)

	cmdEnv, err := getEnv(in, os.Environ())
	if err != nil {
		return &exec.Cmd{}, err
	}
	cmd.Env = cmdEnv

	if in.Pwd != "" {
		cmd.Dir = in.Pwd
	}

	return cmd, nil
}

// Run runs a command based on the CmdIn input and returns a CmdOut.
func Run(in Input) (Output, error) {
	cmd, err := getCmd(in)
	if err != nil {
		return Output{}, fmt.Errorf("error parsing command: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	return Output{Stdout: stdout.String(), Stderr: stderr.String(), Code: cmd.ProcessState.ExitCode()}, err
}

// RunFmtErr runs a command and on any errors returns an error which contains stdout and stderr output if any.
func RunFmtErr(in Input) (Output, error) {
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

// RunErrOnly runs a command and on any errors returns an error and ignores the output.
func RunErrOnly(in Input) error {
	_, err := RunFmtErr(in)
	return err
}

// RunNoOut runs a command and discards its output.
func RunNoOut(in Input) error {
	cmd, err := getCmd(in)
	if err != nil {
		return err
	}

	return cmd.Run()
}
