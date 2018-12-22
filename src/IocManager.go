package main

import (
	"fmt"
	"github.com/facebookgo/inject"
	"math/rand"
	"testing"
	"time"
)

type TypeAnswerStruct struct {
	answer  int
	private int
}

func (t *TypeAnswerStruct) Answer() int {
	return t.answer
}

type TypeNestedStruct struct {
	A *TypeAnswerStruct `inject:""`
}

func (t *TypeNestedStruct) Answer() int {
	return t.A.Answer()
}

var c struct {
	A *TypeNestedStruct `inject:"foo"`
}

var v struct {
	A *TypeAnswerStruct `inject:""`
	B *TypeNestedStruct `inject:""`
}

func init() {
	// we rely on math.Rand in Graph.Objects() and this gives it some randomness.
	//使用 math.Rand 在 Graph.Objects() 给出一些随机性
	rand.Seed(time.Now().UnixNano())
}

//简单注入操作
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
	a2 := &TypeNestedStruct{}
	if err := g.Provide(&inject.Object{Value: a1}, &inject.Object{Value: a2, Name: "foo"}); err != nil {
		t.Fatal(err)
	}

	var c struct {
		A *TypeNestedStruct `inject:"foo"`
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
func main() {
	//TestNamedInstanceWithDependencies(&testing.T{})
	TestInjectSimple(&testing.T{})
}
