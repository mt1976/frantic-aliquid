package jobs

import (
	"reflect"
	"runtime"

	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/stringHelpers"
	"github.com/mt1976/frantic-core/timing"
)

// HostJob represents a job that can be scheduled and run periodically.
//
// Fields:
//
//	databaseAccessors []func() ([]*database.DB, error) (optional): A slice of functions that return database instances.
//	  Uncomment this field for multi-database jobs.
//
// Methods:
//
//	Name() string:
//	  Returns the name of the job.
//
//	Schedule() string:
//	  Returns the schedule for the job in cron format.
//
//	Description() string:
//	  Returns a description of the job.
//
//	Run() error:
//	  Executes the job. Starts a timing clock, runs pre-processing, processes the job, runs post-processing, and stops the clock.
//	  Returns an error if any step fails.
//
//	Service() func():
//	  Returns a function that runs the job and logs any errors.
//
//	AddDatabaseAccessFunctions(fn func() ([]*database.DB, error)):
//	  Adds a function to the databaseAccessors slice. This method is currently not implemented and will panic if called.
//
// Example usage:
//
//	job := &HostJob{}
//	job.Service()()
type HostJob struct {
	// Uncomment the following line for multi-database jobs
	databaseAccessors []func() ([]*database.DB, error)
}

// Name returns the name of the job.
//
// Returns:
//
//	string: The name of the job.
func (j *HostJob) Name() string {
	return "Host Status Monitoring"
}

// Schedule returns the schedule for the job in cron format.
//
// Returns:
//
//	string: The schedule for the job in quartz cron format.
func (j *HostJob) Schedule() string {
	// Every 2 minutes
	return "0 0/1 * * * *"
}

// Description returns a description of the job.
//
// Returns:
//
//	string: A description of the job.
func (j *HostJob) Description() string {
	return "This job monitors the status of all hosts in the system, pinging each host to check if it is active."
}

// Run executes the job. Starts a timing clock, runs pre-processing, processes the job, runs post-processing, and stops the clock.
//
// Returns:
//
//	error: An error if any step fails, otherwise nil.
func (j *HostJob) Run() error {
	clock := timing.Start(jobs.CodedName(j), actions.PROCESS.GetCode(), j.Description())
	jobs.PreRun(j)

	if len(j.databaseAccessors) > 0 && j.databaseAccessors != nil {
		logHandler.ServiceLogger.Printf("[%v] Running '%v' job across %v locations", jobs.CodedName(j), j.Name(), len(j.databaseAccessors))
		for _, f := range j.databaseAccessors {

			logHandler.ServiceLogger.Printf("[%v] Getting list of databases from '%v'", jobs.CodedName(j), runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
			dbList, err := f()
			if err != nil {
				logHandler.ErrorLogger.Printf("[%v] Error=[%v]", j.Name(), err.Error())
				continue
			}
			logHandler.ServiceLogger.Printf("[%v] Running '%v' job across %v database(s)", jobs.CodedName(j), j.Name(), len(dbList))
			for _, db := range dbList {
				jobProcessor(j, db)
			}
		}

	} else {
		jobProcessor(j, nil)
	}
	jobs.PostRun(j)
	clock.Stop(1)
	return nil
}

// Service returns a function that runs the job and logs any errors.
//
// Returns:
//
//	func(): A function that runs the job and logs any errors.
func (j *HostJob) Service() func() {
	return func() {
		err := j.Run()
		if err != nil {
			logHandler.ServiceLogger.Panicf("[%v] %v Error=[%v]", jobs.CodedName(j), j.Name(), err.Error())
			panic(err)
		}
	}
}

// AddDatabaseAccessFunctions adds a function to the databaseAccessors slice.
//
// This method is currently not implemented and will panic if called.
//
// Parameters:
//
//	fn func() ([]*database.DB, error): A function that returns a slice of pointers to `database.DB` and an error.
func (job *HostJob) AddDatabaseAccessFunctions(fn func() ([]*database.DB, error)) {
	job.databaseAccessors = append(job.databaseAccessors, fn)
}

// jobProcessor is the main function that processes the job.
//
// This function is called by the Run method of the HostJob struct to perform the main processing logic of the job.
//
// Parameters:
//
//	j *HostJob: A pointer to the HostJob instance that is being processed.
func jobProcessor(j *HostJob, db *database.DB) {
	// This is the main function
	jobName := stringHelpers.SQuote(j.Name())
	appName := cfg.GetApplication_Name()
	// Ensure the job has the correct database connection
	if db == nil {
		logHandler.EventLogger.Printf("[%v] Running %v tasks with default database for %v", jobs.CodedName(j), jobName, appName)
	} else {
		logHandler.EventLogger.Printf("[%v] Running %v tasks with database=[%v-%v.db]", jobs.CodedName(j), jobName, appName, db.Name)
	}

	//authorityStore.Worker(j, db)

	// Report the completion of the job
	if db == nil {
		logHandler.EventLogger.Printf("[%v] Completed %v tasks with default database for %v", jobs.CodedName(j), jobName, appName)
	} else {
		logHandler.EventLogger.Printf("[%v] Completed %v tasks with database=[%v-%v.db]", jobs.CodedName(j), jobName, appName, db.Name)
	}
}
