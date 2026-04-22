package main

import "github.com/mohae/deepcopy"

func DeepcopyCopy() string {
	original := map[string]string{"key": "value"}
	copied := deepcopy.Copy(original).(map[string]string)
	return copied["key"]
}

func DeepcopyIface() string {
	original := map[string]string{"key": "value"}
	copied := deepcopy.Iface(original).(map[string]string)
	return copied["key"]
}
