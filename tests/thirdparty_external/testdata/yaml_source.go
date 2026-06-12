package main

import "gopkg.in/yaml.v3"

func YamlMarshal() string {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Bob", Age: 30}
	b, err := yaml.Marshal(p)
	if err != nil {
		return "ERR"
	}
	return string(b)
}

func YamlUnmarshal() string {
	type Person struct {
		Name string
		Age  int
	}
	data := "name: Alice\nage: 25\n"
	var p Person
	err := yaml.Unmarshal([]byte(data), &p)
	if err != nil {
		return "ERR"
	}
	return p.Name
}
