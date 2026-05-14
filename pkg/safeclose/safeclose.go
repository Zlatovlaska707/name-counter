package safeclose

import (
	"fmt"
	"io"
	"log"
	"os"
)

func writeToStderr(msg string) {
	if _, err := fmt.Fprint(os.Stderr, msg); err != nil {
		log.Printf("failed to write to stderr: %v", err)
		log.Printf("%s", msg)
	}
}

func Close(closer io.Closer, name string) {
	if err := closer.Close(); err != nil {
		writeToStderr(fmt.Sprintf("Предупреждение: не удалось закрыть %s: %v\n", name, err))
	}
}

func CloseFile(file interface{ Close() error }, path string) {
	if err := file.Close(); err != nil {
		writeToStderr(fmt.Sprintf("Предупреждение: не удалось закрыть файл %s: %v\n", path, err))
	}
}

func CloseWithLog(closer io.Closer, format string, args ...interface{}) {
	if err := closer.Close(); err != nil {
		msg := fmt.Sprintf(format+": %v\n", append(args, err)...)
		writeToStderr(msg)
	}
}
