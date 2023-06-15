package injector

type Container interface {
	Bind(instance interface{}) error
	NamedBind(instance interface{}, qualifier string) error
	ResolveDependencyTree() error
}
