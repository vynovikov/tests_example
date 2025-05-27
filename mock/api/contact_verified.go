package apiprivate

import (
	"context"

	"somerepo/helpers"
	protoUserCommon "somerepo/user_common"
	protoUserPrivate "somerepo/user_private"

	"github.com/google/uuid"
)

func (api API) ContactVerified(
	ctx context.Context,
	request *protoUserPrivate.ContactVerifiedRequest,
) (*protoUserPrivate.ContactVerifiedResponse, error) {
	earlyExit := func(code uint16, message string, params map[string]string) (*protoUserPrivate.ContactVerifiedResponse, error) {
		return &protoUserPrivate.ContactVerifiedResponse{
			Success: false,
			Error:   api.earlyError(code, message, params),
		}, nil
	}

	userID, contact, contactType := request.GetUserId(), request.GetContact(), request.GetType()

	if userID == "" {
		return earlyExit(10, "UserID is empty", nil)
	}

	err := uuid.Validate(userID)
	if err != nil {
		return earlyExit(11, "User ID \"{{orig}}\" format is not UUID", map[string]string{"orig": userID})
	}

	if contact == "" {
		return earlyExit(5, "Contact is empty", nil)
	}

	if contactType == protoUserCommon.ContactType_CONTACT_EMAIL && !helpers.EmailIsValid(contact) {
		return earlyExit(6, "Email \"{{orig}}\" is incorrect", map[string]string{"orig": contact})
	}

	err = api.domain.ContactVerified(ctx, userID, contact, contactType)

	responseError := api.responseError(err)

	return &protoUserPrivate.ContactVerifiedResponse{
		Success: responseError == nil,
		Error:   responseError,
	}, nil
}
