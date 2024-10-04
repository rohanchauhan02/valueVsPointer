package main

import "testing"

func BenchmarkPassByValue(t *testing.B) {
	obj := BigStruct{}

	for n := 0; n < t.N; n++ {
		PassByValue(obj)
	}
}

func BenchmarkPassByPointer(t *testing.B) {
	obj := BigStruct{}
	for n := 0; n < t.N; n++ {
		PassByPointer(&obj)
	}
}
