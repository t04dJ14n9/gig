package vm

import (
	"reflect"
	"strings"
)

func externalBoundaryReflectArgType(fnType reflect.Type, argIndex int) reflect.Type {
	if fnType == nil || fnType.Kind() != reflect.Func || argIndex < 0 {
		return nil
	}
	numIn := fnType.NumIn()
	if argIndex < numIn {
		if fnType.IsVariadic() && argIndex == numIn-1 {
			return fnType.In(argIndex).Elem()
		}
		return fnType.In(argIndex)
	}
	if fnType.IsVariadic() && numIn > 0 {
		return fnType.In(numIn - 1).Elem()
	}
	return nil
}

func isStdlibExternalPath(path string) bool {
	if path == "" || path == "command-line-arguments" || path == "main" {
		return true
	}
	firstSlash := strings.IndexByte(path, '/')
	firstSegment := path
	if firstSlash >= 0 {
		firstSegment = path[:firstSlash]
	}
	return !strings.ContainsRune(firstSegment, '.')
}
