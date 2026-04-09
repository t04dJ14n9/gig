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
	_ "text/template"
	_ "text/tabwriter"
	_ "text/scanner"
	
	// html (for template)
	_ "html/template"
)
