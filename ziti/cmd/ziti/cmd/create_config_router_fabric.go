/*
	Copyright NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package cmd

import (
	_ "embed"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/templates"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	createConfigRouterFabricLong = templates.LongDesc(`
		Creates the fabric router config
`)

	createConfigRouterFabricExample = templates.Examples(`
		# Create the fabric router config for a router named my_router
		ziti create config router fabric --routerName my_router
	`)
)

// CreateConfigRouterFabricOptions the options for the fabric command
type CreateConfigRouterFabricOptions struct {
	CreateConfigRouterOptions
}

//go:embed config_templates/router.yml
var routerConfigFabricTemplate string

// NewCmdCreateConfigRouterFabric creates a command object for the "fabric" command
func NewCmdCreateConfigRouterFabric(data *ConfigTemplateValues) *cobra.Command {
	options := &CreateConfigRouterFabricOptions{}

	cmd := &cobra.Command{
		Use:     "fabric",
		Short:   "Create a fabric router config",
		Aliases: []string{"fabric"},
		Long:    createConfigRouterFabricLong,
		Example: createConfigRouterFabricExample,
		PreRun: func(cmd *cobra.Command, args []string) {
			// Setup logging
			var logOut *os.File
			if options.Verbose {
				logrus.SetLevel(logrus.DebugLevel)
				// Only print log to stdout if not printing config to stdout
				if strings.ToLower(options.Output) != "stdout" {
					logOut = os.Stdout
				} else {
					logOut = os.Stderr
				}
				logrus.SetOutput(logOut)
			}

			// Update fabric router specific values with options passed in
			data.Router.Name = options.RouterName
			SetZitiRouterIdentity(&data.Router, data.Router.Name)
			data.Router.IsFabric = true
		},
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.run(data)
			cmdhelper.CheckErr(err)
		},
	}

	options.addCreateFlags(cmd)
	options.addFlags(cmd)

	return cmd
}

func (options *CreateConfigRouterFabricOptions) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&options.RouterName, optionRouterName, "n", "", "name of the router")
	err := cmd.MarkPersistentFlagRequired(optionRouterName)
	if err != nil {
		return
	}
}

// run implements the command
func (options *CreateConfigRouterFabricOptions) run(data *ConfigTemplateValues) error {

	tmpl, err := template.New("fabric-router-config").Parse(routerConfigFabricTemplate)
	if err != nil {
		return err
	}

	var f *os.File
	if strings.ToLower(options.Output) != "stdout" {
		// Check if the path exists, fail if it doesn't
		basePath := filepath.Dir(options.Output) + "/"
		if _, err := os.Stat(filepath.Dir(basePath)); os.IsNotExist(err) {
			logrus.Fatalf("Provided path: [%s] does not exist\n", basePath)
			return err
		}

		f, err = os.Create(options.Output)
		logrus.Debugf("Created output file: %s", options.Output)
		if err != nil {
			return errors.Wrapf(err, "unable to create config file: %s", options.Output)
		}
	} else {
		f = os.Stdout
	}
	defer func() { _ = f.Close() }()

	if err := tmpl.Execute(f, data); err != nil {
		return errors.Wrap(err, "unable to execute template")
	}

	logrus.Debugf("Fabric Router configuration generated successfully and written to: %s", options.Output)

	return nil
}