package behaviourStore

import (
	"github.com/mt1976/frantic-aliquid/app/business/domains"
	"github.com/mt1976/frantic-core/dao/actions"
)

var LIST Behaviour_Store
var VIEW Behaviour_Store
var EDIT Behaviour_Store
var UPDATE Behaviour_Store
var ENABLE Behaviour_Store
var DISABLE Behaviour_Store

func DeclareBehaviours() {
	LIST, _ = Declare(actions.LIST, domains.BEHAVIOUR, "List Behaviours", "")
	LIST.DefaultRights()
	VIEW, _ = Declare(actions.VIEW, domains.BEHAVIOUR, "View Behaviour", "")
	VIEW.DefaultRights()
	EDIT, _ = Declare(actions.EDIT, domains.BEHAVIOUR, "Edit Behaviour", "")
	EDIT.DefaultRights()
	UPDATE, _ = Declare(actions.UPDATE, domains.BEHAVIOUR, "Update Behaviour", "")
	UPDATE.DefaultRights()

}
