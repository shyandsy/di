package di

import (
	"reflect"

	"github.com/pkg/errors"
)

func (c *container) Resolve(object interface{}) error {
	objectTp := reflect.TypeOf(object)
	objectVal := reflect.ValueOf(object)

	if object == nil {
		return errors.New("object cannot be nil")
	}
	if objectTp.Kind() != reflect.Pointer || objectTp.Elem().Kind() != reflect.Struct {
		return errors.New("object must be a pointer of struct")
	}

	objectVal = objectVal.Elem()

	element := objectTp.Elem()
	for i := 0; i < element.NumField(); i++ {
		tpField := element.Field(i)
		if _, ok := tpField.Tag.Lookup("inject"); !ok {
			continue
		}

		if tpField.Type.Kind() != reflect.Interface &&
			(tpField.Type.Kind() != reflect.Pointer || tpField.Type.Elem().Kind() != reflect.Struct) {
			return errors.New("inject field only for interface or *struct")
		}

		field := objectVal.FieldByName(tpField.Name)
		if !field.IsValid() || !field.CanSet() {
			return errors.New("inject field must CanSet")
		}

		s, err := c.parseStruct(tpField.Type)
		if err != nil {
			return err
		}

		dep, ok := c.singletonStore.Load(s.FullType())
		if ok {
			field.Set(reflect.ValueOf(dep))
			continue
		}

		if tpField.Type.Kind() == reflect.Pointer && tpField.Type.Elem().Kind() == reflect.Struct {
			ptrValue := reflect.New(tpField.Type.Elem()).Interface()
			if err = c.Resolve(ptrValue); err != nil {
				return err
			}
			field.Set(reflect.ValueOf(ptrValue))
			continue
		}

		return errors.New("dependency not found: " + s.FullType())
	}

	return nil
}
