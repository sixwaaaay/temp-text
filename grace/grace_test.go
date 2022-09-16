package grace

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func TestShutdown_Run(t *testing.T) {
	e := make(chan os.Signal)
	var run, onEv atomic.Bool
	endless := NewEndless(func() {
		run.Store(true)
	}, func() {
		onEv.Store(true)
		run.Store(false)
	}, func() chan os.Signal {
		return e
	})
	go endless.Run()
	time.Sleep(time.Millisecond)
	assert.True(t, run.Load())
	assert.False(t, onEv.Load())
	e <- os.Interrupt
	time.Sleep(time.Millisecond)
	assert.True(t, onEv.Load())
	assert.False(t, run.Load())
}
