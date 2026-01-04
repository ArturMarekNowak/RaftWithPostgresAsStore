package setup

import (
	"github.com/hashicorp/go-hclog"
	"os"
)

func ConfigureLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:                     "",
		Level:                    3,
		Output:                   os.Stdout,
		Mutex:                    nil,
		JSONFormat:               true,
		JSONEscapeDisabled:       false,
		IncludeLocation:          false,
		AdditionalLocationOffset: 0,
		TimeFormat:               "",
		TimeFn:                   nil,
		DisableTime:              false,
		Color:                    0,
		ColorHeaderOnly:          false,
		ColorHeaderAndFields:     false,
		Exclude:                  nil,
		IndependentLevels:        false,
		SyncParentLevel:          false,
		SubloggerHook:            nil,
	})
}
