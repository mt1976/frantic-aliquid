package authorityStore

import (
	"github.com/mt1976/frantic-aliquid/app/business/domains"
	"github.com/mt1976/frantic-aliquid/app/dao/behaviourStore"
	"github.com/mt1976/frantic-core/dao/actions"
)

var LIST behaviourStore.Behaviour_Store
var VIEW behaviourStore.Behaviour_Store
var EDIT behaviourStore.Behaviour_Store
var UPDATE behaviourStore.Behaviour_Store
var GRANT behaviourStore.Behaviour_Store
var REVOKE behaviourStore.Behaviour_Store
var behaviourDomain = domains.AUTHORITY

func DeclareBehaviours() {
	LIST, _ = behaviourStore.Declare(actions.LIST, behaviourDomain, "Authorities", "")
	LIST.DefaultRights()
	VIEW, _ = behaviourStore.Declare(actions.VIEW, behaviourDomain, "View Authority", "")
	VIEW.DefaultRights()
	EDIT, _ = behaviourStore.Declare(actions.EDIT, behaviourDomain, "Edit Authority", "")
	EDIT.DefaultRights()
	UPDATE, _ = behaviourStore.Declare(actions.UPDATE, behaviourDomain, "Update Authority", "")
	UPDATE.DefaultRights()
}
