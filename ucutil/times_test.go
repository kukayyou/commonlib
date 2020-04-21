package ucutil

import (
	"testing"
)

func TestGetDateMaxTime(t *testing.T) {
	time := GetDateMaxTime(1555516800000)
	t.Logf("time: %d", time)
}

func TestFormatUnixTime(t *testing.T) {
	time := FormatUnixTime(1555516800)
	t.Logf("time: %s", time)
}
