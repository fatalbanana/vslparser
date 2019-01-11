# vslparser

Simple and efficient parser for varnishlog output.

## DESCRIPTION

`vslparser` is a Go module which provides basic functionality required for
further processing of `varnishlog` output.

`varnishlog` output is made up of entries, such as `Request` or `BeReq`. The
provided `Parse` function parses a single such entry and returns a convenient
`Entry` struct. Little processing is required to obtain the `Entry`, only line
splitting and basic sanity checks are performed.

The `Entry` provides a range of convenience methods which further digest the
entry. This way, only the fields which actually need to be parsed are ever
processed. The parsing process is therefore about as efficient as it gets, and
the API is easy to use at the same time.


## EXAMPLE

Compile this minimal example to start "varnishlog" and use `vslparser` to parse
its standard output.

```go
package main

import (
	"fmt"
	"github.com/Showmax/vslparser"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("varnishlog") // Add "-b" to only process back-end requests.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	for {
		entry, err := vslparser.Parse(scanner)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(entry)
	}
}
```

## CONTRIBUTING

Contributions are welcome. Open a PR and we'll get to you soon.

## LICENSE

Apache 2.0.