# go-xplat

A go cross platform library with abstractions for standard lib global functions.

## using

```bash
go get github.com/patrickhuber/go-xplat
```

### xos

```go
import(
  "github.com/patrickhuber/go-xplat/xos"
)
func main(){
  o := xos.New()
  fmt.Println(o.Executable())
}
```

### xenv

```go
import(
  "github.com/patrickhuber/go-xplat/xenv"
)
func main(){
  e := xenv.NewOS()
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

### xfilepath

```go
func main(){
  fp, err := xfs.Parse("/some/path/to/parse")
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