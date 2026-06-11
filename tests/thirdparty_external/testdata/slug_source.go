package main

import "github.com/gosimple/slug"

func SlugMake() string {
	return slug.Make("Hellö Wörld хелло ворлд")
}

func SlugMakeLang() string {
	return slug.MakeLang("Diese & Dass", "de")
}

func SlugIsSlug() bool {
	return slug.IsSlug("hello-world")
}

func SlugIsSlugInvalid() bool {
	return slug.IsSlug("Hello World!")
}

func SlugSubstitute() string {
	return slug.Substitute("Hello World", map[string]string{"World": "Go"})
}
