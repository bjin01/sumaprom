package collector

import (
	"log"
	"sumaprom/auth"
	"sumaprom/request"
	"sumaprom/schedules"
	"sumaprom/system"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (cc *PkgsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc *PkgsCollector) Collect(ch chan<- prometheus.Metric) {
	MysumaLogin := auth.Sumalogin{Login: cc.Sumainfo.Userid, Passwd: cc.Sumainfo.Password}
	request.Sumahost = &cc.Sumainfo.Server
	SessionKey, err := auth.Login("auth.login", MysumaLogin)
	if err != nil {
		log.Fatal(err)
	}

	listactivesystems := new(system.ListActiveSystem)
	time.Sleep(1 * time.Second)
	_ = listactivesystems.GetActiveSystems(&SessionKey)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	log.Output(2, "listactivesystems.GetActiveSystems Done")

	time.Sleep(1 * time.Second)
	err = listactivesystems.GetUpgPkgs(&SessionKey)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	log.Output(2, "listactivesystems.Getpackages Done")

	time.Sleep(1 * time.Second)
	err = listactivesystems.Getpatches(&SessionKey)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	log.Output(2, "listactivesystems.Getpatches Done")

	cc.PatchPgksInfo = listactivesystems
	cc.GetPkgsResult(ch)
	log.Output(2, "GetPkgResult Metrics Done")
	cc.GetPatchResult(ch)
	log.Output(2, "GetPkgResult Metrics Done")

	listjobs := new(schedules.ListJobs)
	_ = listjobs.GetCompletedJobs(&SessionKey)
	_ = listjobs.GetFailedJobs(&SessionKey)
	_ = listjobs.GetPendingjobs(&SessionKey)

	log.Output(2, "listjobs.GetCompletedJobs Done")
	cc.Completed = &listjobs.Completed

	cc.Failed = &listjobs.Failed
	log.Output(2, "listjobs.GetFailedJobs Done")

	err = auth.Logout("auth.logout", SessionKey)
	if err != nil {
		log.Fatal(err)
	}

	job_status := []string{"completed", "failed"}

	for _, j := range job_status {
		cc.GetJobsResult(j, ch)
	}
	log.Output(2, "GetJobsResult Metrics Done.")
}
