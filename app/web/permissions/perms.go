package permissions

import (
	"fmt"

	logger "github.com/mt1976/frantic-core/logHandler"
)

type Rights struct {
	List       bool
	New        bool
	View       bool
	Edit       bool
	Update     bool
	Delete     bool
	Activate   bool
	Deactivate bool
	Action     bool
}

func (r *Rights) Spew() {
	msgTemplate := "List=[%t],New=[%t],View=[%t],Edit=[%t],Update=[%t],Delete=[%t],Activate=[%t],Deactivate=[%t],Action=[%t]"

	msg := fmt.Sprintf(msgTemplate, r.List, r.New, r.View, r.Edit, r.Update, r.Delete, r.Activate, r.Deactivate, r.Action)

	prefix := "[PERMISSIONS] %v"
	logger.InfoLogger.Printf(prefix, msg)
}

func (r *Rights) Defaults() {
	r.List = true
	r.New = true
	r.View = true
	r.Edit = true
	r.Update = true
	r.Delete = true
	r.Activate = true
	r.Deactivate = true
	r.Action = true
}
