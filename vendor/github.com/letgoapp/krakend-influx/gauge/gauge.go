package gauge

import (
	"time"

	"github.com/devopsfaith/krakend/logging"
	"github.com/influxdata/influxdb/client/v2"
)

func Points(hostname string, now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	res := make([]*client.Point, 4)

	in := map[string]interface{}{
		"gauge": int(counters["krakend.router.connected-gauge"]),
	}
	incoming, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "in"}, in, now)
	if err != nil {
		logger.Error("creating incoming connection counters point:", err.Error())
		return res
	}
	res[0] = incoming

	out := map[string]interface{}{
		"gauge": int(counters["krakend.router.disconnected-gauge"]),
	}
	outgoing, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "out"}, out, now)
	if err != nil {
		logger.Error("creating outgoing connection counters point:", err.Error())
		return res
	}
	res[1] = outgoing

	debug := map[string]interface{}{}
	runtime := map[string]interface{}{}

	for k, v := range counters {
		if k == "krakend.router.connected-gauge" || k == "krakend.router.disconnected-gauge" {
			continue
		}
		if k[:22] == "krakend.service.debug." {
			debug[k[22:]] = int(v)
			continue
		}
		if k[:24] == "krakend.service.runtime." {
			runtime[k[24:]] = int(v)
			continue
		}
		logger.Debug("unknown gauge key:", k)
	}

	debugPoint, err := client.NewPoint("debug", map[string]string{"host": hostname}, debug, now)
	if err != nil {
		logger.Error("creating debug counters point:", err.Error())
		return res
	}
	res[2] = debugPoint

	runtimePoint, err := client.NewPoint("runtime", map[string]string{"host": hostname}, runtime, now)
	if err != nil {
		logger.Error("creating runtime counters point:", err.Error())
		return res
	}
	res[3] = runtimePoint

	return res
}
