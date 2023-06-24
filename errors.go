package injector

import (
	"fmt"
)

type DependencyResolveError struct {
	notFound      []bindToTypeValue
	multipleFound map[bindToTypeValue][]bindToTypeValue
}

func (d DependencyResolveError) Error() string {
	return fmt.Sprintf(`dependency resolve error - not found: %d, multiple found: %d`, len(d.notFound), len(d.multipleFound))
}

type DependencyInjectError struct {
	notFound bindToTypeValue
}

func (d DependencyInjectError) Error() string {
	typeString := fmt.Sprintf(`bind value not found for type: %v`, d.notFound.targetType)
	if d.notFound.qualifier == `` {
		return typeString
	}

	return fmt.Sprintf(`%s and qualifier: %s`, typeString, d.notFound.qualifier)
}
