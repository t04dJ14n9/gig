package thirdparty

import "path/filepath"

// FilepathJoin tests filepath.Join.
func FilepathJoin() string {
	return filepath.Join("dir1", "dir2", "file.txt")
}

// FilepathBase tests filepath.Base.
func FilepathBase() string {
	return filepath.Base("/path/to/file.txt")
}

// FilepathDir tests filepath.Dir.
func FilepathDir() string {
	return filepath.Dir("/path/to/file.txt")
}

// FilepathExt tests filepath.Ext.
func FilepathExt() string {
	return filepath.Ext("/path/to/file.txt")
}

// FilepathClean tests filepath.Clean.
func FilepathClean() string {
	return filepath.Clean("/path/../path/to/./file.txt")
}

// FilepathAbs tests filepath.Abs.
func FilepathAbs() string {
	abs, _ := filepath.Abs("file.txt")
	return abs
}

// FilepathRel tests filepath.Rel.
func FilepathRel() string {
	rel, _ := filepath.Rel("/path/to", "/path/to/file.txt")
	return rel
}

// FilepathSplit tests filepath.Split.
func FilepathSplit() string {
	dir, file := filepath.Split("/path/to/file.txt")
	return dir + file
}

// FilepathFromSlash tests filepath.FromSlash.
func FilepathFromSlash() string {
	return filepath.FromSlash("path/to/file")
}
