package pool

import (
	"fmt"
	"testing"

	"github.com/panenming/go-im/libs/commonspool"
	"github.com/stretchr/testify/assert"
)

type MyPoolObject struct {
}

func TestExample(t *testing.T) {
	p := pool.NewObjectPoolWithDefaultConfig(pool.NewPooledObjectFactorySimple(
		func() (interface{}, error) {
			return &MyPoolObject{}, nil
		}))
	obj, _ := p.BorrowObject()
	p.ReturnObject(obj)
}

type MyObjectFactory struct {
}

func (f *MyObjectFactory) MakeObject() (*pool.PooledObject, error) {
	return pool.NewPooledObject(&MyPoolObject{}), nil
}

func (f *MyObjectFactory) DestroyObject(object *pool.PooledObject) error {
	//do destroy
	return nil
}

func (f *MyObjectFactory) ValidateObject(object *pool.PooledObject) bool {
	//do validate
	return true
}

func (f *MyObjectFactory) ActivateObject(object *pool.PooledObject) error {
	//do activate
	return nil
}

func (f *MyObjectFactory) PassivateObject(object *pool.PooledObject) error {
	//do passivate
	return nil
}

func TestCustomFactoryExample(t *testing.T) {
	p := pool.NewObjectPoolWithDefaultConfig(new(MyObjectFactory))
	obj, _ := p.BorrowObject()
	p.ReturnObject(obj)
}

func TestStringExample(t *testing.T) {
	p := pool.NewObjectPoolWithDefaultConfig(pool.NewPooledObjectFactorySimple(
		func() (interface{}, error) {
			var stringPointer = new(string)
			*stringPointer = "hello"
			return stringPointer, nil
		}))
	obj, _ := p.BorrowObject()
	fmt.Println(obj)
	assert.Equal(t, "hello", *obj.(*string))
	p.ReturnObject(obj)
}
