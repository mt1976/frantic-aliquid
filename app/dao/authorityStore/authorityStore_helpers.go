package authorityStore

// Data Access Object Authority - Authoritystore
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"context"
	"fmt"
	"strings"

	"github.com/mt1976/frantic-aliquid/app/dao/behaviourStore"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/dao/lookup"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
	"github.com/mt1976/frantic-core/timing"
)

func (record *Authority_Store) prepare() (Authority_Store, error) {
	//os.Exit(0)
	//logHandler.ErrorlogHandler.Printf("ACT: VAL Validate")
	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	return *record, nil
}

func (record *Authority_Store) calculate() error {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	// Calculate the duration in days between the start and end dates
	return nil
}

func (record *Authority_Store) IsDuplicateOf(id string) (Authority_Store, error) {
	return record.isDuplicateOf(id)
}
func (record *Authority_Store) isDuplicateOf(id string) (Authority_Store, error) {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	//TODO: Could be replaced with a simple read...

	// Get all status
	recordList, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("Getting all %v failed %v", domain, err.Error())
		return Authority_Store{}, err
	}

	// range through status list, if status code is found and deletedby is empty then return error
	for _, checkRecord := range recordList {
		//s.Dump(!,strings.ToUpper(code) + "-uchk-" + s.Code)
		testValue := strings.ToUpper(id)
		checkValue := strings.ToUpper(checkRecord.Key)
		//logHandler.InfologHandler.Printf("CHK: TestValue:[%v] CheckValue:[%v]", testValue, checkValue)
		//logHandler.InfologHandler.Printf("CHK: Code:[%v] s.Code:[%v] s.Audit.DeletedBy:[%v]", testCode, s.Code, s.Audit.DeletedBy)
		if checkValue == testValue && checkRecord.Audit.DeletedBy == "" {
			logHandler.WarningLogger.Printf("Duplicate %v, %v already in use", strings.ToUpper(domain), record.ID)
			return checkRecord, commonErrors.ErrorDuplicate
		}
	}

	return Authority_Store{}, nil
}

func ClearDown(ctx context.Context) error {
	logHandler.InfoLogger.Printf("Clearing %v", domain)

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	clock := timing.Start(domain, actions.CLEAR.GetCode(), "INITIALISE")

	// Delete all active session recordList
	recordList, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Print(commonErrors.WrapDAOInitialisationError(domain, err).Error())
		clock.Stop(0)
		return commonErrors.WrapDAOInitialisationError(domain, err)
	}

	noRecords := len(recordList)
	count := 0

	for thisRecord, record := range recordList {
		logHandler.InfoLogger.Printf("Deleting %v (%v/%v) %v", domain, thisRecord, noRecords, record.Key)
		delErr := Delete(ctx, record.ID, fmt.Sprintf("Clearing %v %v @ initialisation ", domain, record.ID))
		if delErr != nil {
			logHandler.ErrorLogger.Print(commonErrors.WrapDAOInitialisationError(domain, delErr).Error())
			continue
		}
		count++
	}

	clock.Stop(count)

	return nil
}

func BuildLookup() (lookup.Lookup, error) {
	// Added to support older code.
	return GetLookup(FIELD_Key, FIELD_Raw)
}

func GetByUserKey(id string) []Authority_Store {
	var rtnList []Authority_Store
	// Get all status
	activityList, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		return rtnList
	}

	// range through status list, if status code is found and deletedby is empty then return error
	for _, a := range activityList {
		if a.UserKey == id {
			rtnList = append(rtnList, a)
		}
	}

	return rtnList
}

// func GetByUserID(id int) []Authority_Store {
// 	var rtnList []Authority_Store
// 	// Get all status
// 	activityList, err := GetAll()
// 	if err != nil {
// 		logHandler.ErrorLogger.Printf("ERROR Getting all status: %v", err)
// 		return rtnList
// 	}

// 	// range through status list, if status code is found and deletedby is empty then return error
// 	for _, a := range activityList {
// 		if a.UserID == id {
// 			rtnList = append(rtnList, a)
// 		}
// 	}

// 	return rtnList
// }

func GetByBehaviourKey(id string) []Authority_Store {
	var rtnList []Authority_Store
	// Get all status
	activityList, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		return rtnList
	}

	// range through status list, if status code is found and deletedby is empty then return error
	for _, a := range activityList {
		if a.Key == id {
			rtnList = append(rtnList, a)
		}
	}

	return rtnList
}

func BuildBaseRecord(usr messageHelpers.UserMessage, ben behaviourStore.Behaviour_Store) (Authority_Store, error) {
	g := Authority_Store{}
	g.Raw = usr.Code + SEP + ben.Raw
	g.Key = idHelpers.Encode(g.Raw)
	g.UserKey = usr.Key
	g.User = usr
	g.UserCode = usr.Code
	g.Behaviour = ben
	g.BehaviourKey = ben.Key
	g.Source = usr.Source
	if !strings.EqualFold(ben.Source, usr.Source) {
		logHandler.ErrorLogger.Printf("Source mismatch: User: '%v', Behaviour: '%v'", usr.Code, ben.Raw)
		return Authority_Store{}, fmt.Errorf("source mismatch: User: '%v', Behaviour: '%v'", usr.Code, ben.Raw)
	}
	g.Info = fmt.Sprintf("User: '%v', Behaviour: '%v', Key='%v%v%v'", usr.Code, ben.Raw, usr.Code, SEP, ben.Raw)

	return g, nil
}

func (a Authority_Store) BuildMessage() (messageHelpers.AuthorityMessage, error) {
	am := messageHelpers.AuthorityMessage{}
	am.Key = a.Key
	am.User = messageHelpers.UserMessage{Key: a.User.Code, Code: a.User.Code, Source: a.User.Source}
	am.Behaviour = messageHelpers.BehaviourMessage{Key: a.Behaviour.Key}
	return am, nil
}
