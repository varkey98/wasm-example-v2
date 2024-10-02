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

const GoRoutines = 100

func main() {
	MultipleGoRoutinesTest()

	start := time.Now()
	MultipleGoRoutinesTestWithoutWasm()
	fmt.Printf("Elapsed Time for Normal Invocation: %v\n", time.Since(start))
}

func MultipleGoRoutinesTest() {
	initStart := time.Now()
	arr := make([]*Processor, GoRoutines)
	for i, _ := range arr {
		processor, _ := Initialise()
		arr[i] = processor
	}
	fmt.Println("Elapsed Time for WASM init: %v\n", time.Since(initStart))

	wg := sync.WaitGroup{}
	wg.Add(GoRoutines)

	start := time.Now()
	for i := 0; i < GoRoutines; i++ {
		go func() {
			ctx := context.Background()
			ctx = arena.Initialize(ctx)
			module := arr[i].module

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
				//fmt.Printf("Processed Value: %s\n", obj.Description)
			}
			wg.Done()
			arena.Reset(ctx)
		}()
	}
	wg.Wait()
	fmt.Printf("Elapsed Time for WASM based invocation: %v\n", time.Since(start))
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
