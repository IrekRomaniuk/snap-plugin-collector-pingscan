/*
http://www.apache.org/licenses/LICENSE-2.0.txt

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

package main

import (
	"fmt"
	"os"

	"github.com/IrekRomaniuk/snap-plugin-collector-pingcount/pingcount"
	"github.com/intelsdi-x/snap/control/plugin"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Fprintf(os.Stderr, "Plugin must be run as root\n")
		os.Exit(1)
	}
	plugin.Start(pingcount.Meta(), pingcount.New(), os.Args[1])
}
