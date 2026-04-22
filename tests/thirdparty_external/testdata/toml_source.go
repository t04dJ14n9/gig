package main

import "github.com/pelletier/go-toml/v2"

func TomlMarshal() string {
	type Config struct {
		Name    string
		Version int
	}
	cfg := Config{Name: "test", Version: 2}
	b, err := toml.Marshal(cfg)
	if err != nil {
		return "ERR"
	}
	return string(b)
}

func TomlUnmarshal() string {
	type Config struct {
		Name    string
		Version int
	}
	data := `Name = "hello"
Version = 42`
	var cfg Config
	err := toml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		return "ERR"
	}
	return cfg.Name
}
