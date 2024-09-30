package main

import "testing"

func BenchmarkWasm(b *testing.B) {
	b.Run("bench", func(b *testing.B) {
		MultipleGoRoutinesTestWithoutWasm()
	})
}
