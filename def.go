package injector

type Container interface {
	Bind(instance interface{}) error
	NamedBind(instance interface{}, qualifier string) error
	ResolveDependencyTree() error
	RunModules() error
}

type Runnable interface {
	Run() error
}
