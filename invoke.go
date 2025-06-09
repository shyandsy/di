package di

import (
	"reflect"

	"github.com/pkg/errors"
)

func (c *container) Invoke(f interface{}) error {
	tp := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if f == nil {
		return errors.New("f cannot be nil")
	}
	if tp.Kind() != reflect.Func {
		return errors.New("f must be a function")
	}

	if tp.NumIn() > 0 {
		args := []reflect.Value{}
		for i := 0; i < tp.NumIn(); i++ {
			inType := tp.In(i)
			s, err := c.parseStruct(inType)
			if err != nil {
				return err
			}
			dep, ok := c.singletonStore.Load(s.FullType())
			if !ok {
				return errors.New("dependency not found: " + s.FullType())
			}

			if inType.Kind() == reflect.Struct || inType.Kind() == reflect.Interface {
				args = append(args, reflect.ValueOf(dep).Elem())
			} else if inType.Kind() == reflect.Pointer {
				args = append(args, reflect.ValueOf(dep))
			}

		}
		val.Call(args)
	} else {
		val.Call(nil)
	}

	return nil
}
