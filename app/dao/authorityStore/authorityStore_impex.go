package authorityStore

import (
	"context"

	"github.com/mt1976/frantic-aliquid/app/dao/behaviourStore"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
)

// importProcessor is a helper function to create a new entry instance and save it to the database
// It should be customised to suit the specific requirements of the entryination table/DAO.
func importProcessor(inOriginal **Authority_Store) (string, error) {

	importedData := **inOriginal

	userKey := importedData.UserKey
	userCOde := importedData.UserCode
	behaviourKey := importedData.BehaviourKey

	importRawKey := importedData.Raw

	// usr, err := userStore.GetBy(userStore.FIELD_Key, userKey)
	// if err != nil {
	// 	logHandler.ImportLogger.Panicf("Error importing %v: %v", domain, err.Error())
	// 	return importRawKey, err
	// }

	ben, err := behaviourStore.GetBy(behaviourStore.FIELD_Key, behaviourKey)
	if err != nil {
		logHandler.ImportLogger.Panicf("Error importing %v: %v", domain, err.Error())
		return importRawKey, err
	}

	usr := messageHelpers.UserMessage{
		Code: userCOde,
		Key:  userKey,
	}

	newA, err := BuildBaseRecord(usr, ben)
	if err != nil {
		logHandler.ImportLogger.Panicf("Error importing %v: %v", domain, err.Error())
		return importRawKey, err
	}

	// Check if the record exists
	_, err = GetBy(FIELD_Raw, newA.Raw)
	if err == nil {
		// Record exists, do nothing
		logHandler.WarningLogger.Printf("%v Record exists, skipping: %v", domain, newA.Raw)
		return importRawKey, nil

	} else {
		// Record does not exist, create

		_, err = newA.Create(context.TODO(), audit.IMPORT, "Imported")
		if err != nil {
			logHandler.ImportLogger.Panicf("Error importing %v: %v", domain, err.Error())
			return importRawKey, err
		}
	}

	return importRawKey, nil
}
