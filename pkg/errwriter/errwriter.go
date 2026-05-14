package errwriter

import (
	"fmt"
	"log"
	"os"
)

func PrintErr(a ...interface{}) {
	if _, err := fmt.Fprint(os.Stderr, a...); err != nil {
		log.Printf("failed to write to stderr: %v", err)
	}
}

func PrintErrln(a ...interface{}) {
	if _, err := fmt.Fprintln(os.Stderr, a...); err != nil {
		log.Printf("failed to write to stderr: %v", err)
	}
}

func PrintErrf(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, format, a...); err != nil {
		log.Printf("failed to write to stderr: %v", err)
	}
}
