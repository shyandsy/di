package di

import (
	"reflect"

	"github.com/pkg/errors"
)

// Find
// Parameters:
//
//	object: struct pointer which
func (c *container) Find(object interface{}) error {
	objectTp := reflect.TypeOf(object)
	objectVal := reflect.ValueOf(object)

	if object == nil {
		return errors.New("object cannot be nil")
	}
	if objectTp.Kind() != reflect.Pointer {
		return errors.New("object must be a pointer of struct or a pointer of interface")
	}
	if objectTp.Elem().Kind() != reflect.Struct && objectTp.Elem().Kind() != reflect.Interface {
		return errors.New("object must be a pointer of struct or a pointer of interface")
	}

	s, err := c.parseStruct(objectTp)
	if err != nil {
		return err
	}

	target, ok := c.singletonStore.Load(s.FullType())
	if !ok {
		return errors.New("object not found in container")
	}
	targetVal := reflect.ValueOf(target)

	if !objectVal.Elem().CanSet() {
		return errors.New("object cannot be set")
	}

	if objectTp.Elem().Kind() == reflect.Interface {
		objectVal.Elem().Set(targetVal)
	} else if objectTp.Elem().Kind() == reflect.Struct {
		objectVal.Elem().Set(targetVal.Elem())
	}

	return nil
}
