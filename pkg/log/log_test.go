package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewZapLog(t *testing.T) {
	logger := NewZapLog(defaultTestConfig)
	assert.NotNil(t, logger)
}

func Test_newZapLogWithCallerSkip(t *testing.T) {
	type args struct {
		c          *OutputConfig
		callerSkip int
	}
	tests := []struct {
		name string
		args args
		want Logger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, newZapLogWithCallerSkip(tt.args.c, tt.args.callerSkip), "newZapLogWithCallerSkip(%v, %v)", tt.args.c, tt.args.callerSkip)
		})
	}
}
