# goconfig

# usage example

```go
package main

import (
	"fmt"
	"your-package/goconfig"
)

type Config struct {
	Hello string `yaml:"hello"`
}

func main() {
	var cfg *Config
	goconfig.LoadConfig(&cfg, "", "config.yaml")
}
```

# yaml config example
```yaml
hello: world
```