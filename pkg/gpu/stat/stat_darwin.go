//go:build darwin

package stat

import (
	"github.com/naiba/nezha/pkg/gpu"
)

func GetGPUStat() (float64, error) {
	usage, err := gpu.FindUtilization("PerformanceStatistics", "Device Utilization %")
	return float64(usage), err
}
