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
	if ok {
		targetVal := reflect.ValueOf(target)

		if !objectVal.Elem().CanSet() {
			return errors.New("object cannot be set")
		}

		if targetVal.Kind() == reflect.Func {
			objs, err := c.Invoke(target)
			if err != nil {
				return err
			}
			targetVal = objs[0]
		}

		objectVal.Elem().Set(targetVal.Elem())
	} else if objectTp.Elem().Kind() == reflect.Struct {
		ptrValue := reflect.New(objectTp.Elem()).Interface()
		if err = c.Resolve(ptrValue); err != nil {
			return err
		}
		objectVal.Elem().Set(reflect.ValueOf(ptrValue).Elem())
	} else {
		return errors.New("dependency not found: " + s.FullType())
	}

	return nil
}
