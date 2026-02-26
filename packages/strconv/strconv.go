// Package strconv registers the Go standard library strconv package.
package strconv

import (
	"strconv"

	"gig/importer"
	"gig/value"
)

func init() {
	pkg := importer.RegisterPackage("strconv", "strconv")

	// Parsing
	pkg.AddFunction("ParseBool", strconv.ParseBool, "", directParseBool)
	pkg.AddFunction("ParseInt", strconv.ParseInt, "", nil)
	pkg.AddFunction("ParseUint", strconv.ParseUint, "", nil)
	pkg.AddFunction("ParseFloat", strconv.ParseFloat, "", nil)
	pkg.AddFunction("ParseComplex", strconv.ParseComplex, "", nil)

	// Formatting
	pkg.AddFunction("FormatBool", strconv.FormatBool, "", directFormatBool)
	pkg.AddFunction("FormatInt", strconv.FormatInt, "", directFormatInt)
	pkg.AddFunction("FormatUint", strconv.FormatUint, "", directFormatUint)
	pkg.AddFunction("FormatFloat", strconv.FormatFloat, "", nil)
	pkg.AddFunction("FormatComplex", strconv.FormatComplex, "", nil)

	// Conversions
	pkg.AddFunction("Atoi", strconv.Atoi, "", nil)
	pkg.AddFunction("Itoa", strconv.Itoa, "", directItoa)
	pkg.AddFunction("AppendBool", strconv.AppendBool, "", nil)
	pkg.AddFunction("AppendInt", strconv.AppendInt, "", nil)
	pkg.AddFunction("AppendUint", strconv.AppendUint, "", nil)
	pkg.AddFunction("AppendFloat", strconv.AppendFloat, "", nil)
	pkg.AddFunction("AppendQuote", strconv.AppendQuote, "", nil)
	pkg.AddFunction("AppendQuoteRune", strconv.AppendQuoteRune, "", nil)
	pkg.AddFunction("AppendQuoteRuneToASCII", strconv.AppendQuoteRuneToASCII, "", nil)
	pkg.AddFunction("AppendQuoteRuneToGraphic", strconv.AppendQuoteRuneToGraphic, "", nil)
	pkg.AddFunction("AppendQuoteToASCII", strconv.AppendQuoteToASCII, "", nil)
	pkg.AddFunction("AppendQuoteToGraphic", strconv.AppendQuoteToGraphic, "", nil)

	// Quoting
	pkg.AddFunction("Quote", strconv.Quote, "", directQuote)
	pkg.AddFunction("QuoteToASCII", strconv.QuoteToASCII, "", directQuoteToASCII)
	pkg.AddFunction("QuoteToGraphic", strconv.QuoteToGraphic, "", directQuoteToGraphic)
	pkg.AddFunction("QuoteRune", strconv.QuoteRune, "", nil)
	pkg.AddFunction("QuoteRuneToASCII", strconv.QuoteRuneToASCII, "", nil)
	pkg.AddFunction("QuoteRuneToGraphic", strconv.QuoteRuneToGraphic, "", nil)
	pkg.AddFunction("Unquote", strconv.Unquote, "", nil)
	pkg.AddFunction("UnquoteChar", strconv.UnquoteChar, "", nil)

	// Errors
	pkg.AddFunction("CanBackquote", strconv.CanBackquote, "", directCanBackquote)
	pkg.AddFunction("IsGraphic", strconv.IsGraphic, "", nil)
	pkg.AddFunction("IsPrint", strconv.IsPrint, "", nil)

	// Types
	pkg.AddType("NumError", nil, "numeric parsing error")
}

// Direct wrappers

func directParseBool(args []value.Value) value.Value {
	b, err := strconv.ParseBool(args[0].String())
	return value.FromInterface([]any{b, err})
}

func directFormatBool(args []value.Value) value.Value {
	return value.MakeString(strconv.FormatBool(args[0].Bool()))
}

func directFormatInt(args []value.Value) value.Value {
	return value.MakeString(strconv.FormatInt(args[0].Int(), int(args[1].Int())))
}

func directFormatUint(args []value.Value) value.Value {
	return value.MakeString(strconv.FormatUint(args[0].Uint(), int(args[1].Int())))
}

func directItoa(args []value.Value) value.Value {
	return value.MakeString(strconv.Itoa(int(args[0].Int())))
}

func directQuote(args []value.Value) value.Value {
	return value.MakeString(strconv.Quote(args[0].String()))
}

func directQuoteToASCII(args []value.Value) value.Value {
	return value.MakeString(strconv.QuoteToASCII(args[0].String()))
}

func directQuoteToGraphic(args []value.Value) value.Value {
	return value.MakeString(strconv.QuoteToGraphic(args[0].String()))
}

func directCanBackquote(args []value.Value) value.Value {
	return value.MakeBool(strconv.CanBackquote(args[0].String()))
}
