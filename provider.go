package di

import (
	"reflect"

	"github.com/pkg/errors"
)

// Provide
// Parameters:
//
//	object		- dependency pointer. must be a pointer of struct, and it must be not nil
func (c *container) Provide(object interface{}) error {
	tp := reflect.TypeOf(object)

	if object == nil {
		return errors.New("object cannot be nil")
	}

	if tp.Kind() == reflect.Pointer && tp.Elem().Kind() == reflect.Struct {
		return c.ProvideAs(object, object)
	} else if tp.Kind() == reflect.Func {
		if tp.NumOut() != 1 {
			return errors.New("the function for dependency creation should have exactly 1 return value")
		}
		if tp.Out(0).Kind() != reflect.Interface {
			return errors.New("the function for dependency creation should return an interface")
		}
		ptrValue := reflect.New(tp.Out(0)).Interface()
		return c.ProvideAs(object, ptrValue)
	}

	return errors.New("object must be pointer of struct")

}

// ProvideAs
// Parameters:
//
//	object		- dependency pointer. must be a pointer of struct, and it must be not nil
//	tp			- specific the target type of the dependency, it must be a pointer of interface
func (c *container) ProvideAs(object interface{}, targetType interface{}) error {
	objectTp := reflect.TypeOf(object)
	targetTp := reflect.TypeOf(targetType)

	if object == nil {
		return errors.New("object cannot be nil")
	}
	if targetTp == nil {
		return errors.New("target cannot be nil")
	}
	if reflect.ValueOf(object).IsNil() {
		return errors.New("object value cannot be nil")
	}
	if (objectTp.Kind() != reflect.Pointer || objectTp.Elem().Kind() != reflect.Struct) && objectTp.Kind() != reflect.Func {
		return errors.New("object must be pointer of struct or function")
	}
	if targetTp.Kind() != reflect.Pointer || (targetTp.Elem().Kind() != reflect.Interface && targetTp.Elem().Kind() != reflect.Struct) {
		return errors.New("target must be pointer of interface")
	}

	target := targetTp.Elem()

	if targetTp.Elem().Kind() == reflect.Interface {
		if objectTp.Kind() != reflect.Func && !objectTp.Implements(target) {
			return errors.New("object must implement target interface")
		}
	}

	s, err := c.parseStruct(targetTp)
	if err != nil {
		return err
	}

	fullTypeName := s.FullType()

	c.mu.Lock()
	defer c.mu.Unlock()
	c.singletonStore.Store(fullTypeName, object)

	return nil
}
