package tests

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"git.woa.com/youngjin/gig/model/value"
)

// 原生 Go 结构体
type NativeStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (s NativeStruct) String() string {
	return fmt.Sprintf("[Native: %s, %d]", s.Name, s.Age)
}

// 模拟 gig 创建的匿名结构体（无方法）
func createGigLikeStruct() reflect.Type {
	fields := []reflect.StructField{
		{Name: "Name", Type: reflect.TypeOf(""), Tag: `json:"name" gig:"#test.AnonStruct"`},
		{Name: "Age", Type: reflect.TypeOf(0), Tag: `json:"age" gig:"#test.AnonStruct"`},
	}
	return reflect.StructOf(fields)
}

var anonStructType = createGigLikeStruct()

// Benchmark: ExternWrap 的真实开销（输入已经是 Value）
func BenchmarkExternWrap_Only(b *testing.B) {
	v := reflect.New(anonStructType).Elem()
	v.Field(0).SetString("Alice")
	v.Field(1).SetInt(30)

	// 预先创建 Value 对象（模拟 VM 中的情况）
	val := value.MakeFromReflect(v)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = value.FmtWrap(val)
	}
}

// Benchmark: ExternWrap + isGigStruct 开销分解
func BenchmarkExternWrap_Breakdown(b *testing.B) {
	v := reflect.New(anonStructType).Elem()
	v.Field(0).SetString("Alice")
	v.Field(1).SetInt(30)
	val := value.MakeFromReflect(v)

	b.Run("Interface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = val.Interface()
		}
	})

	b.Run("ExternWrap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = value.FmtWrap(val)
		}
	})
}

// Benchmark: 完整调用链 - 公平比较
func BenchmarkFullChain(b *testing.B) {
	// 原生结构体
	native := NativeStruct{Name: "Alice", Age: 30}

	// gig 匿名结构体
	v := reflect.New(anonStructType).Elem()
	v.Field(0).SetString("Alice")
	v.Field(1).SetInt(30)
	gigVal := value.MakeFromReflect(v)

	b.Run("Native/Direct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(native)
		}
	})

	b.Run("Native/Reflect", func(b *testing.B) {
		marshalFunc := reflect.ValueOf(json.Marshal)
		for i := 0; i < b.N; i++ {
			args := []reflect.Value{reflect.ValueOf(native)}
			marshalFunc.Call(args)
		}
	})

	b.Run("Gig/Reflect_NoWrap", func(b *testing.B) {
		marshalFunc := reflect.ValueOf(json.Marshal)
		for i := 0; i < b.N; i++ {
			args := []reflect.Value{reflect.ValueOf(gigVal.Interface())}
			marshalFunc.Call(args)
		}
	})

	b.Run("Gig/Reflect_WithWrap", func(b *testing.B) {
		marshalFunc := reflect.ValueOf(json.Marshal)
		for i := 0; i < b.N; i++ {
			wrapped := value.FmtWrap(gigVal)
			args := []reflect.Value{reflect.ValueOf(wrapped)}
			marshalFunc.Call(args)
		}
	})
}

// Benchmark: fmt.Sprintf 调用链
func BenchmarkFmtChain(b *testing.B) {
	native := NativeStruct{Name: "Alice", Age: 30}

	v := reflect.New(anonStructType).Elem()
	v.Field(0).SetString("Alice")
	v.Field(1).SetInt(30)
	gigVal := value.MakeFromReflect(v)

	b.Run("Native/Direct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%v", native)
		}
	})

	b.Run("Native/Reflect", func(b *testing.B) {
		sprintfFunc := reflect.ValueOf(fmt.Sprintf)
		for i := 0; i < b.N; i++ {
			args := []reflect.Value{reflect.ValueOf("%v"), reflect.ValueOf(native)}
			sprintfFunc.Call(args)
		}
	})

	b.Run("Gig/Reflect_NoWrap", func(b *testing.B) {
		sprintfFunc := reflect.ValueOf(fmt.Sprintf)
		for i := 0; i < b.N; i++ {
			args := []reflect.Value{reflect.ValueOf("%v"), reflect.ValueOf(gigVal.Interface())}
			sprintfFunc.Call(args)
		}
	})

	b.Run("Gig/Reflect_WithWrap", func(b *testing.B) {
		sprintfFunc := reflect.ValueOf(fmt.Sprintf)
		for i := 0; i < b.N; i++ {
			wrapped := value.FmtWrap(gigVal)
			args := []reflect.Value{reflect.ValueOf("%v"), reflect.ValueOf(wrapped)}
			sprintfFunc.Call(args)
		}
	})
}

// Benchmark: DirectCall vs Reflect 路径
func BenchmarkDirectCall_vs_Reflect(b *testing.B) {
	v := reflect.New(anonStructType).Elem()
	v.Field(0).SetString("Alice")
	v.Field(1).SetInt(30)
	gigVal := value.MakeFromReflect(v)

	b.Run("Reflect_NoWrap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(gigVal.Interface())
		}
	})

	b.Run("Reflect_WithWrap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wrapped := value.FmtWrap(gigVal)
			json.Marshal(wrapped)
		}
	})
}
