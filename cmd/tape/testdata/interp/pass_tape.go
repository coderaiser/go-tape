package main

import . "tapeapi"

func main() {
	Test("interp: pass", func(t T) {
		t.Equal(1, 1)
		t.End()
	})
}