package injector

import (
	"errors"
	env "github.com/caarlos0/env/v8"
	"reflect"
)

type bind struct {
	instance  interface{}
	singleton bool
	qualifier string
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

func (c *container) ResolveDependencyTree() error {
	if c.treeResolved {
		return errors.New(`cannot resolve dependencies of already resolved container`)
	}

	err := c.resolveEnvironmentVariables()
	if err != nil {
		return err
	}

	tree, err := buildTree(c)
	if err != nil {
		return err
	}

	err = tree.injectDependencies(c)
	if err != nil {
		return err
	}
	c.treeResolved = true

	return nil
}

func (c *container) RunModules() error {
	if !c.treeResolved {
		return errors.New(`cannot run modules if dependencies are not resolved`)
	}
	for _, b := range c.binds {
		if runnable, isRunnable := b.instance.(Runnable); isRunnable {
			err := runnable.Run()
			if err != nil {
				return err
			}
		}
	}

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

func (c *container) resolveEnvironmentVariables() error {
	for _, b := range c.binds {
		instance := b.instance
		if err := env.Parse(instance); err != nil {
			return err
		}
	}

	return nil
}
