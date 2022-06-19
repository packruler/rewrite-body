package logger

import (
	"bytes"
	"testing"
)

func TestLogging(t *testing.T) {
	tests := []struct {
		desc        string
		configLevel LogLevel
		targetLevel LogLevel
		shouldWrite bool
	}{
		{
			desc:        "Trace level should write when targeting Trace",
			configLevel: Trace,
			targetLevel: Trace,
			shouldWrite: true,
		},
		{
			desc:        "Trace level should write when targeting Debug",
			configLevel: Trace,
			targetLevel: Debug,
			shouldWrite: true,
		},
		{
			desc:        "Trace level should write when targeting Info",
			configLevel: Trace,
			targetLevel: Info,
			shouldWrite: true,
		},
		{
			desc:        "Trace level should write when targeting Warning",
			configLevel: Trace,
			targetLevel: Warning,
			shouldWrite: true,
		},
		{
			desc:        "Trace level should write when targeting Error",
			configLevel: Trace,
			targetLevel: Error,
			shouldWrite: true,
		},
		{
			desc:        "Debug level should NOT write when targeting Trace",
			configLevel: Debug,
			targetLevel: Trace,
			shouldWrite: false,
		},
		{
			desc:        "Debug level should write when targeting Debug",
			configLevel: Debug,
			targetLevel: Debug,
			shouldWrite: true,
		},
		{
			desc:        "Debug level should write when targeting Info",
			configLevel: Debug,
			targetLevel: Info,
			shouldWrite: true,
		},
		{
			desc:        "Debug level should write when targeting Warning",
			configLevel: Debug,
			targetLevel: Warning,
			shouldWrite: true,
		},
		{
			desc:        "Debug level should write when targeting Error",
			configLevel: Debug,
			targetLevel: Error,
			shouldWrite: true,
		},
		{
			desc:        "Info level should NOT write when targeting Trace",
			configLevel: Info,
			targetLevel: Trace,
			shouldWrite: false,
		},
		{
			desc:        "Info level should NOT write when targeting Debug",
			configLevel: Info,
			targetLevel: Debug,
			shouldWrite: false,
		},
		{
			desc:        "Info level should write when targeting Info",
			configLevel: Info,
			targetLevel: Info,
			shouldWrite: true,
		},
		{
			desc:        "Info level should write when targeting Warning",
			configLevel: Info,
			targetLevel: Warning,
			shouldWrite: true,
		},
		{
			desc:        "Info level should write when targeting Error",
			configLevel: Info,
			targetLevel: Error,
			shouldWrite: true,
		},
		{
			desc:        "Warning level should NOT write when targeting Trace",
			configLevel: Warning,
			targetLevel: Trace,
			shouldWrite: false,
		},
		{
			desc:        "Warning level should NOT write when targeting Debug",
			configLevel: Warning,
			targetLevel: Debug,
			shouldWrite: false,
		},
		{
			desc:        "Warning level should NOT write when targeting Info",
			configLevel: Warning,
			targetLevel: Info,
			shouldWrite: false,
		},
		{
			desc:        "Warning level should write when targeting Warning",
			configLevel: Warning,
			targetLevel: Warning,
			shouldWrite: true,
		},
		{
			desc:        "Warning level should write when targeting Error",
			configLevel: Warning,
			targetLevel: Error,
			shouldWrite: true,
		},
		{
			desc:        "Warning level should NOT write when targeting Trace",
			configLevel: Warning,
			targetLevel: Trace,
			shouldWrite: false,
		},
		{
			desc:        "Warning level should NOT write when targeting Debug",
			configLevel: Warning,
			targetLevel: Debug,
			shouldWrite: false,
		},
		{
			desc:        "Warning level should NOT write when targeting Info",
			configLevel: Warning,
			targetLevel: Info,
			shouldWrite: false,
		},
		{
			desc:        "Warning level should write when targeting Warning",
			configLevel: Warning,
			targetLevel: Warning,
			shouldWrite: true,
		},
		{
			desc:        "Warning level should write when targeting Error",
			configLevel: Warning,
			targetLevel: Error,
			shouldWrite: true,
		},
		{
			desc:        "Error level should NOT write when targeting Trace",
			configLevel: Error,
			targetLevel: Trace,
			shouldWrite: false,
		},
		{
			desc:        "Error level should NOT write when targeting Debug",
			configLevel: Error,
			targetLevel: Debug,
			shouldWrite: false,
		},
		{
			desc:        "Error level should NOT write when targeting Info",
			configLevel: Error,
			targetLevel: Info,
			shouldWrite: false,
		},
		{
			desc:        "Error level should NOT write when targeting Warning",
			configLevel: Error,
			targetLevel: Warning,
			shouldWrite: false,
		},
		{
			desc:        "Error level should write when targeting Error",
			configLevel: Error,
			targetLevel: Error,
			shouldWrite: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			var buffer bytes.Buffer

			logger := createLoggerWithBuffer(test.configLevel, &buffer)

			switch test.targetLevel {
			case Trace:
				logger.LogTrace("test")
			case Debug:
				logger.LogDebug("test")
			case Info:
				logger.LogInfo("test")
			case Warning:
				logger.LogWarning("test")
			case Error:
				logger.LogError("test")
			}

			loggedContent := buffer.Bytes()
			if !test.shouldWrite {
				if len(loggedContent) != 0 {
					t.Errorf("Log level %v written with test value set to %v: '%s'",
						test.targetLevel, test.configLevel, loggedContent)
				}
			} else {
				if len(loggedContent) == 0 {
					t.Errorf("Log level %v NOT written with test value set to %v: '%s'",
						test.targetLevel, test.configLevel, loggedContent)
				}
			}
		})
	}
}
