package resp

import (
	"strings"
	"testing"
)

func TestReadCommand(t *testing.T) {
	input := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	r := NewReader(strings.NewReader(input))

	parts, err := r.ReadCommand()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"SET", "foo", "bar"}
	if len(parts) != len(want) {
		t.Fatalf("got %d parts, want %d: %q", len(parts), len(want), parts)
	}
	for i := range want {
		if parts[i] != want[i] {
			t.Errorf("part %d: got %q, want %q", i, parts[i], want[i])
		}
	}
}
