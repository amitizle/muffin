package logger

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

type logLine struct {
	Level   string
	Message string
}

func TestLoggerInContext(t *testing.T) {
	ctx := context.Background()

	buf := bytes.NewBufferString("")
	writer := bufio.NewWriter(buf)
	reader := bufio.NewReader(buf)
	readWriter := bufio.NewReadWriter(reader, writer)
	l := zerolog.New(readWriter).Level(zerolog.DebugLevel)
	l.Debug().Msg("a log line")
	newCtx := StoreContext(ctx, l)
	l2, err := GetContext(newCtx)
	l2.Debug().Msg("second line")
	if err != nil {
		t.Fatalf("unexpected error when getting logger from context: %v", err)
	}
	readWriter.Flush()
	var firstLogLine, secondLogLine logLine
	firstLine, _, _ := readWriter.ReadLine()
	fmt.Println(firstLine)
	secondLine, _, _ := readWriter.ReadLine()
	if err := json.Unmarshal(firstLine, &firstLogLine); err != nil {
		t.Fatalf("error while unmarshaling log line to a struct: %v", err)
	}
	if err := json.Unmarshal(secondLine, &secondLogLine); err != nil {
		t.Fatalf("error while unmarshaling log line to a struct: %v", err)
	}
	if strings.ToUpper(firstLogLine.Level) != "DEBUG" {
		t.Fatalf("expected log level line to be info, got %s", firstLogLine.Level)
	}

	if firstLogLine.Message != "a log line" {
		t.Fatalf("expected the log msg to be 'a log line', got %s", firstLogLine.Message)
	}

	if secondLogLine.Message != "second line" {
		t.Fatalf("expected the second log msg to be 'second line', got %s", secondLogLine.Message)
	}
}
