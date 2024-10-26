package main

//
//import (
//	"context"
//	"fmt"
//	"io"
//	"log"
//	"net/http"
//
//	"github.com/tetratelabs/wazero"
//	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
//)
//
//func main() {
//
//	ctx := context.Background()
//	cc, err := wazero.NewCompilationCacheWithDir("/Users/varkeychanjacob/Projects/wasm-example-v2/wazerocache")
//	if err != nil {
//		log.Fatal(err)
//	}
//	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCompilationCache(cc))
//	defer r.Close(ctx)
//
//	wasi_snapshot_preview1.MustInstantiate(ctx, r)
//
//	// Works first time no problem
//	c1 := getCompiledModule(ctx, r)
//	doAdd(ctx, r, c1)
//
//	// Fails second time because first module is closed - why?
//	c2 := getCompiledModule(ctx, r)
//	c1.Close(ctx) // I want to close c1, but keep c2 ready for use
//	doAdd(ctx, r, c2)
//}
//
//func getCompiledModule(ctx context.Context, r wazero.Runtime) wazero.CompiledModule {
//	resp, _ := http.Get("https://github.com/tetratelabs/wazero/raw/main/examples/basic/testdata/add.wasm")
//	wasm, _ := io.ReadAll(resp.Body)
//	cm, _ := r.CompileModule(ctx, wasm)
//	return cm
//}
//
//func doAdd(ctx context.Context, r wazero.Runtime, cm wazero.CompiledModule) {
//	cfg := wazero.NewModuleConfig()
//	mod, err := r.InstantiateModule(ctx, cm, cfg)
//	if err != nil {
//		// On second call, panics here with "source module must be compiled before instantiation"
//		// Why? The module is indeed compiled.
//		panic(err)
//	}
//
//	x := uint64(10)
//	y := uint64(20)
//	res, _ := mod.ExportedFunction("add").Call(ctx, x, y)
//	fmt.Printf("%d + %d = %d\n", x, y, res[0])
//}
