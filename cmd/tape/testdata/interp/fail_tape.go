package main

import . "tapeapi"

func main() {
	Test("interp: fail", func(t T) {
		t.Equal(1, 2)
		t.End()
	})
}