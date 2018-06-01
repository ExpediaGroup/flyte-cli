/*
Copyright (C) 2018 Expedia Group.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/spf13/cobra"
	"bytes"
	"github.com/stretchr/testify/assert"
)

func TestTestCommand_ShouldExecuteStepAndReturnOutputForJsonInput(t *testing.T) {
	output, err := executeCommand(rootCmd, "test", "testdata/step-test.json")
	require.NoError(t, err)

	assert.Equal(t, jsonOutput, output)
}

func TestTestCommand_ShouldExecuteStepAndReturnOutputForYamlInput(t *testing.T) {
	output, err := executeCommand(rootCmd, "test", "testdata/step-test.yaml")
	require.NoError(t, err)

	assert.Equal(t, jsonOutput, output)
}

func TestTestCommand_ShouldExecuteStepAndReturnOutputForYmlInput(t *testing.T) {
	output, err := executeCommand(rootCmd, "test", "testdata/step-test.yml")
	require.NoError(t, err)

	assert.Equal(t, jsonOutput, output)
}

func TestTestCommand_ShouldExecuteStepAndReturnOutputAsYaml(t *testing.T) {
	format = "yaml"
	output, err := executeCommand(rootCmd, "test", "testdata/step-test.json")
	require.NoError(t, err)

	assert.Equal(t, yamlOutput, output)
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	_, err = root.ExecuteC()
	return buf.String(), err
}

const jsonOutput = `{
	"name": "SendMessage",
	"packName": "Slack",
	"input": {
		"channelId": "123",
		"message": "Hey \u003c@johnny\u003e, I'm up and running :run:"
	},
	"context": {
		"ChannelID": "123",
		"UserID": "johnny"
	}
}`

const yamlOutput = `context:
  ChannelID: "123"
  UserID: johnny
input:
  channelId: "123"
  message: 'Hey <@johnny>, I''m up and running :run:'
name: SendMessage
packName: Slack
`