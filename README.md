# go-abstractions

A go abstractions library for difficult to test functions in stdlib

## using

```bash
go get github.com/patrickhuber/go-abstractions
```

### os

```go
import(
  "github.com/patrickhuber/go-abstractions/os"
)
func main(){
  o := os.New()
  fmt.Println(o.Executable())
}
```

### env

```go
import(
  "github.com/patrickhuber/go-abstractions/env"
)
func main(){
  e := env.NewOS()
  fmt.Println(e.Get("MY_ENV_VAR"))
}
```

### console

```go
import(
  "github.com/patrickhuber/go-abstractions/console"
)  
func main(){
  c := console.NewOS()
  fmt.FPrintln(c.Out(), "hello world")
}
```