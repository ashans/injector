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

func buildTree(c container) (*dependencyTree, error) {
	depMap := buildDependencyBindMap(c)

	typeMap := buildDependencyTypeMap(depMap)
	singleInstanceMap, err := findSingleInstances(typeMap)
	if err != nil {
		return nil, err
	}

	return &dependencyTree{typeToBind: singleInstanceMap}, nil
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

func findSingleInstances(typeMap map[reflect.Type]map[string][]*bind) (map[reflect.Type]map[string]*bind, error) {
	singleInstances := make(map[reflect.Type]map[string]*bind)
	notFound := make([]bindToTypeValue, 0)
	multipleFound := make(map[bindToTypeValue][]bindToTypeValue)

	for targetType, qualifierMap := range typeMap {
		for qualifier, binds := range qualifierMap {
			switch len(binds) {
			case 0:
				notFound = append(notFound, bindToTypeValue{
					targetType: targetType,
					qualifier:  qualifier,
				})
			case 1:
				singleInstances[targetType] = map[string]*bind{qualifier: binds[0]}
			default:
				matches := make([]bindToTypeValue, 0)
				for _, bind := range binds {
					matches = append(matches, bindToTypeValue{
						targetType: reflect.TypeOf(bind),
						qualifier:  bind.qualifier,
					})
				}

				multipleFound[bindToTypeValue{
					targetType: targetType,
					qualifier:  qualifier,
				}] = matches
			}
		}
	}

	if len(notFound) > 0 || len(multipleFound) > 0 {
		return singleInstances, DependencyResolveError{
			notFound:      notFound,
			multipleFound: multipleFound,
		}
	}

	return singleInstances, nil
}
