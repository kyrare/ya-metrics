package metrics

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// GetPstil получение pstil метрик
func GetPstil() map[string]float64 {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, false)

	var cpuUtilization float64
	if len(c) > 0 {
		cpuUtilization = c[0]
	}

	return map[string]float64{
		"TotalMemory":     float64(v.Total),
		"FreeMemory":      float64(v.Free),
		"CPUutilization1": cpuUtilization,
	}
}
