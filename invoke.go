package di

import (
	"reflect"

	"github.com/pkg/errors"
)

func (c *container) Invoke(f interface{}) ([]reflect.Value, error) {
	tp := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if f == nil {
		return nil, errors.New("f cannot be nil")
	}
	if tp.Kind() != reflect.Func {
		return nil, errors.New("f must be a function")
	}

	var args []reflect.Value
	if tp.NumIn() > 0 {
		for i := 0; i < tp.NumIn(); i++ {
			inType := tp.In(i)

			if (inType.Kind() != reflect.Pointer || inType.Elem().Kind() != reflect.Struct) && inType.Kind() != reflect.Interface {
				return nil, errors.New("parameter must be interface or *struct")
			}

			s, err := c.parseStruct(inType)
			if err != nil {
				return nil, err
			}

			dep, ok := c.singletonStore.Load(s.FullType())
			if ok {
				arg := reflect.ValueOf(dep)
				if inType.Kind() == reflect.Struct || inType.Kind() == reflect.Interface {
					arg = arg.Elem()
				}
				args = append(args, arg)
				continue
			} else {
				if inType.Kind() == reflect.Pointer {
					ptrValue := reflect.New(inType.Elem()).Interface()
					if err = c.Resolve(ptrValue); err != nil {
						return nil, err
					}
					args = append(args, reflect.ValueOf(ptrValue))
					continue
				}
			}

			return nil, errors.New("dependency not found: " + s.FullType())
		}
	}

	values := val.Call(args)

	return values, nil
}
