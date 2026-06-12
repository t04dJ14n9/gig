package vm

func runJumpIf(frame *Frame, offset uint16, cond bool) {
	if cond {
		frame.ip = int(offset)
	}
}
