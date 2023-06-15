package injector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithoutQualifier(t *testing.T) {
	a := &StructA{}
	b := &StructB{}

	c := New()
	err := c.Bind(a)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Bind(b)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ResolveDependencyTree()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, a, b.a)
}

type InterfaceA interface {
	MethodA()
}

type InterfaceB interface {
	MethodB()
}

type InterfaceAB interface {
	MethodA()
	MethodB()
}

type StructA struct{}

func (s *StructA) MethodA() {}

type StructB struct {
	a InterfaceA `inject:""`
}

func (s *StructB) MethodB() {}
