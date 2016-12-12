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

package pingcount

import (
	"fmt"
	"time"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	"github.com/IrekRomaniuk/snap-plugin-collector-pingcount/pingcount/targets"
	"github.com/IrekRomaniuk/snap-plugin-collector-pingcount/pingcount/scan"
)
const (
	vendor        = "niuk"
	fs            = "pingcount"
	pluginName    = "pingcount"
	pluginVersion = 1
	pluginType    = plugin.CollectorPluginType
)
var (
	metricNames = []string{
		"total-up",
	}
)
type PingcountCollector struct {
}

func New() *PingcountCollector {
	pingcount := &PingcountCollector{}
	return pingcount
}

/*  CollectMetrics collects metrics for testing.

CollectMetrics() will be called by Snap when a task that collects one of the metrics returned from this plugins
GetMetricTypes() is started. The input will include a slice of all the metric types being collected.

The output is the collected metrics as plugin.Metric and an error.
*/
func (pingcount *PingcountCollector) CollectMetrics(mts []plugin.MetricType) (metrics []plugin.MetricType, err error) {
	var (
		target string
	)
	conf := mts[0].Config().Table()
	targetConf, ok := conf["target"]
	if !ok || targetConf.(ctypes.ConfigValueStr).Value == "" {
		return nil, fmt.Errorf("target missing from config, %v", conf)
	} else {
		target = targetConf.(ctypes.ConfigValueStr).Value
	}

	hosts, err := targets.ReadTargets(target)
	if err != nil { return nil, fmt.Errorf("Error reading target: %v", err) }

	if len(hosts) == 0 { return nil, fmt.Errorf("No host defined in file %v", target)}

	for _, mt := range mts {
		ns := mt.Namespace()

		val := scan.Ping(hosts)

		metric := plugin.MetricType{
			Namespace_: ns,
			Data_:      val,
			Timestamp_: time.Now(),
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}



/*
	GetMetricTypes returns metric types for testing.
	GetMetricTypes() will be called when your plugin is loaded in order to populate the metric catalog(where snaps stores all
	available metrics).

	Config info is passed in. This config information would come from global config snap settings.

	The metrics returned will be advertised to users who list all the metrics and will become targetable by tasks.
*/
func (pingcount *PingcountCollector) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	for _, metricName := range metricNames {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace("niuk", "pingcount", metricName),
		})
	}
	return mts, nil
}


// GetConfigPolicy returns plugin configuration
func (pingcount *PingcountCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule0, _ := cpolicy.NewStringRule("target", true)
	cp := cpolicy.NewPolicyNode()
	cp.Add(rule0)
	c.Add([]string{"niuk", "pingcount"}, cp)
	return c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		pluginName,
		pluginVersion,
		pluginType,
		[]string{plugin.SnapGOBContentType},//[]string{},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.ConcurrencyCount(1),
	)
}

