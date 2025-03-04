package jobs

import (
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
)

//var SessionExpiryWorker jobs.Job = &jobs.Job{}

func Start() {
	// Start the job
	logHandler.ServiceLogger.Printf("[%v] Queue - Starting", domain.String())
	// Database Backup
	// DatabaseBackup.AddDatabaseAccessFunctions(settingStore.FetchDatabaseInstances())
	// DatabaseBackup.AddDatabaseAccessFunctions(sessionStore.GetDatabaseConnections())
	// DatabaseBackup.AddDatabaseAccessFunctions(userStore.FetchDatabaseInstances())
	// jobs.AddJobToScheduler(DatabaseBackup)
	// // Prune the archive of backups
	// jobs.AddJobToScheduler(DatabasePrune)
	// // Check the status of the hosts
	// jobs.AddJobToScheduler(HostWorker)
	// // Check the status of the sessions
	// SessionExpiryWorker.AddDatabaseAccessFunctions(sessionStore.GetDatabaseConnections())
	// jobs.AddJobToScheduler(SessionExpiryWorker)
	// // Start all the background jobs

	jobs.StartScheduler()

	//logger.ServiceLogger.Printf("[%+v]", scheduledTasks.Entries())
	logHandler.ServiceLogger.Printf("[%v] Queue - Started", domain.String())
}
