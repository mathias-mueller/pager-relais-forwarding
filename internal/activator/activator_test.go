package activator

import (
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

type testActivation struct {
	count int
	lock  sync.Mutex
}

func (a *testActivation) activate() {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.count++
}

func TestActivator_EnableActivation(t *testing.T) {
	type args struct {
		inputs []bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "direct activation",
			args: args{inputs: []bool{true}},
			want: 1,
		},
		{
			name: "no activation",
			args: args{inputs: []bool{false}},
			want: 0,
		},
		{
			name: "delayed activation",
			args: args{inputs: []bool{false, true}},
			want: 1,
		},
		{
			name: "multiple true inputs",
			args: args{inputs: []bool{true, true, true}},
			want: 1,
		},
		{
			name: "multiple activations",
			args: args{inputs: []bool{true, false, true}},
			want: 2,
		},
		{
			name: "complex",
			args: args{inputs: []bool{false, true, true, true, false, false, true, true}},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activator := New()

			testChan := make(chan bool)
			t2 := &testActivation{
				count: 0,
				lock:  sync.Mutex{},
			}
			defer close(testChan)
			go activator.EnableActivation(
				testChan,
				[]Activation{t2},
			)
			for _, i := range tt.args.inputs {
				testChan <- i
			}
			assert.Eventually(t, func() bool {
				t2.lock.Lock()
				defer t2.lock.Unlock()
				return tt.want == t2.count
			},
				time.Millisecond*100,
				time.Millisecond,
			)
			prometheus.Unregister(activator.counter)
		})
	}
}
