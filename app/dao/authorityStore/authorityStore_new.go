package authorityStore

import (
	"context"
	"fmt"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao"
	"github.com/mt1976/frantic-core/dao/actions"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func (u *Authority_Store) Create(ctx context.Context, act audit.Action, note string) (Authority_Store, error) {
	dao.CheckDAOReadyState(domain, audit.CREATE, initialised) // Check the DAO has been initialised, Mandatory.

	//logHandler.InfoLogger.Printf("New %v (%v=%v)", domain, FIELD_ID, field1)
	clock := timing.Start(domain, actions.CREATE.GetCode(), fmt.Sprintf("%v", u.Key))

	// Record the create action in the audit data
	auditErr := u.Audit.Action(context.TODO(), act.WithMessage(fmt.Sprintf("New %v created %v", domain, u.Key)))
	if auditErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panic(commonErrors.WrapDAOUpdateAuditError(domain, u.Key, auditErr))
	}

	// Save the status instance to the database
	writeErr := activeDB.Create(u)
	if writeErr != nil {
		// Log and panic if there is an error creating the status instance
		logHandler.ErrorLogger.Panic(commonErrors.WrapDAOCreateError(domain, u.Key, writeErr))
		//	panic(writeErr)
	}

	//logHandler.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.CREATE, strings.ToUpper(domain), record.ID, fmt.Sprintf("New %v: %v", domain, field1))
	clock.Stop(1)
	return *u, nil
}
