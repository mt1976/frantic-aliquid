package auth

import (
	"context"
	"net/http"
	"strings"

	secHelper "github.com/mt1976/frantic-aegis/app/web/security"
	"github.com/mt1976/frantic-aliquid/app/business/followOnPermissions"
	"github.com/mt1976/frantic-aliquid/app/dao/authorityStore"
	"github.com/mt1976/frantic-aliquid/app/dao/behaviourStore"
	"github.com/mt1976/frantic-core/contextHandler"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
)

//var domain = domains.SECURITY

func CheckUserAuthority(ctx context.Context, w http.ResponseWriter, r *http.Request, b behaviourStore.Behaviour_Store) bool {
	// Check if the user has the required permissions to perform the action
	// If the user has the required permissions, return true

	userKey := contextHandler.GetSession_UserKey(ctx)

	//tempID := ctx.Value(cfg.GetSecuritySessionKey_UserKey())
	// if userKey == "" {
	// 	logHandler.SecurityLogger.Printf("[%v] Error getting user ID: %v", strings.ToUpper(domain.String()), "error")
	// 	return false
	// }
	//logHandler.EventLogger.Printf("tempID: %v %v", tempID, reflect.TypeOf(tempID))

	//userKey := tempID.(string)

	// userMessage := messageHelpers.UserMessage{
	// 	Key:      userKey,
	// 	Code:     contextHandler.GetSession_UserCode(ctx),
	// 	Locale:   contextHandler.GetSession_Locale(ctx),
	// 	Theme:    contextHandler.GetSession_Theme(ctx),
	// 	Timezone: contextHandler.GetSession_Timezone(ctx),
	// }
	userMessage := secHelper.UserMessageFromContext(ctx)

	//uID, err := strconv.Atoi(id.(string))
	if userKey == "" {
		logHandler.SecurityLogger.Printf("[%v] Error getting user Key: %v", strings.ToUpper(domain.String()), "error")
		return false
	}
	auth, err := GetAuthorities(ctx, userMessage, b.Key)
	if err != nil {
		logHandler.SecurityLogger.Printf("[%v] %v", strings.ToUpper(domain.String()), err)
		return false
	}
	////spew.Dump(auth)
	if auth.Key != "" {
		logHandler.SecurityLogger.Printf("[%v] User [%v] has required authority for behaviour [%v] [%v]", appName, userKey, b.Note, b.Raw)
	} else {
		logHandler.SecurityLogger.Printf("[%v] User [%v] does not have required authority for behaviour [%v] [%v]", appName, userKey, b.Note, b.Raw)
		//Violation(w, r, "You do not have required permissions for behaviour: "+b.Text)
	}

	return true
}

func GetRights(ctx context.Context, b behaviourStore.Behaviour_Store) followOnPermissions.FollowOnPermissions {

	//TODO: Sort this out its a mess!!!!!

	// Get the rights of the user for the given behaviour
	// If the user has the required permissions, return the rights
	// If the user does not have the required permissions, return an empty string

	// fmt.Printf("userCode: %v\n", userCode)
	// fmt.Printf("b: %+v\n", b)
	userCode := cfg.GetServiceUser_UserCode()

	logHandler.SecurityLogger.Printf("[%v] Get Right for userCode=[%v], behaviour=[%v]", strings.ToUpper(domain.String()), userCode, b.Raw)

	// u, err := userStore.GetByUserCode(userCode)
	// if err != nil {
	// 	logHandler.SecurityLogger.Printf("[%v] Error getting user: %v", strings.ToUpper(domain.String()), err)
	// 	return permissions.Rights{}
	// }

	um := messageHelpers.UserMessage{
		Code:   userCode,
		Source: cfg.GetApplication_Name(),
	}

	a, err := GetAuthorities(ctx, um, b.Key)
	if err != nil {
		logHandler.SecurityLogger.Printf("[%v] %v", strings.ToUpper(domain.String()), err)
		return followOnPermissions.FollowOnPermissions{}
	}
	//logHandler.InfoLogger.Printf("****** a: %+v", a)
	authority, err := authorityStore.GetByKey(a.Key)
	if err != nil {
		logHandler.SecurityLogger.Printf("[%v] %v", strings.ToUpper(domain.String()), err)
		return followOnPermissions.FollowOnPermissions{}
	}

	logHandler.SecurityLogger.Printf("[%v] Rights Found for userCode=[%v], behaviour=[%v], rights=[%+v]", strings.ToUpper(domain.String()), userCode, b.Raw, authority.Behaviour.Permissions)

	return authority.Behaviour.Permissions
}

func GetAuthorities(ctx context.Context, usr messageHelpers.UserMessage, behaviorKey string) (messageHelpers.AuthorityMessage, error) {

	userKey := usr.Code
	// Get all status
	authorityList, err := authorityStore.GetAll()
	if err != nil {
		logHandler.SecurityLogger.Printf("ERROR Getting all status: %v", err)
		return messageHelpers.AuthorityMessage{}, err
	}
	// //translate the user into a user key
	// u, err := userStore.GetBy(userStore.FIELD_Key, userKey)
	// if err != nil {
	// 	logger.SecurityLogger.Printf("ERROR Getting user: %v", err)
	// 	return authorityStore.Authority_Store{}, err
	// }
	//if err != nil {

	// range through status list, if status code is found and deletedby is empty then return error
	for _, a := range authorityList {
		//logHandler.InfoLogger.Printf("** a: %v %v %v %v", a.UserCode, usr.Code, a.Behaviour.Key, behaviorKey)
		if a.UserCode == userKey && a.Behaviour.Key == behaviorKey {
			//logHandler.InfoLogger.Printf("** a: %v", a)
			return a.BuildMessage()
		}
	}

	msg := "[%v] Authority not found for user [%v] and behaviour [%v]"

	logHandler.SecurityLogger.Printf(msg, strings.ToUpper(domain.String()), userKey, behaviorKey)
	logHandler.WarningLogger.Printf(msg, strings.ToUpper(domain.String()), userKey, behaviorKey)
	return messageHelpers.AuthorityMessage{}, nil
}
