package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
	"github.com/analog-substance/tengomod/types"
)

type ExecCmd struct {
	types.PropObject
	Value *exec.Cmd
}

// TypeName should return the name of the type.
func (c *ExecCmd) TypeName() string {
	return "exec-cmd"
}

// String should return a string representation of the type's value.
func (c *ExecCmd) String() string {
	return fmt.Sprintf("<exec-cmd>: %s", strings.Join(c.Value.Args, ", "))
}

// IsFalsy should return true if the value of the type should be considered
// as falsy.
func (c *ExecCmd) IsFalsy() bool {
	return c.Value == nil
}

// CanIterate should return whether the Object can be Iterated.
func (c *ExecCmd) CanIterate() bool {
	return false
}

func (c *ExecCmd) run(args ...tengo.Object) (tengo.Object, error) {
	err := RunCmdWithSigHandler(c.Value)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

func (c *ExecCmd) setStdin(args ...tengo.Object) (tengo.Object, error) {
	file, err := interop.TStrToGoStr(args[0], "file")
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}
	c.Value.Stdin = f

	return nil, nil
}

func makeExecCmd(cmd *exec.Cmd) *ExecCmd {
	execCmd := &ExecCmd{
		Value: cmd,
	}

	objectMap := map[string]tengo.Object{
		"run": &tengo.UserFunction{
			Name:  "run",
			Value: execCmd.run,
		},
		"set_stdin": &tengo.UserFunction{
			Name:  "set_stdin",
			Value: interop.NewCallable(execCmd.setStdin, interop.WithExactArgs(1)),
		},
	}

	execCmd.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: make(map[string]types.Property),
	}

	return execCmd
}
