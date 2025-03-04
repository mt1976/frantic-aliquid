package behaviourStore

import (
	"context"

	"github.com/mt1976/frantic-aliquid/app/business/domains"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/importExportHelper"
	"github.com/mt1976/frantic-core/logHandler"
)

func ExportCSV() error {

	Initialise(context.TODO())

	exportsList, err := GetAll()
	if err != nil {
		logHandler.ExportLogger.Panicf("Error Getting all %v: %v", domain, err.Error())
	}

	return importExportHelper.ExportCSV(domain, exportsList)
}

func ImportCSV() error {
	logHandler.ImportLogger.Printf("Importing %v", domain)

	Initialise(context.TODO())
	err := importExportHelper.ImportCSV(domain, &Behaviour_Store{}, load)

	return err
}

// load is a helper function to create a new entry instance and save it to the database
// It should be customised to suit the specific requirements of the entryination table/DAO.
func load(inOriginal **Behaviour_Store) (string, error) {

	original := **inOriginal
	//	logHandler.ImportLogger.Printf("Declaring %v [%v] [%v]", domain, original.Raw, original.Text)

	//logger.InfoLogger.Printf("ACT: NEW New %v %v %v", tableName, name, entryination)
	// u := Behaviour_Store{}
	// u.Key = idHelpers.Encode(strings.ToUpper(original.Raw))
	// u.Raw = original.Raw
	// u.Text = original.Text
	// // u.Action = original.Action
	// u.Domain = original.Domain

	importAction := actions.New(original.Action.Name)
	_, err := Declare(importAction, domains.Special(original.Domain), original.Note, original.Source)
	if err != nil {
		logHandler.ImportLogger.Panicf("Error importing behaviour: %v", err.Error())
	}

	//logHandler.ImportLogger.Printf("Declared %v [%v] [%v]", domain, bh.Raw, bh.Text)
	// Return the created entry and nil error
	return original.Raw, nil
}
