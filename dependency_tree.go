package injector

import (
	"fmt"
	"reflect"
	"unsafe"
)

const (
	tagInject = `inject`
)

type dependencyTree struct {
	typeToBind map[reflect.Type]map[string]*bind
}

func (d *dependencyTree) PrintMatch() {
	for t, qMap := range d.typeToBind {
		for q, b := range qMap {
			fmt.Printf("type: %v with qualifier: %s matched with bind type: %v with qualifier: %s\n", t, q,
				reflect.TypeOf(b.instance), b.qualifier)
		}
	}
}

type bindToTypeValue struct {
	targetType reflect.Type
	qualifier  string
}

func (b bindToTypeValue) String() string {
	if b.qualifier == `` {
		return fmt.Sprintf(`type: %v`, b.targetType)
	}
	return fmt.Sprintf(`type: %v | qualifier: %s`, b.targetType, b.qualifier)
}

func (b bindToTypeValue) DebugString() string {
	return b.String()
}

func buildTree(c *container) (*dependencyTree, error) {
	depMap := buildDependencyBindMap(c)
	s := fmt.Sprintf(`%v`, depMap[c.binds[3]][0].targetType)
	println(s)
	typeMap := buildDependencyTypeMap(c, depMap)
	singleInstanceMap, err := findSingleInstances(typeMap)
	if err != nil {
		return nil, err
	}

	return &dependencyTree{typeToBind: singleInstanceMap}, nil
}

func buildDependencyBindMap(c *container) map[*bind][]bindToTypeValue {
	bindToType := make(map[*bind][]bindToTypeValue)

	for _, b := range c.binds {
		receiverType := reflect.TypeOf(b.instance)
		if receiverType.Kind() != reflect.Pointer || receiverType.Elem().Kind() != reflect.Struct {
			bindToType[b] = make([]bindToTypeValue, 0)
			continue
		}

		fields := reflect.ValueOf(b.instance).Elem()

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

func buildDependencyTypeMap(c *container, mapping map[*bind][]bindToTypeValue) (typeMap map[reflect.Type]map[string][]*bind) {
	typeMap = make(map[reflect.Type]map[string][]*bind)

	for _, bindToTypeValues := range mapping {
		for _, value := range bindToTypeValues {
			if _, ok := typeMap[value.targetType]; !ok {
				typeMap[value.targetType] = make(map[string][]*bind)
			}
			typeMap[value.targetType][value.qualifier] = make([]*bind, 0)
		}
	}

	for _, bind := range c.binds {
		for targetType, qualifierMap := range typeMap {
			if reflect.TypeOf(bind.instance).ConvertibleTo(targetType) {
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

func (t *dependencyTree) injectDependencies(c *container) error {
	for _, b := range c.binds {
		receiverType := reflect.TypeOf(b.instance)
		if receiverType.Kind() != reflect.Pointer || receiverType.Elem().Kind() != reflect.Struct {
			continue
		}

		fields := reflect.ValueOf(b.instance).Elem()

		for i := 0; i < fields.NumField(); i++ {
			field := fields.Field(i)

			if qualifier, tagExists := fields.Type().Field(i).Tag.Lookup(tagInject); tagExists {
				targetType := field.Type()
				bindVal, hasBind := t.typeToBind[targetType][qualifier]
				if !hasBind || bindVal == nil {
					return DependencyInjectError{bindToTypeValue{targetType: targetType, qualifier: qualifier}}
				}
				instance := bindVal.instance

				pointer := reflect.NewAt(targetType, unsafe.Pointer(field.UnsafeAddr())).Elem()
				pointer.Set(reflect.ValueOf(instance))
			}
		}
	}

	return nil
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
				if _, hasTargetTypeMap := singleInstances[targetType]; !hasTargetTypeMap {
					singleInstances[targetType] = make(map[string]*bind)
				}
				singleInstances[targetType][qualifier] = binds[0]
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
