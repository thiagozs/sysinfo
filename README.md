# Sysinfo

Proof of concept about system information, like RAM, disk space, memory usage and free tier.

## Usage

```go
package main

import (
    "github.com/thiagozs/sysinfo"
    "fmt"
    )

func main() {
    sys, err := sysinfo.Get()
    if err != nil {
        panic(err.Error())
    }

    fmt.Println(sys.ToString())
}
```