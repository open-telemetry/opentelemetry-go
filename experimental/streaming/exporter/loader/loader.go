// Copyright 2019, OpenTelemetry Authors
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

package loader

import (
	"fmt"
	"os"
	"plugin"
	"time"

	"go.opentelemetry.io/experimental/streaming/exporter/observer"
)

// TODO add buffer support directly, eliminate stdout

func init() {
	pluginName := os.Getenv("OPENTELEMETRY_LIB")
	if pluginName == "" {
		return
	}
	sharedObj, err := plugin.Open(pluginName)
	if err != nil {
		fmt.Println("Open failed", pluginName, err)
		return
	}

	obsPlugin, err := sharedObj.Lookup("Observer")
	if err != nil {
		fmt.Println("Observer not found", pluginName, err)
		return
	}

	f, ok := obsPlugin.(func() observer.Observer)
	if !ok {
		fmt.Printf("Observer not valid\n")
		return
	}
	observer.RegisterObserver(f())
}

func Flush() {
	// TODO implement for exporter/{stdout,stderr,buffer}
	time.Sleep(1 * time.Second)
}
