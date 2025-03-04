package authorityStore

// Data Access Object Authority
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"context"

	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

var sessionExpiry = 20 // default to 20 mins
var activeDB *database.DB
var initialised bool = false // default to false
var cfg *commonConfig.Settings
var SEP = "∘" // separator for key generation"

func Initialise(ctx context.Context) {
	//logHandler.EventLogger.Printf("Initialising %v", domain)
	timing := timing.Start(domain, actions.INITIALISE.GetCode(), "Initialise")
	cfg = commonConfig.Get()

	// For a specific database connection, use NamedConnect, otherwise use Connect
	//activeDB = database.ConnectToNamedDB("Authority")
	activeDB = database.Connect()
	initialised = true
	timing.Stop(1)
	logHandler.EventLogger.Printf("Initialised %v", domain)
}
