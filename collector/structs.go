package collector

import (
	"sumaprom/crypt"
	"sumaprom/schedules"
	"sumaprom/system"

	"github.com/prometheus/client_golang/prometheus"
)

/* type Sumaconf struct {
	Server   string `yaml:"server"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
} */

type PkgsCollector struct {
	Sumainfo      crypt.Sumaconf
	Completed     *schedules.ListCompletedJobs
	Failed        *schedules.ListFailedJobs
	PatchPgksInfo *system.ListActiveSystem
}

var (
	sumaJobsDesc = prometheus.NewDesc(
		"suma_list_of_patching_actions",
		"No of upgradable pkgs per host",
		[]string{"jobtype", "jobstatus"}, nil,
	)

	pkgsDesc = prometheus.NewDesc(
		"suma_list_of_upgrades",
		"No of upgradable pkgs per host",
		[]string{"hostname", "patchtype"}, nil,
	)

	errataDesc = prometheus.NewDesc(
		"suma_list_of_patches",
		"No of patches per host",
		[]string{"hostname", "patchtype"}, nil,
	)
)
