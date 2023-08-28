//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"rye/contrib"
	"rye/env"
	"rye/evaldo"
	"rye/loader"
	"syscall/js"
)

type TagType int
type RjType int
type Series []interface{}

type anyword struct {
	kind RjType
	idx  int
}

type node struct {
	kind  RjType
	value interface{}
}

var CODE []interface{}

//
// main function. Dispatches to appropriate mode function
//

func main1() {
	evaldo.ShowResults = true
	// main_rye_string("print $Hello world$", false, false)
}

//
// main for awk like functionality with rye language
//

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("RyeEvalString", js.FuncOf(RyeEvalString))
	<-c
}

func RyeEvalString(this js.Value, args []js.Value) interface{} {
	sig := false
	subc := true

	code := args[0].String()

	//util.PrintHeader()
	//defer profile.Start(profile.CPUProfile).Stop()

	block, genv := loader.LoadString(code, sig)
	switch val := block.(type) {
	case env.Block:
		es := env.NewProgramState(block.(env.Block).Series, genv)
		evaldo.RegisterBuiltins(es)
		contrib.RegisterBuiltins(es, &evaldo.BuiltinNames)

		if subc {
			ctx := es.Ctx
			es.Ctx = env.NewEnv(ctx)
		}

		evaldo.EvalBlock(es)
		evaldo.MaybeDisplayFailureOrError(es, genv)
		return es.Res.Probe(*es.Idx)
	case env.Error:
		fmt.Println(val.Message)
		return "Error"
	}
	return "Other"
}
