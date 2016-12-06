/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

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
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/IrekRomaniuk/snap-plugin-collector-pingscan/pingscan/targets"
	"fmt"
)

const (
	target = "./examples/pinglist.txt"
)

func TestPingscanPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, pluginName )
		So(meta.Version, ShouldResemble, pluginVersion)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})
	Convey("Create Pingscan Collector", t, func() {
		collector := New()
		Convey("So Pingscan collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})
		Convey("So Pingscan collector should be of Pingscan type", func() {
			So(collector, ShouldHaveSameTypeAs, &PingscanCollector{})
		})
		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := collector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
				t.Log(configPolicy)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			Convey("So config policy namespace should be /niuk/pingscan", func() {
				conf := configPolicy.Get([]string{"niuk", "pingscan"})
				So(conf, ShouldNotBeNil)
				So(conf.HasRules(), ShouldBeTrue)
				tables := conf.RulesAsTable()
				So(len(tables), ShouldEqual, 1)
				for _, rule := range tables {
					So(rule.Name, ShouldEqual, "target") //ShouldBeIn if more rules
					switch rule.Name {
					case "target":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					}
				}
			})
		})
	})
}

func TestReadTargets(t *testing.T) {
	Convey("Read pinglist.txt from examples ", t, func() {
	target := "../examples/pinglist.txt"
		hosts, _ := ReadTargets(target)
		Convey("So pinglist.txt should contain 3 items", func() {
			So(len(hosts), ShouldEqual,1100)
		})
	})
}

func TestPingscanCollector_CollectMetrics(t *testing.T) {
	cfg := setupCfg("../examples/pinglist.txt")
	Convey("Pingscan collector", t, func() {
		p := New()
		mt, err := p.GetMetricTypes(cfg)
		if err != nil {
			t.Fatal(err)
		}
		So(len(mt), ShouldEqual, 1)
		Convey("collect metrics", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"niuk", "pingscan", "total-up"),
					Config_: cfg.ConfigDataNode,
				},
			}
			//fmt.Println(mts[0].Config().Table())
			metrics, err := p.CollectMetrics(mts)
			So(err, ShouldBeNil)
			So(metrics, ShouldNotBeNil)
			So(len(metrics), ShouldEqual, 1)
			So(metrics[0].Namespace()[0].Value, ShouldEqual, "niuk")
			So(metrics[0].Namespace()[1].Value, ShouldEqual, "pingscan")
			for _, m := range metrics {
				fmt.Println(m.Namespace()[2].Value,m.Data())
				So(m.Namespace()[2].Value, ShouldEqual, "total-up")
				So(m.Data(), ShouldEqual, 2) //Assuming 8.8.8.8 and 4.2.2.2 are reachable
				t.Log(m.Namespace()[2].Value, m.Data())
			}
		})
	})
}

func setupCfg(target string) plugin.ConfigType {
	node := cdata.NewNode()
	node.AddItem("target", ctypes.ConfigValueStr{Value: target})
	return plugin.ConfigType{ConfigDataNode: node}
}

