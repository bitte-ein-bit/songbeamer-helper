// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package log provides a universal logger for martian packages.
package log

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"cloud.google.com/go/logging"
	"github.com/fatih/color"
	"google.golang.org/api/option"
)

const (
	// Silent is a level that logs nothing.
	Silent int = iota
	// Error is a level that logs error logs.
	Error
	// Warning is a level that logs error, and warning logs.
	Warning
	// Info is a level that logs error, warning, and info logs.
	Info
	// Debug is a level that logs error, info, and debug logs.
	Debug
)

// credsJSON is injected at build time as base64 encoded string. It contains the credentials for the Google Cloud Logging API.
// The file is not included in the repository and should be created by the user. It can be downloaded from the Google Cloud Console.
// The file should be named .googlecloud.json and placed in the root of the project.
var credsJSON string

// GoogleCloudCreds represents the structure of the Google Cloud credentials JSON
type GoogleCloudCreds struct {
	Type       string `json:"type"`
	ProjectID  string `json:"project_id"`
	PrivateKey string `json:"private_key"`
	// Add other fields as needed
}

// Default log level is Info.
var (
	level   = Info
	lock    sync.Mutex
	loggers = make(map[string]*log.Logger) // stores all loggers by level name

	client   *logging.Client
	usingGCP bool = false // tracks whether GCP logging is active
)

func init() {
	// Try to initialize GCP logging, fall back to local logging on failure
	if !initGCPLogging() {
		initLocalLogging()
	}
}

// initGCPLogging attempts to set up Google Cloud Platform logging
// Returns true on success, false on failure
func initGCPLogging() bool {
	var err error
	ctx := context.Background()

	// decode the base64 encoded credentials
	decodedCreds, err := base64.StdEncoding.DecodeString(credsJSON)
	if err != nil {
		log.Printf("Warning: Failed to decode GCP credentials: %v. Falling back to local logging.", err)
		return false
	}

	// parse the JSON credentials
	var creds GoogleCloudCreds
	err = json.Unmarshal(decodedCreds, &creds)
	if err != nil {
		log.Printf("Warning: Failed to parse GCP credentials: %v. Falling back to local logging.", err)
		return false
	}

	// Sets your Google Cloud Platform project ID dynamically from credentials.
	projectID := fmt.Sprintf("projects/%s", creds.ProjectID)

	// Creates a client.
	client, err = logging.NewClient(ctx, projectID, option.WithCredentialsJSON(decodedCreds))
	if err != nil {
		log.Printf("Warning: Failed to create GCP logging client: %v. Falling back to local logging.", err)
		return false
	}

	// Sets the name of the log to write to.
	logName := "songbeamer-helper"
	gcpLogger := client.Logger(logName)

	// Initialize all loggers with their respective severity levels
	loggers["info"] = gcpLogger.StandardLogger(logging.Info)
	loggers["debug"] = gcpLogger.StandardLogger(logging.Debug)
	loggers["warning"] = gcpLogger.StandardLogger(logging.Warning)
	loggers["error"] = gcpLogger.StandardLogger(logging.Error)
	loggers["fatal"] = gcpLogger.StandardLogger(logging.Emergency)

	usingGCP = true
	loggers["info"].Println("GCP logging initialized successfully")
	return true
}

// initLocalLogging sets up local file-based logging as a fallback
func initLocalLogging() {
	log.Println("Initializing local logging...")

	// Create simple loggers that write to stderr
	for _, level := range []string{"info", "debug", "warning", "error", "fatal"} {
		loggers[level] = log.New(os.Stderr, "", 0)
	}

	log.Println("Local logging initialized")
}

// getLogger retrieves the appropriate logger for a given level
func getLogger(level string) *log.Logger {
	if l, ok := loggers[level]; ok {
		return l
	}
	return loggers["info"] // fallback
}

// logWithLevel is a helper function that handles logging with level checks and formatting
func logWithLevel(loggerName, prefix string, minLevel int, useColor bool, colorAttr color.Attribute, format string, args ...any) {
	lock.Lock()
	defer lock.Unlock()

	// Always log to remote logger (GCP or local file)
	logger := getLogger(loggerName)

	// For GCP debug logs, include file/line info
	if usingGCP && loggerName == "debug" {
		_, filename, line, _ := runtime.Caller(2)
		msg := fmt.Sprintf("[%s][%d] %s", filename, line, format)
		logger.Printf(msg, args...)
	} else {
		logger.Printf(format, args...)
	}

	// Check if we should also log to console
	if level < minLevel {
		return
	}

	// Format and log to console for visibility
	msg := fmt.Sprintf("%s: %s", prefix, format)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	if useColor {
		color.Set(colorAttr)
		log.Println(msg)
		color.Unset()
	} else {
		log.Println(msg)
	}
}

// SetLevel sets the global log level.
func SetLevel(l int) {
	lock.Lock()
	defer lock.Unlock()

	level = l
}

func Printf(format string, args ...any) {
	lock.Lock()
	defer lock.Unlock()

	msg := fmt.Sprintf(format, args...)
	getLogger("info").Print(msg)

	if usingGCP {
		fmt.Print(msg)
	}
}

func Println(args ...any) {
	lock.Lock()
	defer lock.Unlock()

	msg := fmt.Sprintln(args...)
	getLogger("info").Print(msg)

	if usingGCP {
		fmt.Print(msg)
	}
}

func Print(args ...any) {
	lock.Lock()
	defer lock.Unlock()

	msg := fmt.Sprint(args...)
	getLogger("info").Print(msg)

	if usingGCP {
		fmt.Print(msg)
	}
}

// Infof logs an info message.
func Infof(format string, args ...any) {
	logWithLevel("info", "INFO", Info, false, 0, format, args...)
}

// Debugf logs a debug message.
func Debugf(format string, args ...any) {
	logWithLevel("debug", "DEBUG", Debug, true, color.FgBlue, format, args...)
}

// Warnf logs a warning message.
func Warnf(format string, args ...any) {
	logWithLevel("warning", "WARNING", Warning, true, color.FgHiYellow, format, args...)
}

// Errorf logs an error message.
func Errorf(format string, args ...any) {
	logWithLevel("error", "ERROR", Error, true, color.FgRed, format, args...)
}

// Fatalf logs a fatal message and exits.
func Fatalf(format string, args ...any) {
	lock.Lock()
	defer lock.Unlock()

	getLogger("fatal").Printf(format, args...)

	msg := fmt.Sprintf("FATAL: %v", format)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	Finalize()
	log.Fatal(msg)
}

// Fatal logs a fatal message and exits.
func Fatal(v ...any) {
	lock.Lock()
	defer lock.Unlock()

	msg := fmt.Sprint(v...)
	getLogger("fatal").Print(msg)

	Finalize()
	log.Fatal(msg)
}

// Finalize makes sure all logs are sent to GCP
func Finalize() {
	if !usingGCP || client == nil {
		return
	}
	err := client.Close()
	if err != nil {
		log.Printf("Error closing GCP logging client: %v", err)
	}
}
