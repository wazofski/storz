package schema

type NestedWorld interface {
    Description() string
    SetDescription(v string)
    Counter() int
    SetCounter(v int)
    Alive() bool
    SetAlive(v bool)
}

type NestedWorldWrapper struct {
    description string
    counter int
    alive bool
}

func (o NestedWorldWrapper) Description() string {
    return o.description
}

func (o NestedWorldWrapper) SetDescription(v string) {
    o.description = v
}

func (o NestedWorldWrapper) Counter() int {
    return o.counter
}

func (o NestedWorldWrapper) SetCounter(v int) {
    o.counter = v
}

func (o NestedWorldWrapper) Alive() bool {
    return o.alive
}

func (o NestedWorldWrapper) SetAlive(v bool) {
    o.alive = v
}


type WorldSpec interface {
    Name() string
    SetName(v string)
    Nested() NestedWorld
    SetNested(v NestedWorld)
}

type WorldSpecWrapper struct {
    name string
    nested NestedWorld
}

func (o WorldSpecWrapper) Name() string {
    return o.name
}

func (o WorldSpecWrapper) SetName(v string) {
    o.name = v
}

func (o WorldSpecWrapper) Nested() NestedWorld {
    return o.nested
}

func (o WorldSpecWrapper) SetNested(v NestedWorld) {
    o.nested = v
}


type WorldStatus interface {
    Description() string
    SetDescription(v string)
    List() []string
    SetList(v []string)
}

type WorldStatusWrapper struct {
    description string
    list []string
}

func (o WorldStatusWrapper) Description() string {
    return o.description
}

func (o WorldStatusWrapper) SetDescription(v string) {
    o.description = v
}

func (o WorldStatusWrapper) List() []string {
    return o.list
}

func (o WorldStatusWrapper) SetList(v []string) {
    o.list = v
}



type World interface {
    Spec() WorldSpec
    Status() WorldStatus
}

type WorldWrapper struct {
    spec WorldSpec
    status WorldStatus
}

func (o WorldWrapper) Spec() WorldSpec {
    return o.spec
}

func (o WorldWrapper) SetSpec(v WorldSpec) {
    o.spec = v
}

func (o WorldWrapper) Status() WorldStatus {
    return o.status
}

func (o WorldWrapper) SetStatus(v WorldStatus) {
    o.status = v
}



