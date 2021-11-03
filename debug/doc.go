// Copyright The OpenTelemetry Authors
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

/*
Package debug provides methods to control internal logging of opentelemetry.

# Usage:
The default logger will discard all messages, to see the logs the application
developer needs to set a logger.

```
import (
	"os"

	"go.opentelemetry.io/otel/debug"
	"github.com/go-logr/stdr"
)

func main() {
	debug.SetLogger(
		stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))
	)

...
}
```

# Verbosity:
There are two different levels of logging currently the default (0) and
additional (5).  To see the additional debug logs you need to configure your
logger to display the level 5 logs.

For example when using stdr use `stdr.SetVerbosity(5)` to see the aditional logs
*/
package debug // import "go.opentelemetry.io/otel/debug"
