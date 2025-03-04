package behaviourStore

// Data Access Object Behaviour - Behaviour_Store
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"context"
	"fmt"
	"strings"

	"github.com/mt1976/frantic-aliquid/app/business/domains"
	"github.com/mt1976/frantic-aliquid/app/business/translation"
	"github.com/mt1976/frantic-aliquid/app/web/permissions"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/dao/lookup"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/paths"
	"github.com/mt1976/frantic-core/stringHelpers"
	"github.com/mt1976/frantic-core/timing"
)

func (record *Behaviour_Store) prepare() (Behaviour_Store, error) {
	//os.Exit(0)
	//logger.ErrorLogger.Printf("ACT: VAL Validate")
	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	return *record, nil
}

func (record *Behaviour_Store) calculate() error {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	// Calculate the duration in days between the start and end dates
	return nil
}

func (record *Behaviour_Store) isDuplicateOf(inRawKey string) (Behaviour_Store, error) {

	dao.CheckDAOReadyState(domain, audit.PROCESS, initialised) // Check the DAO has been initialised, Mandatory.

	valRec, err := GetBy(FIELD_Raw, inRawKey)
	if err == nil {
		return valRec, commonErrors.ErrorDuplicate
	}
	return Behaviour_Store{}, nil
	// // Get all status
	// recordList, err := GetAll()
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Getting all %v failed %v", domain, err.Error())
	// 	return Behaviour_Store{}, err
	// }

	// // range through status list, if status code is found and deletedby is empty then return error
	// for _, checkRecord := range recordList {
	// 	//s.Dump(!,strings.ToUpper(code) + "-uchk-" + s.Code)
	// 	testValue := strings.ToUpper(inRawKey)
	// 	checkValue := strings.ToUpper(checkRecord.Raw)
	// 	logHandler.InfoLogger.Printf("CHK: TestValue:[%v] CheckValue:[%v]", testValue, checkValue)
	// 	//logger.InfoLogger.Printf("CHK: Code:[%v] s.Code:[%v] s.Audit.DeletedBy:[%v]", testCode, s.Code, s.Audit.DeletedBy)
	// 	if checkValue == testValue && checkRecord.Audit.DeletedBy == "" {
	// 		logHandler.WarningLogger.Printf("Duplicate %v, %v already in use", domain, record.Raw)
	// 		return checkRecord, commonErrors.ErrorDuplicate
	// 	}
	// }

	// return Behaviour_Store{}, nil
}

func ClearDown(ctx context.Context) error {
	logHandler.EventLogger.Printf("Clearing %v", domain)

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

func Declare(inAction actions.Action, inDomain domains.Domain, note, source string) (Behaviour_Store, error) {
	logHandler.EventLogger.Printf("Declaring (%v) (%v) (%v) (%v)", inDomain, inAction.GetName(), source, note)
	b := Behaviour_Store{}
	clock := timing.Start(inDomain.String(), actions.CREATE.GetCode(), fmt.Sprintf("%v", inAction.GetName()))
	// Default all Permissions to True
	b.DefaultRights()
	b.Action = inAction
	b.Domain = strings.ToLower(inDomain.String())
	if source == "" {
		source = cfg.GetApplication_Name()
		logHandler.WarningLogger.Printf("No source provided for %v %v, using %v", inDomain, inAction.GetName(), source)
	}
	b.Source = strings.ToLower(source)
	b.HTMLPageID = strings.ToLower(inDomain.String() + "_" + inAction.GetName())
	b.Raw = strings.ToLower(source + SEP + inAction.GetName() + SEP + inDomain.String())
	b.Key = idHelpers.Encode(b.Raw)
	if note == "" {
		note = inAction.GetName() + " " + inDomain.String()
	}
	b.Note = note
	// Check for duplicates
	xStatus, err := b.isDuplicateOf(b.Raw)
	if err == commonErrors.ErrorDuplicate {
		// This is OK, do nothing as this is a duplicate record
		// we ignore duplicates.
		logHandler.WarningLogger.Printf(translation.Get("Duplicate %v %v already in use, skipping"), domain, stringHelpers.DQuote(b.Raw))
		clock.Stop(1)
		return xStatus, nil
	}
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] New %s %v %v", domain, err.Error(), b.Key, b.Raw)
		clock.Stop(0)
		return Behaviour_Store{}, err
	}

	// Record the create action in the audit data
	_ = b.Audit.Action(context.TODO(), audit.CREATE.WithMessage(fmt.Sprintf("New %v %v", b.ID, b.Raw)))

	// Save the dest instance to the database
	err = activeDB.Create(&b)
	if err != nil {
		// Log and panic if there is an error creating the dest instance
		logHandler.ErrorLogger.Printf("[%v] Create %s %v %v", inDomain, err.Error(), b.ID, b.Raw)
		panic(err)
	}

	logHandler.SecurityLogger.Printf("[%v] [%v] ID=[%v] Raw=[%v]", inDomain, strings.ToUpper(inDomain.String()), b.ID, b.Raw)

	return b, nil
}

func (b *Behaviour_Store) Title() string {
	return b.Note
}
func (b *Behaviour_Store) ActionName() string {
	return strings.ToTitle(b.Action.GetName())
}

func (b *Behaviour_Store) DefineRights(list, new, view, edit, update, delete, activate, deactivate, action bool) {

	b.Permissions = permissions.Rights{
		List:       list,
		New:        new,
		View:       view,
		Edit:       edit,
		Update:     update,
		Delete:     delete,
		Activate:   activate,
		Deactivate: deactivate,
		Action:     action,
	}
}

func (b *Behaviour_Store) CloneRightsOf(source Behaviour_Store) {
	b.Permissions = source.Permissions
}
func (b *Behaviour_Store) DefaultRights() {
	b.Permissions.Defaults()
}

func (b *Behaviour_Store) Rights() permissions.Rights {
	return b.Permissions
}

func (b *Behaviour_Store) Name() string {
	return b.Note
}

func (b *Behaviour_Store) Page() string {
	return b.HTMLPageID
}

func (b *Behaviour_Store) ActionType() actions.Action {
	return b.Action
}

func (b *Behaviour_Store) Template() string {
	html := paths.HTML().String() + b.HTMLPageID + ".html"
	logHandler.TraceLogger.Printf("[BEHAVIOUR] Template=[%v]", html)
	return html
}

func (b *Behaviour_Store) Is(c Behaviour_Store) bool {
	if b.Note == c.Note && b.HTMLPageID == c.HTMLPageID {
		return true
	}
	return false
}

func BuildLookup() (lookup.Lookup, error) {

	//logger.InfoLogger.Printf("BuildLookup")

	// Get all status
	Activities, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("ERROR Getting all status: %v", err)
		return lookup.Lookup{}, err
	}

	// Create a new Lookup
	var rtnList lookup.Lookup
	rtnList.Data = make([]lookup.LookupData, 0)

	// range through Behaviour list, if status code is found and deletedby is empty then return error
	for _, a := range Activities {
		rtnList.Data = append(rtnList.Data, lookup.LookupData{Key: a.Key, Value: a.Raw})
	}

	return rtnList, nil
}
