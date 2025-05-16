package jobs_manager

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var FuncDefinitions = make(map[string]*FuncDefinition, 0)

type FuncDefinition struct {
	Handler      interface{}
	HandlerValue reflect.Value
	HandlerName  string
}

func Register(regFunc interface{}) {
	v := reflect.ValueOf(regFunc)
	typ := v.Type()
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected handler to be func, but got: %T", regFunc))
	}

	funcPath := strings.Split(runtime.FuncForPC(v.Pointer()).Name(), ".")

	def := &FuncDefinition{
		Handler:      regFunc,
		HandlerValue: v,
		HandlerName:  funcPath[1],
	}

	FuncDefinitions[def.HandlerName] = def
}

func FindFuncDefinition(funcName string) *FuncDefinition {
	for _, funcDef := range FuncDefinitions {
		if funcDef.HandlerName == funcName {
			return funcDef
		}
	}
	return nil
}
