package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/facebookgo/inject"
)

type TypeAnswerStruct struct {
	answer  int
	private int
}

func (t *TypeAnswerStruct) Answer() int {
	return t.answer
}

type TypeNesteStruct struct {
	A *TypeAnswerStruct `inject:""`
}

func (t *TypeNesteStruct) Answer() int {
	return t.A.Answer()
}

var c struct {
	A *TypeNesteStruct `inject:"foo"`
}

var v struct {
	A *TypeAnswerStruct `inject:""`
	B *TypeNesteStruct  `inject:""`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 简单注入操作
func TestInjectSimple(t *testing.T) {
	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}

	fmt.Println(v.B.A.Answer())
}

//测试实例别名依赖
func TestNamedInstanceWithDependencies(t *testing.T) {
	var g inject.Graph
	a1 := &TypeAnswerStruct{answer: 100}
	a2 := &TypeNesteStruct{}
	if err := g.Provide(&inject.Object{Value: a1}, &inject.Object{Value: a2, Name: "foo"}); err != nil {
		t.Fatal(err)
	}

	var c struct {
		A *TypeNesteStruct `inject:"foo"`
	}
	if err := g.Provide(&inject.Object{Value: &c}); err != nil {
		t.Fatal(err)
	}

	if err := g.Populate(); err != nil {
		t.Fatal(err)
	}

	if c.A.A == nil {
		t.Fatal("c.A.A was not injected")
	}

	fmt.Println(c.A.A.Answer())
}
