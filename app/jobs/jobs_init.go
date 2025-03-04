package jobs

import (
	"github.com/mt1976/frantic-aliquid/app/business/domains"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/robfig/cron/v3"
)

var scheduledTasks *cron.Cron
var cfg *commonConfig.Settings
var domain = domains.JOBS
var appName string

func init() {

	cfg = commonConfig.Get()
	appName = cfg.GetApplication_Name()
	err := jobs.Initialise(cfg)
	if err != nil {
		logHandler.ServiceLogger.Fatalf("[%v] Error: %v", domain.String(), err)
		panic(err)
	}
}
