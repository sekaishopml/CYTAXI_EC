package containers

type Container struct {
	services map[string]any
}

func New() *Container {
	return &Container{services: make(map[string]any)}
}

func (c *Container) Register(name string, service any) {
	c.services[name] = service
}

func (c *Container) Resolve(name string) any {
	return c.services[name]
}
