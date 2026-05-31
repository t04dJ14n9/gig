package vm

import (
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func TestInterpreterDefinedReflectValueTypeFindsNestedInterfaceValue(t *testing.T) {
	scriptType := reflect.StructOf([]reflect.StructField{
		{Name: "Field", Type: reflect.TypeOf(0)},
	})
	scriptValue := reflect.New(scriptType).Elem()
	scriptValue.Field(0).SetInt(7)

	prog := &bytecode.CompiledProgram{}
	prog.RegisterTypeName(scriptType, "BoundaryScriptStruct")
	v := &vm{program: prog}

	payload := map[string]any{
		"nested": scriptValue.Interface(),
	}

	got, ok := v.interpreterDefinedReflectValueType(reflect.ValueOf(payload), make(map[reflect.Type]bool), 0)

	if !ok {
		t.Fatal("interpreter-defined type hidden behind interface was not detected")
	}
	if got != "BoundaryScriptStruct" {
		t.Fatalf("type name = %q, want BoundaryScriptStruct", got)
	}
}

func TestInterpreterDefinedReflectValueTypeIgnoresNilDynamicValues(t *testing.T) {
	v := &vm{program: &bytecode.CompiledProgram{}}
	payload := map[string]any{
		"nil": nil,
	}

	_, ok := v.interpreterDefinedReflectValueType(reflect.ValueOf(payload), make(map[reflect.Type]bool), 0)

	if ok {
		t.Fatal("nil dynamic value reported as interpreter-defined type")
	}
}
