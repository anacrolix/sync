package sync

import (
	"testing"

	"github.com/anacrolix/envpprof"
)

func TestMain(m *testing.M) {
	println("sync test main")
	envpprof.TestMain(m)
}
