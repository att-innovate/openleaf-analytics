// The MIT License (MIT)
//
// Copyright (c) 2017 AT&T
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package inputs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Ports struct {
	Objects []struct {
		ObjectId string `json:"ObjectId"`
		Object   struct {
			IntfRef       string `json:"IntfRef"`
			IfIndex       int64  `json:"IfIndex"`
			OperState     string `json:"OperState"`
			IfInDiscards  int64  `json:"IfInDiscards"`
			IfOutDiscards int64  `json:"IfOutDiscards"`
			IfEtherPkts   int64  `json:"IfEtherPkts"`
			IfEtherMCPkts int64  `json:"IfEtherMCPkts"`
			IfEtherBCPkts int64  `json:"IfEtherBcastPkts"`
			IfInOctets    int64  `json:"IfInOctets"`
			IfOutOctets   int64  `json:"IfOutOctets"`
		} `json:"Object"`
	} `json:"Objects"`
}

type Status struct {
	ObjectId string `json:"ObjectId"`
	Object   struct {
		Ready bool `json:"Ready"`
	} `json:"Object"`
}

type SnapRoute struct {
	Url               string
	lastTime          time.Time
	lastInDiscards    [512]int64
	lastOutDiscards   [512]int64
	lastEtherPkts     [512]int64
	lastEtherMCPkts   [512]int64
	lastEtherBCPkts   [512]int64
	lastInOctets      [512]int64
	lastOutOctets     [512]int64
}

var sampleConfig = `
  ## URL-prefix for SnapRoute
  url = "http://localhost:8080/public/v1/"
`

func (_ *SnapRoute) Description() string {
	return `Read SnapRoute Metrics`
}

func (_ *SnapRoute) SampleConfig() string {
	return sampleConfig
}

func (s *SnapRoute) Gather(acc telegraf.Accumulator) error {
	defer func() {
		if r := recover(); r != nil {
			glog.Error("E! Problem reading from SnapRoute: ", r)
			acc.AddFields("status", map[string]interface{}{"ready": false}, nil, time.Now())
		}
	}()

	now := time.Now()
	diffTime := now.Sub(s.lastTime).Seconds()

	var ports Ports
	var requestURL = fmt.Sprint(s.Url, "state/Ports")
	content, err := getContent(requestURL)
	if err != nil {
		glog.Error("Error talking to SnapRoute:", err)
		acc.AddFields("status", map[string]interface{}{"ready": false}, nil, time.Now())
		return err
	}

	err = json.Unmarshal(content, &ports)
	if err != nil {
		glog.Error("Error Umarshalling:", err)
		glog.Error("content:", content)
		return err
	}

	for _, port := range ports.Objects {
		if port.Object.OperState == "UP" {
			tags := map[string]string{
				"port": port.Object.IntfRef,
			}

			var (
				metricsCount       = 0
				inDiscards   int64 = 0
				outDiscards  int64 = 0
				etherPkts    int64 = 0
				etherMCPkts  int64 = 0
				etherBCPkts  int64 = 0
				inOctets     int64 = 0
				outOctets    int64 = 0
			)

			if port.Object.IfInDiscards == 0 {
				metricsCount++
			} else if s.lastInDiscards[port.Object.IfIndex] != 0 {
				inDiscards = (port.Object.IfInDiscards - s.lastInDiscards[port.Object.IfIndex])
				metricsCount++
			}
			s.lastInDiscards[port.Object.IfIndex] = port.Object.IfInDiscards

			if port.Object.IfOutDiscards == 0 {
				metricsCount++
			} else if s.lastOutDiscards[port.Object.IfIndex] != 0 {
				outDiscards = (port.Object.IfOutDiscards - s.lastOutDiscards[port.Object.IfIndex])
				metricsCount++
			}
			s.lastOutDiscards[port.Object.IfIndex] = port.Object.IfOutDiscards

			if port.Object.IfEtherPkts == 0 {
				metricsCount++
			} else if s.lastEtherPkts[port.Object.IfIndex] != 0 {
				etherPkts = (port.Object.IfEtherPkts - s.lastEtherPkts[port.Object.IfIndex]) / int64(diffTime)
				metricsCount++
			}
			s.lastEtherPkts[port.Object.IfIndex] = port.Object.IfEtherPkts

			if port.Object.IfEtherMCPkts == 0 {
				metricsCount++
			} else if s.lastEtherMCPkts[port.Object.IfIndex] != 0 {
				etherMCPkts = (port.Object.IfEtherMCPkts - s.lastEtherMCPkts[port.Object.IfIndex]) / int64(diffTime)
				metricsCount++
			}
			s.lastEtherMCPkts[port.Object.IfIndex] = port.Object.IfEtherMCPkts

			if port.Object.IfEtherBCPkts == 0 {
				metricsCount++
				acc.AddGauge("ports", map[string]interface{}{"sent_bc": 0}, tags, now)
			} else if s.lastEtherBCPkts[port.Object.IfIndex] != 0 {
				metricsCount++
				etherBCPkts = (port.Object.IfEtherBCPkts - s.lastEtherBCPkts[port.Object.IfIndex]) / int64(diffTime)
			}
			s.lastEtherBCPkts[port.Object.IfIndex] = port.Object.IfEtherBCPkts

			if port.Object.IfInOctets == 0 {
				metricsCount++
				acc.AddGauge("ports", map[string]interface{}{"in_octets": 0}, tags, now)
			} else if s.lastInOctets[port.Object.IfIndex] != 0 {
				metricsCount++
				inOctets = (port.Object.IfInOctets - s.lastInOctets[port.Object.IfIndex]) / int64(diffTime)
			}
			s.lastInOctets[port.Object.IfIndex] = port.Object.IfInOctets

			if port.Object.IfOutOctets == 0 {
				metricsCount++
				acc.AddGauge("ports", map[string]interface{}{"out_octets": 0}, tags, now)
			} else if s.lastOutOctets[port.Object.IfIndex] != 0 {
				metricsCount++
				outOctets = (port.Object.IfOutOctets - s.lastOutOctets[port.Object.IfIndex]) / int64(diffTime)
			}
			s.lastOutOctets[port.Object.IfIndex] = port.Object.IfOutOctets

			if metricsCount == 7 {
				acc.AddCounter("ports", map[string]interface{}{"discard_in": inDiscards}, tags, now)
				acc.AddCounter("ports", map[string]interface{}{"discard_out": outDiscards}, tags, now)
				acc.AddGauge("ports", map[string]interface{}{"sent": etherPkts}, tags, now)
				acc.AddGauge("ports", map[string]interface{}{"sent_mc": etherMCPkts}, tags, now)
				acc.AddGauge("ports", map[string]interface{}{"sent_bc": etherBCPkts}, tags, now)
				acc.AddGauge("ports", map[string]interface{}{"in_octets": inOctets}, tags, now)
				acc.AddGauge("ports", map[string]interface{}{"out_octets": outOctets}, tags, now)
			}
		}
	}

	s.lastTime = now

	var status Status
	requestURL = fmt.Sprint(s.Url, "state/SystemStatus")
	content, err = getContent(requestURL)
	if err != nil {
		glog.Error("Error talking to SnapRoute:", err)
		acc.AddFields("status", map[string]interface{}{"ready": false}, nil, time.Now())
		return err
	}

	err = json.Unmarshal(content, &status)
	if err != nil {
		glog.Error("Error Unmarshalling:", err)
		glog.Error("content:", content)
		return err
	}

	acc.AddFields("status", map[string]interface{}{"ready": status.Object.Ready}, nil, now)

	return nil
}

func getContent(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func init() {
	inputs.Add("snaproute", func() telegraf.Input {
		return &SnapRoute{}
	})
}
