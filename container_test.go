package injector

import (
	"fmt"
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

func TestWithQualifier(t *testing.T) {
	a := new(StructA)
	b1 := new(StructB2)
	b2 := new(StructB2)
	ab := new(StructAB)
	c := new(StructC)

	container := New()
	err := container.NamedBind(a, `a`)
	if err != nil {
		t.Fatal(err)
	}
	err = container.NamedBind(b1, `b1`)
	if err != nil {
		t.Fatal(err)
	}
	err = container.NamedBind(b2, `b2`)
	if err != nil {
		t.Fatal(err)
	}
	err = container.NamedBind(c, `c`)
	if err != nil {
		t.Fatal(err)
	}
	err = container.NamedBind(ab, `ab`)
	if err != nil {
		t.Fatal(err)
	}

	err = container.ResolveDependencyTree()
	if err != nil {
		if resolveErr, ok := err.(DependencyResolveError); ok {
			for _, value := range resolveErr.notFound {
				fmt.Println(`no `, value.qualifier, value.targetType)
			}
			for value := range resolveErr.multipleFound {
				fmt.Println(`multiple `, value.qualifier, value.targetType)
			}
		}

		t.Fatal(err)
	}

	assert.Equal(t, a, c.a1)
	assert.Equal(t, ab, c.a2)
	assert.Equal(t, a, c.a3)
	assert.Equal(t, b1, c.b1)
	assert.Equal(t, b2, c.b2)
	assert.Equal(t, ab, c.ab)
	assert.Equal(t, c, c.c1)
	assert.Equal(t, c, c.c2)
	assert.Equal(t, c, c.c3)
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

type StructAB struct{}

func (s *StructAB) MethodA() {}
func (s *StructAB) MethodB() {}

type InterfaceC interface {
	MethodC()
}

type StructA struct{}

func (s *StructA) MethodA() {}

type StructB struct {
	a InterfaceA `inject:""`
}

func (s *StructB) MethodB() {}

type StructB2 struct{}

func (s *StructB2) MethodB() {}

type StructC struct {
	a1 InterfaceA  `inject:"a"`
	a2 InterfaceA  `inject:"ab"`
	a3 *StructA    `inject:"a"`
	b1 InterfaceB  `inject:"b1"`
	b2 *StructB2   `inject:"b2"`
	ab InterfaceAB `inject:"ab"`
	c1 InterfaceC  `inject:"c"`
	c2 InterfaceC  `inject:""`
	c3 *StructC    `inject:""`
}

func (s *StructC) MethodC() {}
