package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"wasm-example-v2/arena"
	"wasm-example-v2/book"
)

const GoRoutines = 501

func main() {
	processor, err := Initialise()
	if err != nil {
		panic(err)
	}

	for i := 0; i < GoRoutines; i++ {
		ctx := context.Background()

		ctx = arena.Initialize(ctx)
		defer arena.Reset(ctx)

		fn := processor.module.ExportedFunction("Process")
		if fn == nil {
			fmt.Printf("[%d] no function named Process", i)
			return
		}

		req := book.Book{
			Name: uuid.New().String(),
		}

		res, err := fn.Call(ctx, arena.Store(ctx, &req))
		if err != nil {
			log.Panicf("[%d] failed to invoke Process: %v", i, err)
		}

		if obj, ok := arena.Load(ctx, res[0]).(*book.Book); ok {
			fmt.Printf("Processed Value: %s\n", obj.Name)
		}

	}
}
