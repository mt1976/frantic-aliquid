package behaviourStore

// Data Access Object Behaviour
// Version: 0.2.0
// Updated on: 2021-09-10

import (
	"github.com/mt1976/frantic-aliquid/app/web/permissions"
	"github.com/mt1976/frantic-core/dao/actions"
	audit "github.com/mt1976/frantic-core/dao/audit"
)

// Behaviour_Store represents a Behaviour_Store entity.
type Behaviour_Store struct {
	ID          int                `csv:"-" storm:"id,increment=100000"` // primary key with auto increment
	Key         string             `csv:"-" storm:"unique"`              // key
	Raw         string             `csv:"raw" storm:"unique"`            // raw ID before encoding
	HTMLPageID  string             `csv:"-"`
	Permissions permissions.Rights `csv:"-"`
	Note        string             `csv:"text"`
	Action      actions.Action     `csv:"action" storm:"index"`
	Domain      string             `csv:"domain" storm:"index"`
	Source      string             `csv:"source" storm:"index"`
	Audit       audit.Audit        `csv:"-"`
}

// Define the field set as names
var (
	FIELD_ID          = "ID"
	FIELD_Key         = "Key"
	FIELD_Raw         = "Raw"
	FIELD_HTMLPageID  = "HTMLPageID"
	FIELD_Permissions = "Permissions"
	FIELD_Text        = "Text"
	FIELD_Action      = "Action"
	FIELD_Domain      = "Domain"
	FIELD_Source      = "Source"
	FIELD_Audit       = "Audit"
)

var domain = "Behaviour"
