package container

type Container interface {
	Bind(instance interface{}) error
	NamedBind(instance interface{}, qualifier string) error
	ResolveTree() error
}
