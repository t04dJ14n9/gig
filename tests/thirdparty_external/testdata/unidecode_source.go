package main

import "github.com/mozillazg/go-unidecode"

func UnidecodeUnidecode() string {
	return unidecode.Unidecode("北京kožušček")
}

func UnidecodePlain() string {
	return unidecode.Unidecode("Hello World")
}

func UnidecodeVersion() string {
	return unidecode.Version()
}
