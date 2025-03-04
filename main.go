package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mt1976/frantic-aliquid/app/business/auth"
	"github.com/mt1976/frantic-aliquid/app/business/translation"
	"github.com/mt1976/frantic-aliquid/app/dao"
	"github.com/mt1976/frantic-aliquid/app/dao/behaviourStore"
	"github.com/mt1976/frantic-aliquid/app/jobs"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/dockerHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
	"github.com/mt1976/frantic-core/timing"
)

var cfg *commonConfig.Settings
var appName string

// var systemUserUserUID string
// var systemUserUserName string
// var systemUser string

func init() {
	cfg = commonConfig.Get()
	appName = "frantic-aliquid"

	// systemUserUserUID = cfg.GetServiceUser_UID()
	// systemUserUserName = cfg.GetServiceUser_Name()
	// systemUser = cfg.GetServiceUser_UserCode()
}

func main() {

	logHandler.EventLogger.Printf("[%v] Starting...", appName)
	logHandler.EventLogger.Printf("[%v] Version: %v", appName, cfg.GetApplication_Version())
	logHandler.EventLogger.Printf("[%v] Author: %v", appName, cfg.GetApplication_Author())
	logHandler.EventLogger.Printf("[%v] Build Date: %v", appName, cfg.GetApplication_Copyright())
	logHandler.EventLogger.Printf("[%v] Release Date: %v", appName, cfg.GetApplication_ReleaseDate())
	logHandler.EventLogger.Printf("[%v] Path: %v", appName, cfg.GetApplication_HomePath())
	logHandler.EventLogger.Printf("[%v] Locale: %v", appName, cfg.GetApplication_Locale())
	logHandler.EventLogger.Printf("[%v] Environment: %v", appName, cfg.GetApplication_Environment())
	logHandler.EventLogger.Printf("")

	startupSequence := timing.Start(appName, "Startup", "Sequence")

	logHandler.EventLogger.Printf("[%v] Loading Docker Default Payloads...", appName)
	err := dockerHelpers.DeployDefaultsPayload()
	if err != nil {
		logHandler.ErrorLogger.Fatal(err.Error())
	}
	logHandler.EventLogger.Printf("[%v] Docker Default Payloads Loaded", appName)

	err = error(nil)
	logHandler.EventLogger.Printf("[%v] Database Connection - Starting...", appName)

	database.Connect()

	logHandler.EventLogger.Printf("[%v] Database Connection - Connecting...", appName)
	err = dao.Initialise(cfg)
	if err != nil {
		logHandler.ErrorLogger.Fatal(err.Error())
	}
	logHandler.EventLogger.Printf("[%v] Database Connection - Connected", appName)

	logHandler.EventLogger.Printf("[%v] Backup - Starting...", appName)

	// Belt and braces, reinitialise the database
	err = dao.Initialise(cfg)
	if err != nil {
		logHandler.ErrorLogger.Fatal(err.Error())
	}

	logHandler.EventLogger.Printf("[%v] Backup - Done", appName)

	logHandler.EventLogger.Printf("[%v] Starting %v...", appName, appName)

	importDefaults()

	logHandler.EventLogger.Printf("[%v] Initialise - Routes - Started", appName)

	// Grant the system user all authorities
	// Get all authorities
	logHandler.EventLogger.Printf("[%v] Initialise - System Users - Started", appName)

	//au := assignSystemUserPermissions(su)
	logHandler.EventLogger.Printf("[%v] Initialise - System Users - Completed", appName)

	//logHandler.EventLogger.Printf("Admin User [%v] [%v] Created", au.UserName, au.UserCode)

	router := httprouter.New()

	// ANNOUNCE ROUTES ABOVE
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, translation.Get("404 page not found"), http.StatusNotFound)
	})

	logHandler.EventLogger.Printf("[%v] Initialise - Routes - Completed", appName)

	// Start the job processor
	logHandler.EventLogger.Printf("[%v] Job Scheduler Starting...", appName)
	jobs.Start()
	logHandler.EventLogger.Printf("[%v] Job Scheduler Started", appName)
	startupSequence.Stop(1)
	// Get a list of all users

	logHandler.EventLogger.Printf("[%v] Startup - Granting User Access Rights to System User", appName)
	// Grant all users access rights

	UserKey := cfg.GetServiceUser_UID()
	UserCode := cfg.GetServiceUser_UserCode()

	logHandler.EventLogger.Printf("[%v] Grant User Access Rights [%v] [%v]", appName, UserCode, UserKey)
	su := messageHelpers.UserMessage{
		Code: UserCode,
		Key:  UserKey,
	}
	err = auth.GrantAllUserAccessRights(su)
	if err != nil {
		logHandler.ErrorLogger.Println(err.Error())
	}

	logHandler.EventLogger.Printf("[%v] Startup Complete", appName)
}

func impexTest(exporter func() error, importer func() error, counter func() (int, error)) error {
	var err error
	err = exporter()
	if err != nil {
		logHandler.ErrorLogger.Fatal(err.Error())
		return err
	}

	err = importer()
	if err != nil {
		logHandler.ErrorLogger.Fatal(err.Error())
		return err
	}

	hsc, err := counter()
	if err != nil {
		logHandler.ErrorLogger.Fatal(err.Error())
		return err
	}
	logHandler.EventLogger.Printf("[%v] Impex - Count [%v]", appName, hsc)
	return nil
}

func importDefaults() {
	logHandler.EventLogger.Printf("[%v] Importing Defaults - %v - Started", appName, "Behaviours")

	err := behaviourStore.ImportCSV()
	if err != nil {
		logHandler.ErrorLogger.Fatal(err.Error())
	}

	logHandler.EventLogger.Printf("[%v] Importing Defaults - %v - Completed", appName, "Behaviours")
}
