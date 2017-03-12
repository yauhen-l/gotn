# gotn
Get test name at position in `*_test.go` file

Name `gotn` stands for `go test name`

This tool was written to execute particular go test from editor (e.g. Emacs)

##Requirements
- golang 1.7 (not tested on other versions)

##Installation
`go get github.com/yauhen-l/gotn`
Then make sure `gotn` executable in your `PATH`

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
>go test -v -run ^`gotn -f gotn_test.go -p 550`$
=== RUN   Test
=== RUN   Test/second_level
--- PASS: Test (0.00s)
    --- PASS: Test/second_level (0.00s)
PASS
ok      github.com/yauhenl/gotn 0.002s
```

##Emacs
Add `gotn.el` to `load-path`
```
(require 'gotn)
(add-hook 'go-mode-hook (lambda ()
                        (local-set-key (kbd "C-c t") 'gotn-run-test))
                        (local-set-key (kbd "C-c C-t") 'gotn-run-test-package)))
```

For customizations see group `gotn`.

```
(customer-set-variables
  '(go-test-case-command "gb test -v -test.run")
  '(go-test-package-command "gb test -v")
  '(go-test-package-test-fallback t))
```

##Restrictions
- Supports only standard Go testing framework: https://golang.org/pkg/testing/
- Does not support `"testing"` package import aliases
