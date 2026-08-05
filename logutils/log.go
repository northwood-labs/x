// Copyright 2026, Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logutils

import (
	"log/slog"
	"os"
	"time"

	"charm.land/log/v2"

	clihelpers "go.nwlabs.dev/cli-helpers/v2"
)

const (
	noVerbose    = 0
	verboseInfo  = 1
	verboseDebug = 2
)

// charmlogger defaults to stderr with caller info and timestamps. The
// PersistentPreRunE adjusts the level based on -v count so that normal
// usage is quiet and -vvv gives full debug output with source locations.
var charmlogger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
})

// GetCharmLogger creates a Charmbracelet logger, with default settings applied.
// Only use GetCharmLogger if you need to edit underlying Charmbracelet logger
// settings. Prefer using GetDefaultLogger.
func GetCharmLogger(fVerbose int) *log.Logger {
	charmlogger.SetStyles(clihelpers.GetLoggerStyles())

	if fVerbose < 0 {
		fVerbose = 0
	} else if fVerbose > 3 {
		fVerbose = 3
	}

	// Verbose levels: 0=warn (quiet default), 1=info (-v), 2=debug
	// (-vv), 3+=debug with source file:line (-vvv). ReportCaller is
	// expensive so it's only enabled at the highest level for deep
	// debugging.
	switch fVerbose {
	case noVerbose:
		charmlogger.SetLevel(log.WarnLevel)
		charmlogger.SetReportCaller(false)
	case verboseInfo:
		charmlogger.SetLevel(log.InfoLevel)
		charmlogger.SetReportCaller(false)
	case verboseDebug:
		charmlogger.SetLevel(log.DebugLevel)
		charmlogger.SetReportCaller(false)
	default:
		charmlogger.SetLevel(log.DebugLevel)
		charmlogger.SetReportCaller(true)
	}

	return charmlogger
}

// GetDefaultLogger wraps the Charmbracelet logger in a slog interface.
func GetDefaultLogger(fVerbose int) *slog.Logger {
	return slog.New(charmlogger)
}
