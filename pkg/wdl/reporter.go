package wdl

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"os"
)

// An IVisitorReporter that uses the fmt package to report messages.
type FmtReporter struct{}

func (r *FmtReporter) Warn(ctx antlr.ParserRuleContext, warn string) {
	fmt.Fprintf(os.Stderr, "WARN: %v\n", warn)
}

func (r *FmtReporter) Error(ctx antlr.ParserRuleContext, err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

func (r *FmtReporter) NotImplemented(message fmt.Stringer) {
	fmt.Fprintf(os.Stderr, "NotImplemented: %v\n", message)
}
