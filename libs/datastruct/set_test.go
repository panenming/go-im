package datastruct

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	//初始化
	s := NewSet()

	s.Add("1")
	s.Add("1")
	s.Add("0")
	s.Add("2")
	s.Add("3")
	s.Add("4")
	fmt.Println("s list = ", s.List())

	s.Clear()
	if s.IsEmpty() {
		fmt.Println("0 item")
	}

	s.Add("1")
	s.Add("2")
	s.Add("3")

	fmt.Println("s list = ", s.List())

	if s.Has("2") {
		fmt.Println("2 does exist")
	}

	s.Remove("2")
	s.Remove("3")
	s.Add("0")
	fmt.Println("无序的切片", s.List())
	fmt.Println("有序的切片", s.SortList())

}
