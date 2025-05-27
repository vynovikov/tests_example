package storage

import (
	"context"

	protoUserCommon "somerepo/user_common"
	"somerepo/wrapper/v2"
)

const (
	insertContactQuery = `INSERT INTO user_contacts
								(user_id, contact, contact_type, use_for_authorization, is_primary)
						  VALUES
    							($1, $2, $3, false, false)`
)

func (r Repository) ContactAdd(
	ctx context.Context,
	userID string,
	contact string,
	contactType protoUserCommon.ContactType,
) error {
	_, err := r.client.ExecContext(
		ctx,
		insertContactQuery,
		userID,
		contact,
		contactType.ToValue(),
	)

	return wrapper.Wrap(err)
}
