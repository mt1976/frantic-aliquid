package dao

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/mt1976/frantic-aliquid/app/business/domains"
	"github.com/mt1976/frantic-aliquid/app/dao/authorityStore"
	"github.com/mt1976/frantic-aliquid/app/dao/behaviourStore"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

var name = "DAO"
var Version = 1
var DB *storm.DB
var tableName = "authority"

func Initialise(cfg *commonConfig.Settings) error {
	clock := timing.Start(name, actions.INITIALISE.GetCode(), "")
	logHandler.EventLogger.Printf("Initialising %v - Started", name)

	authorityStore.Initialise(context.TODO())
	behaviourStore.Initialise(context.TODO())

	authorityStore.DeclareBehaviours()
	behaviourStore.DeclareBehaviours()

	_, _ = behaviourStore.Declare(actions.ROUTE, domains.Special("home"), "Home", "")
	_, _ = behaviourStore.Declare(actions.MESSAGE, domains.Special("oops"), "Action Failed", "")
	_, _ = behaviourStore.Declare(actions.MESSAGE, domains.Special("fail"), "Action Failed", "")
	_, _ = behaviourStore.Declare(actions.MESSAGE, domains.Special("success"), "Action Succeeded", "")

	logHandler.EventLogger.Printf("Initialising %v - Complete", name)
	clock.Stop(1)
	return nil
}

func ExportAllToCSV() {
	authorityStore.ExportCSV()
	behaviourStore.ExportCSV()
}
