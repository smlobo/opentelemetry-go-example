// Copyright 2022 Sheldon Lobo
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

package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strconv"
)

var Config map[string]string
var Sleep int

func SetupConfig() {
	// Open the config file
	configFile, err := os.Open("configs/opentelemetry-go-example.yml")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		os.Exit(1)
	}
	defer configFile.Close()

	// Unmarshal into the map
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error unmarshalling config file:", err)
		os.Exit(1)
	}

	// Override with environment variables
	for variable, _ := range Config {
		envVarSet := os.Getenv(variable)
		if len(envVarSet) != 0 {
			Config[variable] = envVarSet
			log.Println("From environment: [", variable, "] =", envVarSet)
		}
	}

	// Set sleep time as an int
	Sleep, _ = strconv.Atoi(Config["SLEEP"])
}