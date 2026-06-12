package thirdparty_complex

// This file lists third-party packages to generate wrappers for
// Each import generates a MethodResolver for calling host methods

import (
	// encoding
	_ "encoding/base64"
	_ "encoding/binary"

	// archive
	_ "archive/tar"
	_ "archive/zip"

	// text
	_ "text/scanner"
	_ "text/tabwriter"
	_ "text/template"

	// html (for template)
	_ "html/template"
)
