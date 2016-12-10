# gotn
Get test name at position in _test.go file

Name `gotn` stands for `go test name`

This tool was written to execute particular go test from editor (e.g. emacs ofr vim)

##Requirements
- golang 1.7 (not tested on other versions)

##Usage
Run `Test/top_level` from `gotn_test.go`
```
>go test -v -run ^`gotn -f gotn_test.go -p 350`$
=== RUN   Test
=== RUN   Test/top_level
--- PASS: Test (0.00s)
    --- PASS: Test/top_level (0.00s)
PASS
ok      github.com/yauhenl/gotn 0.002s
```

Run `Test/second_level` from `gotn_test.go`
```
go test -v -run ^`gotn -f gotn_test.go -p 550`$
=== RUN   Test
=== RUN   Test/second_level
--- PASS: Test (0.00s)
    --- PASS: Test/second_level (0.00s)
PASS
ok      github.com/yauhenl/gotn 0.002s
```

##TODO
- Integrate with emacs
