package di

import (
	"reflect"

	"github.com/pkg/errors"
)

func (c *container) Resolve(object interface{}) error {
	objectTp := reflect.TypeOf(object)
	objectVal := reflect.ValueOf(object).Elem()

	if object == nil {
		return errors.New("object cannot be nil")
	}
	if objectTp.Kind() != reflect.Pointer || objectTp.Elem().Kind() != reflect.Struct {
		return errors.New("object must be a pointer of struct")
	}

	element := objectTp.Elem()
	for i := 0; i < element.NumField(); i++ {
		tpField := element.Field(i)
		if _, ok := tpField.Tag.Lookup("inject"); !ok {
			continue
		}
		if tpField.Type.Kind() != reflect.Interface &&
			(tpField.Type.Kind() != reflect.Pointer || tpField.Type.Elem().Kind() != reflect.Struct) {
			continue
		}
		if !tpField.IsExported() {
			continue
		}

		field := objectVal.FieldByName(tpField.Name)
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		s, err := c.parseStruct(tpField.Type)
		if err != nil {
			return err
		}

		dep, ok := c.singletonStore.Load(s.FullType())
		if !ok {
			continue
		}

		field.Set(reflect.ValueOf(dep))
	}

	return nil
}
