# go-xplat

A go cross platform library with abstractions for standard lib global functions.

## using

```bash
go get github.com/patrickhuber/go-xplat
```

### os

```go
import(
  "github.com/patrickhuber/go-xplat/os"
)
func main(){
  o := os.New()
  fmt.Println(o.Executable())
}
```

### env

```go
import(
  "github.com/patrickhuber/go-xplat/env"
)
func main(){
  e := env.NewOS()
  e.Set("MY_ENV_VAR", "test")
  fmt.Println(e.Get("MY_ENV_VAR"))
}
```

```
test
```

### xstd

```go
import(
  "github.com/patrickhuber/go-xplat/xstd"
)  
func main(){
  c := xstd.NewOS()
  fmt.Fprintln(c.Out(), "hello world")
}
```

```
hello world
```

### filepath

```go
func main(){
  fp, err := fs.Parse("/some/path/to/parse")
  if err != nil{
    fmt.Fprintf("%w", err)
    os.Exit(1)
  }
  for _, seg := range fp.Segments{
    fmt.Println(seg)
  }
}
```

```
some
path
to
parse
```