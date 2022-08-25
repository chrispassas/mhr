[![Go Reference](https://pkg.go.dev/badge/github.com/chrispassas/mhr.svg)](https://pkg.go.dev/github.com/chrispassas/mhr) [![Go Report Card](https://goreportcard.com/badge/github.com/chrispassas/mhr)](https://goreportcard.com/report/github.com/chrispassas/mhr)


# Malware Hash Registry (MHR) 
https://team-cymru.com/community-services/mhr/

This Go Module supports bulk queries to the Team Cymru MHR API.

## Description
Using this module you can submit md5, sha1 or sha256 binary hashes to MHR. The API will return the Anti-virus hit rate and time.

## Go Doc
https://pkg.go.dev/github.com/chrispassas/mhr



## Example

### Parse Whole File
```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/chrispassas/mhr"
)

func main() {
	var hashes = []string{
		"7697561ccbbdd1661c25c86762117613",
		"d48a85139dde1eb00ee7460e80f42c35",
		"8a62d103168974fba9c61edab336038c",
	}
	var results []mhr.Result
	var err error
	if results, err = mhr.Search(context.Background(), hashes); err != nil {
		log.Fatalf("mhr.Search() error:%v", err)
	}

	for _, r := range results {
		log.Printf("hash:%s hit:%d time:%s nodata:%t", r.Hash, r.HitRate, r.Timestamp.Format(time.RFC3339), r.NoData)
	}
}

```

