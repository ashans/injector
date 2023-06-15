package container

import (
	"reflect"
)

const (
	tagInject = `inject`
)

type dependencyTree struct {
	typeToBind map[reflect.Type]map[string][]*bind
}

type bindToTypeValue struct {
	targetType reflect.Type
	qualifier  string
}

func newDependencyTree(c container) dependencyTree {
	depMap := buildDependencyBindMap(c)

	return dependencyTree{typeToBind: buildDependencyTypeMap(depMap)}
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

func buildDependencyTypeMap(mapping map[*bind][]bindToTypeValue) (typeMap map[reflect.Type]map[string][]*bind) {
	typeMap = make(map[reflect.Type]map[string][]*bind)

	for _, bindToTypeValues := range mapping {
		for _, value := range bindToTypeValues {
			if _, ok := typeMap[value.targetType]; !ok {
				typeMap[value.targetType] = make(map[string][]*bind)
			}
			typeMap[value.targetType][value.qualifier] = make([]*bind, 0)
		}
	}

	for bind := range mapping {
		for targetType, qualifierMap := range typeMap {
			if reflect.TypeOf(bind).ConvertibleTo(targetType) {
				if _, hasEmptyQualifier := qualifierMap[``]; hasEmptyQualifier {
					qualifierMap[``] = append(qualifierMap[``], bind)
				}
				if bind.qualifier == `` {
					continue
				}
				if _, hasSpecificQualifier := qualifierMap[bind.qualifier]; hasSpecificQualifier {
					qualifierMap[bind.qualifier] = append(qualifierMap[bind.qualifier], bind)
				}
			}
		}
	}

	return typeMap
}
