package behaviourStore

// func New(field1 int, field2 string) (Behaviour_Store, error) {

// 	dao.CheckDAOReadyState(domain, audit.CREATE, initialised) // Check the DAO has been initialised, Mandatory.

// 	//logHandler.InfoLogger.Printf("New %v (%v=%v)", domain, FIELD_ID, field1)
// 	clock := timing.Start(domain, actions.CREATE.GetCode(), fmt.Sprintf("%v", field1))

// 	sessionID := idHelpers.GetUUID()

// 	// Create a new struct
// 	record := Behaviour_Store{Field1: field1}
// 	record.Key = idHelpers.Encode(sessionID)
// 	record.Raw = sessionID
// 	record.Field1 = field1
// 	record.Field2 = field2
// 	record.Field3 = time.Now().Add(time.Minute * time.Duration(sessionExpiry))

// 	// Record the create action in the audit data
// 	auditErr := record.Audit.Action(context.TODO(), audit.CREATE.WithMessage(fmt.Sprintf("New %v created %v", domain, field1)))
// 	if auditErr != nil {
// 		// Log and panic if there is an error creating the status instance
// 		logHandler.ErrorLogger.Panic(commonErrors.WrapDAOUpdateAuditError(domain, record.ID, auditErr))
// 	}

// 	// Save the status instance to the database
// 	writeErr := activeDB.Create(&record)
// 	if writeErr != nil {
// 		// Log and panic if there is an error creating the status instance
// 		logHandler.ErrorLogger.Panic(commonErrors.WrapDAOCreateError(domain, record.ID, writeErr))
// 		//	panic(writeErr)
// 	}

// 	//logHandler.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.CREATE, strings.ToUpper(domain), record.ID, fmt.Sprintf("New %v: %v", domain, field1))
// 	clock.Stop(1)
// 	return record, nil
// }
