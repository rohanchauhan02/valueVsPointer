package main

type BigStruct struct {
	Buf [1 << 18]byte
}

var obj BigStruct

func main() {
	PassByValue(obj)
	PassByPointer(&obj)
}

func PassByValue(obj BigStruct) {}

func PassByPointer(obj *BigStruct) {}
