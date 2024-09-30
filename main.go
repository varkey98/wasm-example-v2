package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"log"
	"regexp"
	"sync"
	"time"
	"wasm-example-v2/arena"
	"wasm-example-v2/book"
)

//go:embed body.json
var description []byte

const GoRoutines = 501

func main() {
	start := time.Now()
	MultipleGoRoutinesTest()
	fmt.Printf("Elapsed Time for WASM based invocation: %v\n", time.Since(start))

	start = time.Now()
	MultipleGoRoutinesTestWithoutWasm()
	fmt.Printf("Elapsed Time for Normal Invocation: %v\n", time.Since(start))
}

func MultipleGoRoutinesTest() {
	processor, err := Initialise()
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(GoRoutines)
	for i := 0; i < GoRoutines; i++ {
		go func() {
			ctx := context.Background()

			ctx = arena.Initialize(ctx)
			defer arena.Reset(ctx)

			module := processor.GetModule()
			defer processor.ResetModule(module)

			fn := module.ExportedFunction("ProcessRegex")
			if fn == nil {
				fmt.Printf("[%d] no function named Process", i)
				return
			}

			req := book.Book{
				Name:        uuid.New().String(),
				Description: string(description),
			}

			res, err := fn.Call(ctx, arena.Store(ctx, &req))
			if err != nil {
				log.Panicf("[%d] failed to invoke Process: %v", i, err)
			}

			if _, ok := arena.Load(ctx, res[0]).(*book.Book); ok {
				//fmt.Printf("Processed Value: %s\n", obj.Name)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func MultipleGoRoutinesTestWithoutWasm() {

	wg := sync.WaitGroup{}
	wg.Add(GoRoutines)
	for i := 0; i < GoRoutines; i++ {
		go func() {
			regex := regexp.MustCompile(`.*traceable.*`)
			req := book.Book{
				Name:        uuid.New().String(),
				Description: string(description),
			}

			if regex.MatchString(req.Name) {
				req.Name = req.Name + ": processed"
			}

			//fmt.Printf("Processed Value: %s\n", req.Name)
			wg.Done()
		}()
	}

	wg.Wait()
}
