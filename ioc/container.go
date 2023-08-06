package ioc

import (
	"fmt"
	"reflect"
)

type DependencyMap map[string]any

type Container struct {
	store DependencyMap
}

func NewContainer() *Container {
	return &Container{
		store: make(DependencyMap),
	}
}

// Register must pass in the non-pointer instantiated struct
// and the factory function returning a pointer to the
// passed in v struct
// TODO: This must work given an interfaces.
// Just pass in factory, deduce the type from the factory return type.
func (c *Container) Register(factory any) error {
	rfFunc := reflect.ValueOf(factory)
	if rfFunc.Kind() != reflect.Func {
		return fmt.Errorf("not a factory function")
	}

	rfType := rfFunc.Type().Out(0)
	for rfType.Kind() == reflect.Ptr {
		rfType = rfType.Elem()
	}
	typeName := rfType.Name()
	pkgName := rfType.PkgPath()

	// fqName is the package path + the struct name
	fqName := fmt.Sprintf("%s/%s", pkgName, typeName)

	c.store[fqName] = factory
	return nil
}

func resolveFactory(c *Container, factory any) (any, error) {
	rfFunc := reflect.ValueOf(factory)
	rfFuncType := rfFunc.Type()

	// Analyze the factory function
	paramsCount := rfFuncType.NumIn()
	params := make([]reflect.Value, 0)
	for i := 0; i < paramsCount; i++ {
		paramType := rfFuncType.In(i)
		paramFqName := fmt.Sprintf("%s/%s", paramType.Elem().PkgPath(), paramType.Elem().Name())
		paramFactory, has := c.store[paramFqName]
		if !has {
			return nil, fmt.Errorf("no factory of %s is stored", paramFqName)
		}
		fmt.Println("Calling " + paramFqName)
		value, err := resolveFactory(c, paramFactory)
		if err != nil {
			return nil, err
		}
		params = append(params, reflect.ValueOf(value))
	}
	outputs := rfFunc.Call(params)
	return outputs[0].Interface(), nil
}

func Create[T any](c *Container, factory any) (*T, error) {
	v, err := resolveFactory(c, factory)
	if err != nil {
		return nil, err
	}
	return v.(*T), nil
}

func Inject[T any](c *Container, v *T) (*T, error) {
	rfType := reflect.TypeOf(v).Elem()
	rfValue := reflect.ValueOf(v).Elem()
	rfNumField := rfType.NumField()

	for i := 0; i < rfNumField; i++ {
		field := rfType.Field(i)
		_, isInjectable := field.Tag.Lookup("inject")
		if !isInjectable {
			continue
		}

		rfFieldType := field.Type
		for rfFieldType.Kind() == reflect.Ptr {
			rfFieldType = rfFieldType.Elem()
		}

		fqName := fmt.Sprintf("%s/%s", rfFieldType.PkgPath(), rfFieldType.Name())
		fmt.Printf("Dependency %d: %s\n", (i + 1), fqName)
		factory, isExist := c.store[fqName]
		if !isExist {
			return nil, fmt.Errorf("not found")
		}
		value, err := resolveFactory(c, factory)
		if err != nil {
			return nil, err
		}
		rfValue.Field(i).Set(reflect.ValueOf(value))
	}

	return v, nil
}
