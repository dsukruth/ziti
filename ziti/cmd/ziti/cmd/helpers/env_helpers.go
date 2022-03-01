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

package helpers

import (
	"github.com/openziti/ziti/ziti/cmd/ziti/constants"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	h := os.Getenv("USERPROFILE") // windows
	if h == "" {
		h = "."
	}
	return NormalizePath(h)
}

func WorkingDir() (string, error) {
	wd, err := os.Getwd()
	if wd == "" || err != nil {
		return "", err
	}

	return NormalizePath(wd), nil
}

func GetZitiHome() (string, error) {

	// Get path from env variable
	retVal := os.Getenv(constants.ZitiHomeVarName)

	if retVal == "" {
		// If not set, create a default path of the current working directory
		workingDir, err := WorkingDir()
		if err != nil {
			return "", err
		}

		err = os.Setenv(constants.ZitiHomeVarName, workingDir)
		if err != nil {
			return "", err
		}

		retVal = os.Getenv(constants.ZitiHomeVarName)
	}

	return NormalizePath(retVal), nil
}

func GetZitiCtrlAdvertisedAddress() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		err := errors.Wrap(err, "Unable to get hostname")
		if err != nil {
			return "", err
		}
	}

	return getValueOrSetAndGetDefault(constants.ZitiCtrlAdvertisedAddressVarName, hostname, false)
}

func GetZitiCtrlPort() (string, error) {
	return getValueOrSetAndGetDefault(constants.ZitiCtrlPortVarName, constants.DefaultZitiControllerPort, false)
}

func GetZitiCtrlListenerAddress() (string, error) {
	return getValueOrSetAndGetDefault(constants.ZitiCtrlListenerAddressVarName, constants.DefaultZitiControllerListenerAddress, false)
}

func GetZitiCtrlName() (string, error) {
	return getValueOrSetAndGetDefault(constants.ZitiCtrlNameVarName, constants.DefaultZitiControllerName, false)
}

func GetZitiEdgeRouterPort() (string, error) {
	return getValueOrSetAndGetDefault(constants.ZitiEdgeRouterPortVarName, constants.DefaultZitiEdgeRouterPort, false)
}

func GetZitiEdgeCtrlListenerHostPort() (string, error) {
	return getValueOrSetAndGetDefault(constants.ZitiEdgeCtrlListenerHostPortVarName, constants.DefaultZitiEdgeListenerHostPort, false)
}

func GetZitiEdgeCtrlAdvertisedHostPort() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		err := errors.Wrap(err, "Unable to get hostname")
		if err != nil {
			return "", err
		}
	}

	return getValueOrSetAndGetDefault(constants.ZitiEdgeCtrlAdvertisedHostPortVarName, hostname+":"+constants.DefaultZitiEdgeAPIPort, false)
}

func getValueOrSetAndGetDefault(envVarName string, defaultValue string, forceDefault bool) (string, error) {
	retVal := ""
	if !forceDefault {
		// Get path from env variable
		retVal = os.Getenv(envVarName)
		if retVal != "" {
			return retVal, nil
		}
	}

	err := os.Setenv(envVarName, defaultValue)
	if err != nil {
		return "", err
	}

	retVal = os.Getenv(envVarName)

	return retVal, nil
}

// NormalizePath replaces windows \ with / which windows allows for
func NormalizePath(input string) string {
	return strings.ReplaceAll(input, "\\", "/")
}