package auth

import (
	"context"
	"strings"

	"github.com/mt1976/frantic-aliquid/app/business/domains"
	"github.com/mt1976/frantic-aliquid/app/business/translation"
	"github.com/mt1976/frantic-aliquid/app/dao/authorityStore"
	"github.com/mt1976/frantic-aliquid/app/dao/behaviourStore"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
)

var domain = domains.AUTHORITY
var cfg commonConfig.Settings

func init() {
	cfg = *commonConfig.Get()
}

func encode(str string) string {
	// Encode the string
	// Return the encoded string
	return idHelpers.Encode(str)
}

// func GrantUserAuthority(ctx context.Context, usr messageHelpers.UserMessage, ben behaviourStore.Behaviour_Store) (authorityStore.Authority_Store, error) {
// 	user, userError := userStore.GetByUserCode(usr.Code)
// 	if userError != nil {
// 		logHandler.ErrorLogger.Println(userError.Error())
// 		return authorityStore.Authority_Store{}, userError
// 	}
// 	return grantUserAuthority(ctx, user, ben)
// }

func GrantUserAuthority(ctx context.Context, usr messageHelpers.UserMessage, ben behaviourStore.Behaviour_Store) (messageHelpers.AuthorityMessage, error) {

	g, err := authorityStore.BuildBaseRecord(usr, ben)
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] New %s %v %v", strings.ToUpper(domain.String()), err.Error(), g.UserKey, g.Behaviour)
		return messageHelpers.AuthorityMessage{}, err
	}

	xUser, err := g.IsDuplicateOf(g.Key)
	if err == commonErrors.ErrorDuplicate {
		logHandler.WarningLogger.Printf("[%v] DUPLICATE [%v] already exists,skipping", domain.String(), g.Raw)
		rtn, _ := xUser.BuildMessage()
		return rtn, nil
	}

	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] New %s %v %v", strings.ToUpper(domain.String()), err.Error(), g.UserKey, g.Behaviour)
		return messageHelpers.AuthorityMessage{}, err
	}

	if _, err := g.Create(ctx, audit.GRANT, translation.Get("Authority Granted")); err != nil {
		logHandler.ErrorLogger.Printf("[%v] Create %s %v %v", strings.ToUpper(domain.String()), err.Error(), g.UserKey, g.Behaviour)
		panic(err)
	}

	//logHandler.AuditLogger.Printf("[ACTIVITIES] [%v] [%v] ID=[%v] Notes[%v]", audit.CREATE.Text(), strings.ToUpper(domain.String()), g.UserKey, fmt.Sprintf("New %v %v", g.UserKey, g.Behaviour))
	logHandler.SecurityLogger.Printf("Granted '%v' '%v' authority", usr.Code, ben.Raw)
	return g.BuildMessage()
}

func RevokeUserAuthority(ctx context.Context, usr messageHelpers.UserMessage, ben behaviourStore.Behaviour_Store) error {
	// user, userError := userStore.GetByUserCode(usr.Code)
	// if userError != nil {
	// 	logHandler.ErrorLogger.Println(userError.Error())
	// 	return userError
	// }
	return revokeUserAuthority(ctx, usr, ben)
}

func revokeUserAuthority(ctx context.Context, usr messageHelpers.UserMessage, ben behaviourStore.Behaviour_Store) error {
	userErr := usr.Validate(logHandler.SecurityLogger)
	if userErr != nil {
		return userErr
	}
	user := usr.Key
	action := ben.Key
	authority_Key := strings.ToLower(user + cfg.SEP() + action)
	authority_Key = idHelpers.Encode(authority_Key)
	err := authorityStore.DeleteBy(ctx, authorityStore.FIELD_Key, authority_Key, "Revoked")
	if err != nil {
		logHandler.WarningLogger.Printf("[%v] Reading Id=[%v] %v", strings.ToUpper(domain.String()), authority_Key, err.Error())
		return err
	}
	logHandler.SecurityLogger.Printf("Revoked '%v' '%v' authority", usr.Code, ben.Raw)

	return nil
}
