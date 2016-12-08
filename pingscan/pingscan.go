/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

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

package pingscan

import (
	"errors"
	"strings"
	"time"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	"github.com/IrekRomaniuk/snap-plugin-collector-pingscan/pingscan/scan"
)
const (
	vendor        = "niuk"
	fs            = "pingscan"
	pluginName    = "pingscan"
	pluginVersion = 1
	pluginType    = plugin.CollectorPluginType
)
var (
	metricNames = []string{
		"total-up",
	}
)
type PingscanCollector struct {
}

func New() *PingscanCollector {
	pingscan := &PingscanCollector{}
	return pingscan
}

/*  CollectMetrics collects metrics for testing.

CollectMetrics() will be called by Snap when a task that collects one of the metrics returned from this plugins
GetMetricTypes() is started. The input will include a slice of all the metric types being collected.

The output is the collected metrics as plugin.Metric and an error.
*/
func (pingscan *PingscanCollector) CollectMetrics(mts []plugin.MetricType) (metrics []plugin.MetricType, err error) {

	hostsStr := mts[0].Config().Table()["hosts"].(ctypes.ConfigValueStr).Value

	hosts := strings.Fields(hostsStr)

	if len(hosts) == 0 {
		return nil, errors.New("No host requested to ping")
	}

	for _, mt := range mts {
		ns := mt.Namespace()

		val := scan.Ping(hosts)
		/*if err != nil {
			return nil, fmt.Errorf("Error collecting metrics: %v", err)
		}*/
		//fmt.Println(val)
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
func (pingscan *PingscanCollector) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	for _, metricName := range metricNames {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace("niuk", "pingscan", metricName),
		})
	}
	return mts, nil
}


// GetConfigPolicy returns plugin configuration
func (pingscan *PingscanCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule0, _ := cpolicy.NewStringRule("hosts", true)
	cp := cpolicy.NewPolicyNode()
	cp.Add(rule0)
	c.Add([]string{"niuk", "pingscan"}, cp)
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

