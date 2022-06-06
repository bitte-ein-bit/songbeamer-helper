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

// Default log level is Error.
var (
	level         = Info
	lock          sync.Mutex
	debugLogger   *log.Logger = nil
	infoLogger    *log.Logger = nil
	warningLogger *log.Logger = nil
	errorLogger   *log.Logger = nil
	fatalLogger   *log.Logger = nil

	client *logging.Client
)

func init() {
	var err error
	ctx := context.Background()

	// decode the base64 encoded credentials
	decodedCreds, err := base64.StdEncoding.DecodeString(credsJSON)
	if err != nil {
		log.Fatalf("Failed to decode credentials: %v", err)
	}

	// parse the JSON credentials
	var creds GoogleCloudCreds
	err = json.Unmarshal(decodedCreds, &creds)
	if err != nil {
		log.Fatalf("Failed to parse credentials: %v", err)
	}

	// Sets your Google Cloud Platform project ID dynamically from credentials.
	projectID := fmt.Sprintf("projects/%s", creds.ProjectID)

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID, option.WithCredentialsJSON(decodedCreds))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// defer client.Close()

	// Sets the name of the log to write to.
	logName := "songbeamer-helper"

	infoLogger = client.Logger(logName).StandardLogger(logging.Info)
	debugLogger = client.Logger(logName).StandardLogger(logging.Debug)
	warningLogger = client.Logger(logName).StandardLogger(logging.Warning)
	errorLogger = client.Logger(logName).StandardLogger(logging.Error)
	fatalLogger = client.Logger(logName).StandardLogger(logging.Emergency)

	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	infoLogger.Println("logging initialized")
	debugLogger.Println("debug test")
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

	infoLogger.Printf(format, args...)
	msg := fmt.Sprintf("%s", format)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	fmt.Print(msg)
}

func Println(args ...any) {
	Print(args, "\r\n")
}

func Print(args ...any) {
	lock.Lock()
	defer lock.Unlock()

	msg := fmt.Sprint(args...)
	infoLogger.Print(msg)
	fmt.Print(msg)
}

// Infof logs an info message.
func Infof(format string, args ...any) {
	lock.Lock()
	defer lock.Unlock()
	infoLogger.Printf(format, args...)
	if level < Info {
		return
	}

	msg := fmt.Sprintf("INFO: %s", format)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	log.Println(msg)
}

// Debugf logs a debug message.
func Debugf(format string, args ...any) {
	lock.Lock()
	defer lock.Unlock()
	_, filename, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf("[%s][%d] %s", filename, line, format)
	debugLogger.Printf(msg, args...)
	if level < Debug {
		return
	}

	msg = fmt.Sprintf("DEBUG: %s", format)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	color.Set(color.FgBlue)
	log.Println(msg)
	color.Unset()
}

// Warnf logs an error message.
func Warnf(format string, args ...any) {
	lock.Lock()
	defer lock.Unlock()
	warningLogger.Printf(format, args...)
	if level < Warning {
		return
	}

	msg := fmt.Sprintf("WARNING: %s", format)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	color.Set(color.FgHiYellow)
	log.Println(msg)
	color.Unset()
}

// Errorf logs an error message.
func Errorf(format string, args ...any) {
	lock.Lock()
	defer lock.Unlock()
	errorLogger.Printf(format, args...)
	if level < Error {
		return
	}

	msg := fmt.Sprintf("ERROR: %s", format)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	color.Set(color.FgRed)
	log.Println(msg)
	color.Unset()
}

// Fatalf logs an error message.
func Fatalf(format string, args ...any) {
	lock.Lock()
	defer lock.Unlock()
	fatalLogger.Printf(format, args...)
	msg := fmt.Sprintf("FATAL: %v", format)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	Finalize()
	log.Fatal(msg)
}

func Fatal(v ...any) {
	lock.Lock()
	defer lock.Unlock()
	msg := fmt.Sprint(v...)
	fatalLogger.Print(msg)
	Finalize()
	log.Fatal(msg)
}

// Finalize makes sure all logs are sent to GCP
func Finalize() {
	err := client.Close()
	if err != nil {
		log.Fatalf("Fehler beim Hochladen der Logs: %v", err)
	}
}
