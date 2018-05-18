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
	"github.com/spf13/cobra"
	"fmt"
	"github.com/HotelsDotCom/flyte/template"
	"github.com/HotelsDotCom/flyte/execution"
	jsont "github.com/HotelsDotCom/flyte/json"
	"io/ioutil"
	"encoding/json"
	"strings"
	"github.com/ghodss/yaml"
	"errors"
)

var format string

func newTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test file",
		Short: "Test step execution with trigger event and optional context",
		Long:  testCmdLong,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you need to provide a test file")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			output, err := runTestCmd(args[0], format)
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(c.OutOrStdout(), output)
			return err
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "json", "Output format. One of: json|yaml")
	return cmd
}

func runTestCmd(testFilePath, format string) (string, error) {
	var step testStep
	if err := unmarshalFile(testFilePath, &step); err != nil {
		return "", err
	}

	action, err := step.execute()
	if err != nil {
		return "", err
	}

	out, err := marshal(action, format)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

type testStep struct {
	Step      execution.Step
	Event     execution.Event
	Context   map[string]string
	Datastore map[string]interface{}
}

func (t testStep) execute() (*testAction, error) {
	//override default datastore func which is using mongo to get data item
	//use static map instead which can be passed in the input file
	template.AddStaticContextEntry("datastore", datastoreFn(t.Datastore))

	action, err := t.Step.Execute(t.Event, t.Context)
	if err != nil {
		return nil, err
	}
	return &testAction{
		Name:       action.Name,
		PackName:   action.PackName,
		PackLabels: action.PackLabels,
		Input:      action.Input,
		Context:    action.Context,
	}, nil
}

type testAction struct {
	Name       string            `json:"name"`
	PackName   string            `json:"packName"`
	PackLabels map[string]string `json:"packLabels,omitempty"`
	Input      jsont.Json        `json:"input,omitempty"`
	Context    map[string]string `json:"context,omitempty"`
}

func unmarshalFile(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	ext := filename[strings.LastIndex(filename, "."):]
	switch ext {
	case ".json":
		return json.Unmarshal(data, v)
	case ".yaml":
		return yaml.Unmarshal(data, v)
	default:
		return fmt.Errorf("cannot unmarshal: unsuported file %s", ext)
	}
}

func marshal(v interface{}, format string) ([]byte, error) {
	switch format {
	case "yaml":
		return yaml.Marshal(v)
	default:
		return json.MarshalIndent(v, "", "\t")
	}
}

func datastoreFn(datastore map[string]interface{}) func(string) interface{} {
	if datastore == nil {
		datastore = map[string]interface{}{}
	}
	return func(key string) interface{} {
		v, ok := datastore[key]
		if !ok {
			panic(fmt.Errorf("cannot find datastore item key=%s", key))
		}
		return v
	}
}

const testCmdLong = `

Test step execution with provided test file containing step, trigger event and optional 
context and datastore items. It should be in json or yaml format.

Example:
---
step:
  id: status
  event:
    packName: Slack
    name: ReceivedMessage
  command:
    packName: Slack
    name: SendMessage
    input:
      message: 'Hello'
event:
  pack:
    name: Slack
  event: ReceivedMessage
`
