package process

import (
	"bytes"
	"fmt"
	"grits/types"

	"golang.org/x/exp/slices"
)

// Common environment used both by the typechecker and interpreter/runtime environment

type GlobalEnvironment struct {
	// Contains all the functions definitions
	FunctionDefinitions *[]FunctionDefinition

	// Contains all the type definitions
	Types *[]types.SessionTypeDefinition

	// Logging levels
	LogLevels []LogLevel
}

/////////////////////////////////////////////////////
////////////////////// Logging //////////////////////
/////////////////////////////////////////////////////

// Similar to Println
func (globalEnv *GlobalEnvironment) log(level LogLevel, message string) {
	if slices.Contains(globalEnv.LogLevels, level) {
		var buffer bytes.Buffer

		buffer.WriteString(message)
		buffer.WriteString("\n")

		fmt.Print(buffer.String())
	}
}

// Similar to Printf
func (globalEnv *GlobalEnvironment) logf(level LogLevel, message string, args ...interface{}) {
	if slices.Contains(globalEnv.LogLevels, level) {
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf(message, args...))
		fmt.Print(buffer.String())
	}
}
