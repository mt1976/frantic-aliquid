package auth

import (
	"context"
	"strings"

	"github.com/mt1976/frantic-aliquid/app/dao/behaviourStore"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
)

var appName = "frantic-aliquid"

func GrantAllUserAccessRights(usr messageHelpers.UserMessage) error {

	usrErr := usr.Validate(logHandler.SecurityLogger)
	if usrErr != nil {
		logHandler.ErrorLogger.Println(usrErr.Error())
		return usrErr
	}

	if usr.Source == "" {
		logHandler.ErrorLogger.Printf("[%v] User [%v] has no source using [%v]", appName, usr.Code, cfg.GetApplication_Name())
		usr.Source = cfg.GetApplication_Name()
	}

	bList, err := behaviourStore.GetAllWhere(behaviourStore.FIELD_Source, usr.Source)
	if err != nil {
		logHandler.ErrorLogger.Println(err.Error())
		return err
	}

	for _, b := range bList {
		if strings.EqualFold(b.Source, usr.Source) {
			logHandler.InfoLogger.Printf("[%v] Granting [%v] to [%v]", appName, b.Raw, usr.Code)
			_, err = GrantUserAuthority(context.TODO(), usr, b)
			if err != nil {
				logHandler.ErrorLogger.Println(err.Error())
				return err
			}
		}
	}

	return nil
}
