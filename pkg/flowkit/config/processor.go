/*
 * Flow CLI
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"encoding/json"
	"fmt"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
)

// ProcessorRun all pre-processors.
func ProcessorRun(raw []byte) ([]byte, map[string]string) {
	raw = processEnv(raw)
	raw, accountFromFiles := processFromFile(raw)

	return raw, accountFromFiles
}

// processEnv finds env variables and insert env values.
func processEnv(raw []byte) []byte {
	_ = godotenv.Load() // try to load .env file

	raw, _ = envsubst.Bytes(raw)
	return raw
}

// processFromFile finds file variables and insert content.
func processFromFile(raw []byte) ([]byte, map[string]string) {
	accountFromFiles := map[string]string{}

	type config struct {
		Accounts    map[string]map[string]any `json:"accounts,omitempty"`
		Contracts   any                       `json:"contracts,omitempty"`
		Networks    any                       `json:"networks,omitempty"`
		Deployments any                       `json:"deployments,omitempty"`
		Emulators   any                       `json:"emulators,omitempty"`
	}

	var conf config
	_ = json.Unmarshal(raw, &conf)

	for name, val := range conf.Accounts {
		if location := val["fromFile"]; location != nil {
			fmt.Println("DEPRECATED: the fromFile directive will be soon deprecated, please use key file type instead: https://developers.flow.com/tools/flow-cli/configuration#advanced-format-1")
			accountFromFiles[name] = location.(string)
			delete(conf.Accounts, name)
		}
	}

	raw, _ = json.Marshal(conf)
	return raw, accountFromFiles
}
