package main

import "github.com/json-iterator/go"

func JsoniterMarshal() string {
	data, err := jsoniter.Marshal(map[string]string{"key": "value"})
	if err != nil {
		return "ERR"
	}
	return string(data)
}

func JsoniterMarshalIndent() string {
	data, err := jsoniter.MarshalIndent(map[string]string{"key": "value"}, "", "  ")
	if err != nil {
		return "ERR"
	}
	return string(data)
}

func JsoniterMarshalToString() string {
	s, err := jsoniter.MarshalToString(map[string]string{"key": "value"})
	if err != nil {
		return "ERR"
	}
	return s
}

func JsoniterUnmarshal() string {
	var m map[string]string
	err := jsoniter.Unmarshal([]byte(`{"key":"value"}`), &m)
	if err != nil {
		return "ERR"
	}
	return m["key"]
}

func JsoniterUnmarshalFromString() string {
	var m map[string]string
	err := jsoniter.UnmarshalFromString(`{"key":"value"}`, &m)
	if err != nil {
		return "ERR"
	}
	return m["key"]
}

func JsoniterValid() bool {
	return jsoniter.Valid([]byte(`{"key":"value"}`))
}

func JsoniterConfigDefaultMarshal() string {
	data, err := jsoniter.ConfigDefault.Marshal(map[string]string{"key": "value"})
	if err != nil {
		return "ERR"
	}
	return string(data)
}

func JsoniterConfigFastestMarshal() string {
	data, err := jsoniter.ConfigFastest.Marshal(map[string]string{"key": "value"})
	if err != nil {
		return "ERR"
	}
	return string(data)
}

func JsoniterConfigCompatibleMarshal() string {
	data, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(map[string]string{"key": "value"})
	if err != nil {
		return "ERR"
	}
	return string(data)
}

func JsoniterGet() string {
	result := jsoniter.Get([]byte(`{"key":"value"}`), "key")
	return result.ToString()
}

func JsoniterWrap() string {
	any := jsoniter.Wrap(float64(3.14))
	return any.ToString()
}

func JsoniterWrapString() string {
	any := jsoniter.WrapString("hello")
	return any.ToString()
}

func JsoniterWrapInt64() string {
	any := jsoniter.WrapInt64(42)
	return any.ToString()
}

func JsoniterWrapUint64() string {
	any := jsoniter.WrapUint64(42)
	return any.ToString()
}

func JsoniterWrapFloat64() string {
	any := jsoniter.WrapFloat64(3.14)
	return any.ToString()
}

func JsoniterWrapInt32() string {
	any := jsoniter.WrapInt32(42)
	return any.ToString()
}

func JsoniterWrapUint32() string {
	any := jsoniter.WrapUint32(42)
	return any.ToString()
}
