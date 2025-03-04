package authorityStore

// Data Access Object Authority
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"github.com/mt1976/frantic-aliquid/app/dao/behaviourStore"
	audit "github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/messageHelpers"
)

// Authority_Store represents a Authority_Store entity.
type Authority_Store struct {
	ID           int                            `storm:"id,increment=100000"` // primary key with auto increment
	Key          string                         `storm:"unique,index"`        // key
	Raw          string                         `storm:"unique,index"`        // raw ID before encoding
	UserKey      string                         `storm:"index"`               // user key
	UserCode     string                         `storm:"index"`               // user code
	User         messageHelpers.UserMessage     `csv:"-"`                     // user details
	BehaviourKey string                         `storm:"index"`               // behaviour ID
	Behaviour    behaviourStore.Behaviour_Store `csv:"-"`                     // behaviour details
	Info         string                         `csv:"-"`                     // info
	Source       string                         `storm:"index"`               // source
	Audit        audit.Audit                    `csv:"-"`                     // audit data
}

// Define the field set as names
var (
	FIELD_ID  = "ID"
	FIELD_Key = "Key"
	FIELD_Raw = "Raw"
	// Add fields below here
	FIELD_UserKey     = "UserKey"
	FIELD_User        = "User"
	FIELD_BehaviourID = "BehaviourID"
	FIELD_Behaviour   = "Behaviour"
	FIELD_Source      = "Source"
	// Add fields above here
	FIELD_Audit = "Audit"
)

var domain = "Authority"
