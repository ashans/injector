package container

import (
	"reflect"
)

const (
	tagInject = `inject`
)

type dependencyTree struct {
	typeToBind map[reflect.Type]map[string]*bind
}

type bindToTypeValue struct {
	targetType reflect.Type
	qualifier  string
}

func newDependencyTree(c container) (tree dependencyTree) {
	_ = buildDependencyBindMap(c)

	return tree
}

func buildDependencyBindMap(c container) map[*bind][]bindToTypeValue {
	bindToType := make(map[*bind][]bindToTypeValue)

	for _, b := range c.binds {
		receiverType := reflect.TypeOf(b.instance)
		if receiverType.Kind() != reflect.Pointer || receiverType.Elem().Kind() != reflect.Struct {
			bindToType[b] = make([]bindToTypeValue, 0)
			continue
		}

		fields := reflect.ValueOf(b).Elem()

		for i := 0; i < fields.NumField(); i++ {
			field := fields.Field(i)

			if qualifier, tagExists := fields.Type().Field(i).Tag.Lookup(tagInject); tagExists {
				bindToType[b] = append(bindToType[b], bindToTypeValue{
					targetType: field.Type(),
					qualifier:  qualifier,
				})
			}
		}
	}

	return bindToType
}
