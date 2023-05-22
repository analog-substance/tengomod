package exec

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

var ErrSignaled error = errors.New("process signaled to close")

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"err_signaled": interop.GoErrToTErr(ErrSignaled),
		"run_with_sig_handler": &interop.AdvFunction{
			Name:    "run_with_sig_handler",
			NumArgs: interop.MinArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("cmd-name"), interop.StrSliceArg("args", true)},
			Value:   tengoRunWithSigHandler,
		},
		"cmd": &interop.AdvFunction{
			Name:    "cmd",
			NumArgs: interop.MinArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("cmd-name"), interop.StrSliceArg("args", true)},
			Value:   tengoCmd,
		},
	}
}

func tengoRunWithSigHandler(args map[string]interface{}) (tengo.Object, error) {
	cmdName := args["cmd-name"].(string)
	cmdArgs := args["args"].([]string)

	err := RunWithSigHandler(cmdName, cmdArgs...)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

func tengoCmd(args map[string]interface{}) (tengo.Object, error) {
	cmdName := args["cmd-name"].(string)
	cmdArgs := args["args"].([]string)

	cmd := exec.CommandContext(context.Background(), cmdName, cmdArgs...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return makeExecCmd(cmd), nil
}

func RunWithSigHandler(name string, args ...string) error {
	cmd := exec.CommandContext(context.Background(), name, args...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return RunCmdWithSigHandler(cmd)
}

func RunCmdWithSigHandler(cmd *exec.Cmd) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// relay trapped signals to the spawned process
	signaled := false
	go func() {
		for sig := range sigs {
			signaled = true
			cmd.Process.Signal(sig)
		}
	}()

	defer func() {
		signal.Stop(sigs)
		close(sigs)
	}()

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if !ok {
			return err
		}

		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			if !signaled {
				signaled = status.Signaled()
			}
		}
	}

	if signaled {
		return ErrSignaled
	}

	return nil
}
