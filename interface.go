package di

type Container interface {
	Provide(object interface{}) error
	ProvideAs(object interface{}, tp interface{}) error
	Find(object interface{}) error
	Resolve(object interface{}) error
}
