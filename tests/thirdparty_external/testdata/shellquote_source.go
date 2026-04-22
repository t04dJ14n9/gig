package main

import "github.com/kballard/go-shellquote"

func ShellquoteJoin() string {
	return shellquote.Join("echo", "hello world", "foo'bar")
}

func ShellquoteSplit() int {
	words, err := shellquote.Split(`echo "hello world"`)
	if err != nil {
		return -1
	}
	return len(words)
}
