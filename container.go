package container

import (
	"errors"
	"reflect"
)

type bind struct {
	instance         interface{}
	instanceResolver interface{}
	singleton        bool
	qualifier        string
}

type container struct {
	binds        []*bind
	treeResolved bool
}

func New() Container {
	return &container{
		binds: make([]*bind, 0),
	}
}

func (c *container) Bind(instance interface{}) error {
	return c.bind(instance, "")
}

func (c *container) NamedBind(instance interface{}, qualifier string) error {
	if qualifier == "" {
		return errors.New(`qualifier name should not be empty`)
	}

	return c.bind(instance, qualifier)
}

func (c *container) ResolveTree() error {
	c.treeResolved = true
	// TODO : this must inject dependencies to components in container
	return nil
}

func (c *container) bind(instance interface{}, qualifier string) error {
	if c.treeResolved {
		return errors.New(`cannot bind new instances after dependency tree is resolved`)
	}

	instanceType := reflect.TypeOf(instance)
	if instanceType == nil {
		return errors.New(`bind instance is null`)
	}

	if instanceType.Kind() != reflect.Ptr {
		return errors.New(`bind instance must be pointer`)
	}

	c.binds = append(c.binds, &bind{
		instance:  instance,
		singleton: true,
		qualifier: qualifier,
	})

	return nil
}
