package bytecode

// SelectMeta stores metadata for a select statement.
// It is stored in the constant pool and referenced by OpSelect.
type SelectMeta struct {
	NumStates int    // total number of select cases (excluding default)
	Blocking  bool   // true if no default branch
	Dirs      []bool // direction for each case: true=send, false=recv
	NumRecv   int    // number of recv cases (determines result tuple size)
}
