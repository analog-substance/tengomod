package log

import (
	"fmt"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

const (
	msgPrefix  string = "[+]"
	warnPrefix string = "[!]"
	infoPrefix string = "[-]"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"msg": &tengo.UserFunction{
			Name:  "msg",
			Value: logMsg,
		},
		"warn": &tengo.UserFunction{
			Name:  "warn",
			Value: logWarn,
		},
		"info": &tengo.UserFunction{
			Name:  "info",
			Value: logInfo,
		},
	}
}

func logMsg(args ...tengo.Object) (tengo.Object, error) {
	err := log(msgPrefix, args...)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

func logWarn(args ...tengo.Object) (tengo.Object, error) {
	err := log(warnPrefix, args...)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

func logInfo(args ...tengo.Object) (tengo.Object, error) {
	err := log(infoPrefix, args...)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

func log(prefix string, args ...tengo.Object) error {
	logArgs, err := getLogArgs(args...)
	if err != nil {
		return err
	}

	fmt.Printf("%s ", prefix)
	fmt.Print(logArgs...)
	fmt.Println()

	return nil
}

func getLogArgs(args ...tengo.Object) ([]interface{}, error) {
	var logArgs []interface{}
	l := 0
	for _, arg := range args {
		s, _ := tengo.ToString(arg)
		slen := len(s)
		// make sure length does not exceed the limit
		if l+slen > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		l += slen
		logArgs = append(logArgs, s)
	}
	return logArgs, nil
}
