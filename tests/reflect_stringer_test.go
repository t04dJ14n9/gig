package tests

import (
	"fmt"
	"reflect"
	"testing"
)

// TestReflectCallStringer 验证通过 reflect.Call 调用 fmt.Sprint 时，
// 结构体的 String() 方法是否被正确调用
//
// 结论：
// 1. 值接收器 String(): 无论传值还是传指针，都会被调用
// 2. 指针接收器 String(): 只有传指针时才会被调用，传值不会
// 3. reflect.Call 行为与直接调用完全一致

type valueStringer struct {
	Name string
}

func (s valueStringer) String() string {
	return fmt.Sprintf("[valueStringer: %s]", s.Name)
}

type pointerStringer struct {
	Name string
}

func (s *pointerStringer) String() string {
	return fmt.Sprintf("[pointerStringer: %s]", s.Name)
}

func TestReflectCall_ValueStringer(t *testing.T) {
	v := valueStringer{Name: "Alice"}

	// 直接调用
	direct := fmt.Sprintf("%v", v)

	// reflect.Call 调用
	sprintFunc := reflect.ValueOf(fmt.Sprint)
	args := []reflect.Value{reflect.ValueOf(v)}
	result := sprintFunc.Call(args)
	viaReflect := result[0].String()

	expected := "[valueStringer: Alice]"
	if direct != expected {
		t.Errorf("直接调用: 期望 %q, 得到 %q", expected, direct)
	}
	if viaReflect != expected {
		t.Errorf("reflect.Call: 期望 %q, 得到 %q", expected, viaReflect)
	}
	t.Logf("值接收器: 直接=%q, reflect=%q", direct, viaReflect)
}

func TestReflectCall_PointerStringer_Value(t *testing.T) {
	v := pointerStringer{Name: "Bob"}

	// 直接调用（传值）- 不调用 String()
	direct := fmt.Sprintf("%v", v)

	// reflect.Call 调用（传值）
	sprintFunc := reflect.ValueOf(fmt.Sprint)
	args := []reflect.Value{reflect.ValueOf(v)}
	result := sprintFunc.Call(args)
	viaReflect := result[0].String()

	// 传值时，指针接收器的方法不可用，打印默认格式
	expected := "{Bob}"
	if direct != expected {
		t.Errorf("直接调用(传值): 期望 %q, 得到 %q", expected, direct)
	}
	if viaReflect != expected {
		t.Errorf("reflect.Call(传值): 期望 %q, 得到 %q", expected, viaReflect)
	}
	t.Logf("指针接收器(传值): 直接=%q, reflect=%q", direct, viaReflect)
}

func TestReflectCall_PointerStringer_Pointer(t *testing.T) {
	v := pointerStringer{Name: "Charlie"}

	// 直接调用（传指针）- 调用 String()
	direct := fmt.Sprintf("%v", &v)

	// reflect.Call 调用（传指针）
	sprintFunc := reflect.ValueOf(fmt.Sprint)
	args := []reflect.Value{reflect.ValueOf(&v)}
	result := sprintFunc.Call(args)
	viaReflect := result[0].String()

	// 传指针时，指针接收器的方法可用
	expected := "[pointerStringer: Charlie]"
	if direct != expected {
		t.Errorf("直接调用(传指针): 期望 %q, 得到 %q", expected, direct)
	}
	if viaReflect != expected {
		t.Errorf("reflect.Call(传指针): 期望 %q, 得到 %q", expected, viaReflect)
	}
	t.Logf("指针接收器(传指针): 直接=%q, reflect=%q", direct, viaReflect)
}

// TestReflectCall_Interface 演示 interface{} 参数的行为
func TestReflectCall_Interface(t *testing.T) {
	// fmt.Sprint 的参数类型是 ...interface{}
	// 当传入结构体时，会被装箱为 interface{}

	v := valueStringer{Name: "Dave"}

	// 即使装箱为 interface{}，String() 仍会被调用
	// 因为 fmt 包会检查 interface{} 值是否实现了 fmt.Stringer
	direct := fmt.Sprint(v)
	viaReflect := reflect.ValueOf(fmt.Sprint).Call([]reflect.Value{reflect.ValueOf(v)})[0].String()

	expected := "[valueStringer: Dave]"
	if direct != expected {
		t.Errorf("直接调用: 期望 %q, 得到 %q", expected, direct)
	}
	if viaReflect != expected {
		t.Errorf("reflect.Call: 期望 %q, 得到 %q", expected, viaReflect)
	}
	t.Logf("interface{} 参数: 直接=%q, reflect=%q", direct, viaReflect)
}
