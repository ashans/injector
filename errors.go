package container

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
