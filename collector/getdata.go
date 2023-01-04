package collector

import (
	"log"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (x *PkgsCollector) GetPkgsResult(ch chan<- prometheus.Metric) {

	for _, s := range x.PatchPgksInfo.Result {

		ch <- prometheus.MustNewConstMetric(
			pkgsDesc,
			prometheus.GaugeValue,
			float64(float64(len(s.PackageList.Result))),
			s.Name,
			"updates",
		)
	}
}

func (x *PkgsCollector) GetPatchResult(ch chan<- prometheus.Metric) {

	for _, s := range x.PatchPgksInfo.Result {
		counter_security := 0
		counter_buxfix := 0
		counter_enhancements := 0
		counter_others := 0
		counter_livepatches := 0
		for _, p := range s.Patches.Result {

			if p.Advisory_type == "Security Advisory" {
				counter_security++
				if strings.Contains(p.Advisory_synopsis, "Linux Kernel (Live Patch") {
					//fmt.Printf("livepatch: %s %s\n", s.Name, p.Advisory_synopsis)

					counter_livepatches++
				}
			} else if p.Advisory_type == "Bug Fix Advisory" {
				counter_buxfix++
			} else if p.Advisory_type == "Product Enhancement Advisory" {
				counter_enhancements++
			} else {
				counter_others++
			}
		}

		if counter_livepatches > 0 {
			log.Printf("Found %d available Live Patches for %s\n", counter_livepatches, s.Name)
		}

		ch <- prometheus.MustNewConstMetric(
			errataDesc,
			prometheus.GaugeValue,
			float64(counter_livepatches),
			s.Name,
			"livepatches",
		)

		ch <- prometheus.MustNewConstMetric(
			errataDesc,
			prometheus.GaugeValue,
			float64(counter_security),
			s.Name,
			"security",
		)

		ch <- prometheus.MustNewConstMetric(
			errataDesc,
			prometheus.GaugeValue,
			float64(counter_buxfix),
			s.Name,
			"bugfix",
		)

		ch <- prometheus.MustNewConstMetric(
			errataDesc,
			prometheus.GaugeValue,
			float64(counter_enhancements),
			s.Name,
			"enhancements",
		)

		ch <- prometheus.MustNewConstMetric(
			errataDesc,
			prometheus.GaugeValue,
			float64(counter_others),
			s.Name,
			"others",
		)
	}
}

func (x *PkgsCollector) GetJobsResult(a string, ch chan<- prometheus.Metric) {
	patch_int := 0
	pkg_refresh_int := 0
	pkg_install_int := 0
	search_pkg_refresh := "Package List Refresh"
	search_patch_updates := "Patch Update"
	search_pkg_installs := "Package Install"
	layout := "Jan 2, 2006, 3:04:05 PM"
	timenow := time.Now().Local()

	if a == "completed" {
		for _, b := range x.Completed.Result {
			mytime, _ := time.Parse(layout, b.Earliest.String())

			//fmt.Printf("%s: %s Jobs at %s today: %s\n", b.Name, b.Type, mytime, time.Now().Local())
			if strings.Contains(b.Type, search_patch_updates) && mytime.YearDay() == timenow.YearDay() {
				patch_int++
			}

			if strings.Contains(b.Type, search_pkg_installs) && mytime.YearDay() == timenow.YearDay() {
				pkg_install_int++
			}

			if strings.Contains(b.Type, search_pkg_refresh) && mytime.YearDay() == timenow.YearDay() {
				pkg_refresh_int++
			}
		}
	}

	if a == "failed" {
		for _, b := range x.Failed.Result {
			mytime, _ := time.Parse(layout, b.Earliest.String())

			//fmt.Printf("%s: %s Jobs at %s today: %s\n", b.Name, b.Type, mytime, time.Now().Local())
			if strings.Contains(b.Type, search_patch_updates) && mytime.YearDay() == timenow.YearDay() {
				patch_int++
			}

			if strings.Contains(b.Type, search_pkg_installs) && mytime.YearDay() == timenow.YearDay() {
				pkg_install_int++
			}

			if strings.Contains(b.Type, search_pkg_refresh) && mytime.YearDay() == timenow.YearDay() {
				pkg_refresh_int++
			}
		}
	}

	ch <- prometheus.MustNewConstMetric(
		sumaJobsDesc,
		prometheus.GaugeValue,
		float64(pkg_refresh_int),
		"pkg_refresh",
		a,
	)

	ch <- prometheus.MustNewConstMetric(
		sumaJobsDesc,
		prometheus.GaugeValue,
		float64(pkg_install_int),
		"pkg_install",
		a,
	)

	ch <- prometheus.MustNewConstMetric(
		sumaJobsDesc,
		prometheus.GaugeValue,
		float64(patch_int),
		"patch_updates",
		a,
	)

}
