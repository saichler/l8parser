/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package tests provides test utilities and test cases for the L8Parser.
// It includes unit tests for parsing rules, integration tests for the parsing service,
// and helper functions for setting up test topologies.
package tests

import (
	"github.com/saichler/l8bus/go/overlay/protocol"
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/l8test/go/infra/t_topology"
	. "github.com/saichler/l8types/go/ifs"
)

// topo is the shared test topology instance used across all tests.
var topo *TestTopology

// init sets the default log level for tests to Trace for detailed debugging.
func init() {
	Log.SetLogLevel(Trace_Level)
}

// setup initializes the test environment by setting up the test topology.
func setup() {
	setupTopology()
}

// tear cleans up the test environment by shutting down the test topology.
func tear() {
	shutdownTopology()
}

// reset logs the end of a test and resets all handlers in the topology.
func reset(name string) {
	Log.Info("*** ", name, " end ***")
	topo.ResetHandlers()
}

// setupTopology creates a new test topology with 4 nodes on ports 20000, 30000, 40000.
func setupTopology() {
	protocol.MessageLog = true
	topo = NewTestTopology(4, []int{20000, 30000, 40000}, Info_Level)
}

// shutdownTopology gracefully shuts down the test topology.
func shutdownTopology() {
	topo.Shutdown()
}
