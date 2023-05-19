package interop

import (
	"context"
	"errors"

	"github.com/analog-substance/tengo/v2"
)

type CompiledFuncRunner struct {
	ctx      context.Context
	compiled *tengo.Compiled
	fn       *tengo.CompiledFunction
}

func NewCompiledFuncRunner(fn *tengo.CompiledFunction, compiled *tengo.Compiled, ctx context.Context) CompiledFuncRunner {
	return CompiledFuncRunner{
		ctx:      ctx,
		compiled: compiled,
		fn:       fn,
	}
}

func (r *CompiledFuncRunner) Run(args ...tengo.Object) (tengo.Object, error) {
	vm := tengo.NewVM(r.compiled.Bytecode(), r.compiled.Globals(), -1)
	ch := make(chan tengo.Object, 1)

	go func() {
		obj, err := vm.RunCompiled(r.fn, args...)
		if err != nil {
			ch <- GoErrToTErr(err)
			return
		}

		ch <- obj
	}()

	var obj tengo.Object
	var err error
	select {
	case <-r.ctx.Done():
		vm.Abort()
		err = r.ctx.Err()
	case obj = <-ch:
	}

	if err != nil {
		return nil, err
	}

	errObj, ok := obj.(*tengo.Error)
	if ok {
		return nil, errors.New(errObj.String())
	}

	return obj, nil
}
