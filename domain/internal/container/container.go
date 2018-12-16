package container

import (
	"fmt"
	"reflect"
)

type singletonContainer struct {
	values map[string]interface{}
}

var instance *singletonContainer

func init() {
	instance = &singletonContainer{
		make(map[string]interface{}),
	}
}

// Submit an instance to container
func Submit(val interface{}) error {
	var key string
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		key = reflect.Indirect(reflect.ValueOf(val)).Type().String()
	} else {
		key = reflect.ValueOf(val).Type().String()
	}

	if instance.values[key] != nil {
		return fmt.Errorf("already submitted such type of %s", key)
	}

	instance.values[key] = val
	return nil
}

// Override an instance to container
func Override(val interface{}) {
	var key string
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		key = reflect.Indirect(reflect.ValueOf(val)).Type().String()
	} else {
		key = reflect.ValueOf(val).Type().String()
	}

	instance.values[key] = val
	return
}

// Get an instance from container
func Get(ptr interface{}) error {
	val := reflect.ValueOf(ptr)
	key := reflect.Indirect(val).Type().String()
	component := instance.values[key]
	if component == nil {
		if reflect.Indirect(val).Type().Kind() != reflect.Interface {
			return fmt.Errorf("component not found. such type of %s", key)
		}
		for _, component := range instance.values {
			value := reflect.ValueOf(component)
			elm := reflect.ValueOf(ptr).Elem()
			if value.Type().Implements(elm.Type()) {
				elm.Set(value)
				return nil
			}
		}
		return fmt.Errorf("component not found. such type of %s", key)
	}

	elm := reflect.ValueOf(ptr).Elem()
	if reflect.TypeOf(component).Kind() == reflect.Ptr {
		elm.Set(reflect.Indirect(reflect.ValueOf(component)))
	} else {
		elm.Set(reflect.ValueOf(component))
	}
	return nil
}
