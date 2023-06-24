# Injector

> GoLang dependency injection library to be used as a IoC container

## Documentation

### Features

- Container management for **Go** applications
- Seamless dependency injection through tags and type matching
- No manual type casting required

### Examples

#### Container bind and resolve

```go
package main

import "github.com/ashans/injector"

type InterfaceA interface {
	MethodA()
}
type ImplA struct{}

func (*ImplA) MethodA() {}

type StructB struct {
	a InterfaceA `inject:""`
}

func main() {
	c := injector.New()

	c.Bind(new(ImplA))
	c.Bind(new(StructB))

	err := c.ResolveDependencyTree()
	if err != nil {
		panic(err)
	}
}
```

#### Resolve with qualifier for multiple implementations

```go
package main

import "github.com/ashans/injector"

type InterfaceA interface {
	MethodA()
}

type ImplA1 struct{}

func (*ImplA1) MethodA() {}

type ImplA2 struct{}

func (*ImplA2) MethodA() {}

type StructB struct {
	a1 InterfaceA `inject:"a1"`
	a2 InterfaceA `inject:"a2"`
}

func main() {
	c := injector.New()

	c.NamedBind(new(ImplA1), `a1`)
	c.NamedBind(new(ImplA2), `a2`)
	c.Bind(new(StructB))

	err := c.ResolveDependencyTree()
	if err != nil {
		panic(err)
	}
}
```

Refer [test cases](container_test.go) for more examples