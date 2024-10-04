# Understanding Pass by Value vs. Pass by Reference in Go

As developers, we often encounter the concepts of pass by value and pass by reference when working with functions. These two methods of argument passing can significantly affect how our code behaves, especially when it comes to handling data in memory.

In this blog post, we will explore the differences between pass by value and pass by reference in Go, and how they impact the behavior of our code.

```go:main.go
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
```

## Pass by Value

Pass by value creates a copy of the value and passes it to the function. Any changes to the copy inside the function don't touch the original value outside.

Think of it like making a photocopy of a document. You can scribble on the copy, but the original stays clean.

In Go, simple types like integers, floats, and booleans are passed by value.

## Pass by Reference

Pass by reference hands over the memory address of the value to the function. Changes made inside the function directly modify the original value.

Imagine sharing a live Google Doc. When someone edits it, everyone sees the changes instantly.

In Go, complex types like slices, maps, and pointers are passed by reference.

## Benchmarking

We can use benchmarks to see the performance difference. Let's test with a large data structure (BigStruct roughly 256KB).

```go
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
```

### With Compiler Optimization (data size ~ 256KB)

```
go test -bench=. -count 1
```

```
goos: darwin
goarch: arm64
pkg: github.com/rohanchauhan02/valuevspointer
BenchmarkPassByValue-8          1000000000               0.3544 ns/op
BenchmarkPassByPointer-8        1000000000               0.3589 ns/op
PASS
ok      github.com/rohanchauhan02/valuevspointer        1.011s
```

With optimization, the compiler might be smart enough to avoid unnecessary copying, so the difference is negligible.

Without Compiler Optimization (data size ~ 256KB)

### benchmark test without compliler optimization for data size 1<<18 bytes ~ 256KB

```
go test -bench=. -count 1 -gcflags="-N -l"
```

```
goos: darwin
goarch: arm64
pkg: github.com/rohanchauhan02/valuevspointer
BenchmarkPassByValue-8            874224              1318 ns/op
BenchmarkPassByPointer-8        560874014                2.116 ns/op
PASS
ok      github.com/rohanchauhan02/valuevspointer        3.421s
```

Without optimization, pass by value is significantly slower because it has to copy the entire 256KB structure. Pass by reference only copies a memory address, which is much faster.

![Benchmark results](https://github.com/rohanchauhan02/valueVsPointer/blob/main/doc/img1.png)

## Assembly code without optimization

```
go tool compile -N -S -l main.go
```

```
main.main STEXT size=144 args=0x0 locals=0x80008 funcid=0x0 align=0x0
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     TEXT    main.main(SB), ABIInternal, $524304-0
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     MOVD    16(g), R16
        0x0004 00004 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     PCDATA  $0, $-2
        0x0004 00004 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     SUBS    $524176, RSP, R17
        0x0010 00016 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     BLO     132
        0x0014 00020 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     CMP     R16, R17
        0x0018 00024 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     BLS     132
        0x001c 00028 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     PCDATA  $0, $-1
        0x001c 00028 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     SUB     $524304, RSP, R20
        0x0024 00036 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     STP     (R29, R30), -8(R20)
        0x0028 00040 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     PCDATA  $0, $-2
        0x0028 00040 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     MOVD    R20, RSP
        0x002c 00044 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     PCDATA  $0, $-1
        0x002c 00044 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     SUB     $8, RSP, R29
        0x0030 00048 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        0x0030 00048 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        0x0030 00048 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:8)     MOVD    $main.obj-262144(SP), R16
        0x0038 00056 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:8)     MOVD    $main.obj-16(SP), R0
        0x0040 00064 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:8)     STP.P   (ZR, ZR), 16(R16)
        0x0044 00068 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:8)     CMP     R0, R16
        0x0048 00072 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:8)     BLE     64
        0x004c 00076 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     MOVD    $8(RSP), R16
        0x0050 00080 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     MOVD    $262136(RSP), R0
        0x0058 00088 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     STP.P   (ZR, ZR), 16(R16)
        0x005c 00092 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     CMP     R0, R16
        0x0060 00096 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     BLE     88
        0x0064 00100 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     PCDATA  $1, $0
        0x0064 00100 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     CALL    main.PassByValue(SB)
        0x0068 00104 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:10)    MOVD    $main.obj-262144(SP), R0
        0x0070 00112 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:10)    CALL    main.PassByPointer(SB)
        0x0074 00116 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    LDP     -8(RSP), (R29, R30)
        0x0078 00120 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    ADD     $524304, RSP
        0x0080 00128 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    RET     (R30)
        0x0084 00132 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    NOP
        0x0084 00132 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     PCDATA  $1, $-1
        0x0084 00132 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     PCDATA  $0, $-2
        0x0084 00132 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     MOVD    R30, R3
        0x0088 00136 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     CALL    runtime.morestack_noctxt(SB)
        0x008c 00140 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     PCDATA  $0, $-1
        0x008c 00140 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:7)     JMP     0
        0x0000 90 0b 40 f9 1b f2 9f d2 fb 00 a0 f2 f1 63 3b eb  ..@..........c;.
        0x0010 a3 03 00 54 3f 02 10 eb 69 03 00 54 f4 43 00 d1  ...T?...i..T.C..
        0x0020 94 02 42 d1 9d fa 3f a9 9f 02 00 91 fd 23 00 d1  ..B...?......#..
        0x0030 f0 03 41 91 10 22 00 91 e0 ff 41 91 00 e0 3f 91  ..A.."....A...?.
        0x0040 1f 7e 81 a8 1f 02 00 eb cd ff ff 54 f0 23 00 91  .~.........T.#..
        0x0050 e0 ff 40 91 00 e0 3f 91 1f 7e 81 a8 1f 02 00 eb  ..@...?..~......
        0x0060 cd ff ff 54 00 00 00 94 e0 03 41 91 00 20 00 91  ...T......A.. ..
        0x0070 00 00 00 94 fd fb 7f a9 ff 43 00 91 ff 03 42 91  .........C....B.
        0x0080 c0 03 5f d6 e3 03 1e aa 00 00 00 94 dd ff ff 17  .._.............
        rel 100+4 t=R_CALLARM64 main.PassByValue+0
        rel 112+4 t=R_CALLARM64 main.PassByPointer+0
        rel 136+4 t=R_CALLARM64 runtime.morestack_noctxt+0
main.PassByValue STEXT size=16 args=0x40000 locals=0x0 funcid=0x0 align=0x0 leaf
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:13)    TEXT    main.PassByValue(SB), LEAF|NOFRAME|ABIInternal, $0-262144
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:13)    FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:13)    FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:13)    FUNCDATA        $5, main.PassByValue.arginfo1(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:13)    RET     (R30)
        0x0000 c0 03 5f d6 00 00 00 00 00 00 00 00 00 00 00 00  .._.............
main.PassByPointer STEXT size=16 args=0x8 locals=0x0 funcid=0x0 align=0x0 leaf
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    TEXT    main.PassByPointer(SB), LEAF|NOFRAME|ABIInternal, $0-8
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    FUNCDATA        $0, gclocals·wgcWObbY2HYnK2SU/U22lA==(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    FUNCDATA        $1, gclocals·J5F+7Qw7O7ve2QcWC7DpeQ==(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    FUNCDATA        $5, main.PassByPointer.arginfo1(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    MOVD    R0, main.obj(FP)
        0x0004 00004 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    RET     (R30)
        0x0000 e0 07 00 f9 c0 03 5f d6 00 00 00 00 00 00 00 00  ......_.........
go:cuinfo.producer.<unlinkable> SDWARFCUINFO dupok size=0
        0x0000 2d 4e 20 2d 6c 20 72 65 67 61 62 69              -N -l regabi
go:cuinfo.packagename.main SDWARFCUINFO dupok size=0
        0x0000 6d 61 69 6e                                      main
main..inittask SNOPTRDATA size=8
        0x0000 00 00 00 00 00 00 00 00                          ........
type:.eqfunc262144 SRODATA dupok size=16
        0x0000 00 00 00 00 00 00 00 00 00 00 04 00 00 00 00 00  ................
        rel 0+8 t=R_ADDR runtime.memequal_varlen+0
runtime.memequal64·f SRODATA dupok size=8
        0x0000 00 00 00 00 00 00 00 00                          ........
        rel 0+8 t=R_ADDR runtime.memequal64+0
runtime.gcbits.0100000000000000 SRODATA dupok size=8
        0x0000 01 00 00 00 00 00 00 00                          ........
type:.namedata._main.BigStruct. SRODATA dupok size=17
        0x0000 01 0f 2a 6d 61 69 6e 2e 42 69 67 53 74 72 75 63  .._main.BigStruc
        0x0010 74                                               t
type:_main.BigStruct SRODATA size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 ec 94 7d 9d 08 08 08 36 00 00 00 00 00 00 00 00  ..}....6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=R_ADDR runtime.memequal64·f+0
        rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
        rel 40+4 t=R_ADDROFF type:.namedata._main.BigStruct.+0
        rel 48+8 t=R_ADDR type:main.BigStruct+0
runtime.gcbits. SRODATA dupok size=0
type:.namedata._[262144]uint8- SRODATA dupok size=16
        0x0000 00 0e 2a 5b 32 36 32 31 34 34 5d 75 69 6e 74 38  .._[262144]uint8
type:_[262144]uint8 SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 3d 92 c2 7b 08 08 08 36 00 00 00 00 00 00 00 00  =..{...6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=R_ADDR runtime.memequal64·f+0
        rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
        rel 40+4 t=R_ADDROFF type:.namedata._[262144]uint8-+0
        rel 48+8 t=R_ADDR type:[262144]uint8+0
type:[262144]uint8 SRODATA dupok size=72
        0x0000 00 00 04 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0010 2c 22 b2 54 0a 01 01 11 00 00 00 00 00 00 00 00  ,".T............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 00 00 04 00 00 00 00 00                          ........
        rel 24+8 t=R_ADDR type:.eqfunc262144+0
        rel 32+8 t=R_ADDR runtime.gcbits.+0
        rel 40+4 t=R_ADDROFF type:.namedata._[262144]uint8-+0
        rel 44+4 t=RelocType(-32763) type:_[262144]uint8+0
        rel 48+8 t=R_ADDR type:uint8+0
        rel 56+8 t=R_ADDR type:[]uint8+0
type:.namedata.Buf. SRODATA dupok size=5
        0x0000 01 03 42 75 66                                   ..Buf
type:.importpath.main. SRODATA dupok size=6
        0x0000 00 04 6d 61 69 6e                                ..main
type:main.BigStruct SRODATA size=120
        0x0000 00 00 04 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0010 99 7c be 3f 0f 01 01 19 00 00 00 00 00 00 00 00  .|.?............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 01 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
        0x0050 00 00 00 00 00 00 00 00 28 00 00 00 00 00 00 00  ........(.......
        0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0070 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=R_ADDR type:.eqfunc262144+0
        rel 32+8 t=R_ADDR runtime.gcbits.+0
        rel 40+4 t=R_ADDROFF type:.namedata.*main.BigStruct.+0
        rel 44+4 t=R_ADDROFF type:*main.BigStruct+0
        rel 56+8 t=R_ADDR type:main.BigStruct+96
        rel 80+4 t=R_ADDROFF type:.importpath.main.+0
        rel 96+8 t=R_ADDR type:.namedata.Buf.+0
        rel 104+8 t=R_ADDR type:[262144]uint8+0
gclocals·g2BeySu+wFnoycgXfElmcg== SRODATA dupok size=8
        0x0000 01 00 00 00 00 00 00 00                          ........
main.PassByValue.arginfo1 SRODATA static dupok size=26
        0x0000 fe fe 00 01 01 01 02 01 03 01 04 01 05 01 06 01  ................
        0x0010 07 01 08 01 09 01 fc fd fd ff                    ..........
gclocals·wgcWObbY2HYnK2SU/U22lA== SRODATA dupok size=10
        0x0000 02 00 00 00 01 00 00 00 01 00                    ..........
gclocals·J5F+7Qw7O7ve2QcWC7DpeQ== SRODATA dupok size=8
        0x0000 02 00 00 00 00 00 00 00                          ........
main.PassByPointer.arginfo1 SRODATA static dupok size=3
        0x0000 00 08 ff                                         ...
➜  valueVsPointer go tool compile -N -S -l main.go
main.main STEXT size=128 args=0x0 locals=0x40008 funcid=0x0 align=0x0
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     TEXT    main.main(SB), ABIInternal, $262160-0
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     MOVD    16(g), R16
        0x0004 00004 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     PCDATA  $0, $-2
        0x0004 00004 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     SUBS    $262032, RSP, R17
        0x0010 00016 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     BLO     116
        0x0014 00020 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     CMP     R16, R17
        0x0018 00024 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     BLS     116
        0x001c 00028 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     PCDATA  $0, $-1
        0x001c 00028 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     SUB     $262160, RSP, R20
        0x0024 00036 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     STP     (R29, R30), -8(R20)
        0x0028 00040 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     PCDATA  $0, $-2
        0x0028 00040 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     MOVD    R20, RSP
        0x002c 00044 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     PCDATA  $0, $-1
        0x002c 00044 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     SUB     $8, RSP, R29
        0x0030 00048 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        0x0030 00048 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        0x0030 00048 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    MOVD    $8(RSP), R17
        0x0034 00052 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    MOVD    $main.obj(SB), R16
        0x003c 00060 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    MOVD    $main.obj+262128(SB), R0
        0x0044 00068 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    PCDATA  $0, $-2
        0x0044 00068 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    LDP.P   16(R16), (R25, R27)
        0x0048 00072 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    STP.P   (R25, R27), 16(R17)
        0x004c 00076 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    PCDATA  $0, $-1
        0x004c 00076 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    CMP     R0, R16
        0x0050 00080 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    BLE     68
        0x0054 00084 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    PCDATA  $1, $0
        0x0054 00084 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:11)    CALL    main.PassByValue(SB)
        0x0058 00088 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:12)    MOVD    $main.obj(SB), R0
        0x0060 00096 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:12)    CALL    main.PassByPointer(SB)
        0x0064 00100 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:13)    LDP     -8(RSP), (R29, R30)
        0x0068 00104 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:13)    ADD     $262160, RSP
        0x0070 00112 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:13)    RET     (R30)
        0x0074 00116 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:13)    NOP
        0x0074 00116 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     PCDATA  $1, $-1
        0x0074 00116 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     PCDATA  $0, $-2
        0x0074 00116 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     MOVD    R30, R3
        0x0078 00120 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     CALL    runtime.morestack_noctxt(SB)
        0x007c 00124 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     PCDATA  $0, $-1
        0x007c 00124 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:9)     JMP     0
        0x0000 90 0b 40 f9 1b f2 9f d2 7b 00 a0 f2 f1 63 3b eb  ..@.....{....c;.
        0x0010 23 03 00 54 3f 02 10 eb e9 02 00 54 f4 43 00 d1  #..T?......T.C..
        0x0020 94 02 41 d1 9d fa 3f a9 9f 02 00 91 fd 23 00 d1  ..A...?......#..
        0x0030 f1 23 00 91 10 00 00 90 10 02 00 91 00 00 00 90  .#..............
        0x0040 00 00 00 91 19 6e c1 a8 39 6e 81 a8 1f 02 00 eb  .....n..9n......
        0x0050 ad ff ff 54 00 00 00 94 00 00 00 90 00 00 00 91  ...T............
        0x0060 00 00 00 94 fd fb 7f a9 ff 43 00 91 ff 03 41 91  .........C....A.
        0x0070 c0 03 5f d6 e3 03 1e aa 00 00 00 94 e1 ff ff 17  .._.............
        rel 52+8 t=R_ADDRARM64 main.obj+0
        rel 60+8 t=R_ADDRARM64 main.obj+262128
        rel 84+4 t=R_CALLARM64 main.PassByValue+0
        rel 88+8 t=R_ADDRARM64 main.obj+0
        rel 96+4 t=R_CALLARM64 main.PassByPointer+0
        rel 120+4 t=R_CALLARM64 runtime.morestack_noctxt+0
main.PassByValue STEXT size=16 args=0x40000 locals=0x0 funcid=0x0 align=0x0 leaf
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    TEXT    main.PassByValue(SB), LEAF|NOFRAME|ABIInternal, $0-262144
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    FUNCDATA        $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    FUNCDATA        $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    FUNCDATA        $5, main.PassByValue.arginfo1(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:15)    RET     (R30)
        0x0000 c0 03 5f d6 00 00 00 00 00 00 00 00 00 00 00 00  .._.............
main.PassByPointer STEXT size=16 args=0x8 locals=0x0 funcid=0x0 align=0x0 leaf
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:17)    TEXT    main.PassByPointer(SB), LEAF|NOFRAME|ABIInternal, $0-8
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:17)    FUNCDATA        $0, gclocals·wgcWObbY2HYnK2SU/U22lA==(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:17)    FUNCDATA        $1, gclocals·J5F+7Qw7O7ve2QcWC7DpeQ==(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:17)    FUNCDATA        $5, main.PassByPointer.arginfo1(SB)
        0x0000 00000 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:17)    MOVD    R0, main.obj(FP)
        0x0004 00004 (/Users/rohanchauhan/Learning/golang/valueVsPointer/main.go:17)    RET     (R30)
        0x0000 e0 07 00 f9 c0 03 5f d6 00 00 00 00 00 00 00 00  ......_.........
go:cuinfo.producer.<unlinkable> SDWARFCUINFO dupok size=0
        0x0000 2d 4e 20 2d 6c 20 72 65 67 61 62 69              -N -l regabi
go:cuinfo.packagename.main SDWARFCUINFO dupok size=0
        0x0000 6d 61 69 6e                                      main
main..inittask SNOPTRDATA size=8
        0x0000 00 00 00 00 00 00 00 00                          ........
type:.eqfunc262144 SRODATA dupok size=16
        0x0000 00 00 00 00 00 00 00 00 00 00 04 00 00 00 00 00  ................
        rel 0+8 t=R_ADDR runtime.memequal_varlen+0
runtime.memequal64·f SRODATA dupok size=8
        0x0000 00 00 00 00 00 00 00 00                          ........
        rel 0+8 t=R_ADDR runtime.memequal64+0
runtime.gcbits.0100000000000000 SRODATA dupok size=8
        0x0000 01 00 00 00 00 00 00 00                          ........
type:.namedata._main.BigStruct. SRODATA dupok size=17
        0x0000 01 0f 2a 6d 61 69 6e 2e 42 69 67 53 74 72 75 63  .._main.BigStruc
        0x0010 74                                               t
type:_main.BigStruct SRODATA size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 ec 94 7d 9d 08 08 08 36 00 00 00 00 00 00 00 00  ..}....6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=R_ADDR runtime.memequal64·f+0
        rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
        rel 40+4 t=R_ADDROFF type:.namedata._main.BigStruct.+0
        rel 48+8 t=R_ADDR type:main.BigStruct+0
runtime.gcbits. SRODATA dupok size=0
type:.namedata._[262144]uint8- SRODATA dupok size=16
        0x0000 00 0e 2a 5b 32 36 32 31 34 34 5d 75 69 6e 74 38  .._[262144]uint8
type:_[262144]uint8 SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 3d 92 c2 7b 08 08 08 36 00 00 00 00 00 00 00 00  =..{...6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=R_ADDR runtime.memequal64·f+0
        rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
        rel 40+4 t=R_ADDROFF type:.namedata._[262144]uint8-+0
        rel 48+8 t=R_ADDR type:[262144]uint8+0
type:[262144]uint8 SRODATA dupok size=72
        0x0000 00 00 04 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0010 2c 22 b2 54 0a 01 01 11 00 00 00 00 00 00 00 00  ,".T............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 00 00 04 00 00 00 00 00                          ........
        rel 24+8 t=R_ADDR type:.eqfunc262144+0
        rel 32+8 t=R_ADDR runtime.gcbits.+0
        rel 40+4 t=R_ADDROFF type:.namedata._[262144]uint8-+0
        rel 44+4 t=RelocType(-32763) type:_[262144]uint8+0
        rel 48+8 t=R_ADDR type:uint8+0
        rel 56+8 t=R_ADDR type:[]uint8+0
type:.namedata.Buf. SRODATA dupok size=5
        0x0000 01 03 42 75 66                                   ..Buf
type:.importpath.main. SRODATA dupok size=6
        0x0000 00 04 6d 61 69 6e                                ..main
type:main.BigStruct SRODATA size=120
        0x0000 00 00 04 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0010 99 7c be 3f 0f 01 01 19 00 00 00 00 00 00 00 00  .|.?............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 01 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
        0x0050 00 00 00 00 00 00 00 00 28 00 00 00 00 00 00 00  ........(.......
        0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0070 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=R_ADDR type:.eqfunc262144+0
        rel 32+8 t=R_ADDR runtime.gcbits.+0
        rel 40+4 t=R_ADDROFF type:.namedata.*main.BigStruct.+0
        rel 44+4 t=R_ADDROFF type:*main.BigStruct+0
        rel 56+8 t=R_ADDR type:main.BigStruct+96
        rel 80+4 t=R_ADDROFF type:.importpath.main.+0
        rel 96+8 t=R_ADDR type:.namedata.Buf.+0
        rel 104+8 t=R_ADDR type:[262144]uint8+0
main.obj SNOPTRBSS size=262144
 SDWARFVAR size=25
        0x0000 0a 6d 61 69 6e 2e 6f 62 6a 00 09 03 00 00 00 00  .main.obj.......
        0x0010 00 00 00 00 00 00 00 00 01                       .........
        rel 12+8 t=R_ADDR main.obj+0
        rel 20+4 t=R_DWARFSECREF go:info.main.BigStruct+0
gclocals·g2BeySu+wFnoycgXfElmcg== SRODATA dupok size=8
        0x0000 01 00 00 00 00 00 00 00                          ........
main.PassByValue.arginfo1 SRODATA static dupok size=26
        0x0000 fe fe 00 01 01 01 02 01 03 01 04 01 05 01 06 01  ................
        0x0010 07 01 08 01 09 01 fc fd fd ff                    ..........
gclocals·wgcWObbY2HYnK2SU/U22lA== SRODATA dupok size=10
        0x0000 02 00 00 00 01 00 00 00 01 00                    ..........
gclocals·J5F+7Qw7O7ve2QcWC7DpeQ== SRODATA dupok size=8
        0x0000 02 00 00 00 00 00 00 00                          ........
main.PassByPointer.arginfo1 SRODATA static dupok size=3
        0x0000 00 08 ff                                         ...
```

![Assembly code](https://github.com/rohanchauhan02/valueVsPointer/blob/main/doc/img2.png)

- More Steps in Pass by Value: The extra steps in pass by value are due to the need to copy data, allocate memory, and perform more instructions to handle the value as opposed to a reference.
- Efficiency: Pass by reference is generally more efficient in terms of memory usage and speed, especially for larger data structures, because it avoids the overhead of copying the entire dataset.

### Conclusion

In Go, the choice between pass by value and pass by reference has significant implications for both safety and efficiency:

- Pass by Value: This method ensures that changes made within the function are isolated from the original data, providing a safer approach. However, it may be less efficient when dealing with large data structures since the compiler needs to generate code to copy the value to the function’s argument.

- Pass by Reference: In contrast, passing by reference can enhance efficiency, particularly with large data, as it avoids copying the entire structure. However, it also means that any modifications made in the function will affect the original data.

Understanding these trade-offs is essential for writing efficient and correct Go code. By carefully selecting the method that best aligns with your needs, you can optimize both performance and safety in your applications. Ultimately, the compiler's role in generating code to handle these different passing methods is crucial in achieving the desired outcomes in your code execution.
