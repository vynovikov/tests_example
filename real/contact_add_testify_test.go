//go:build postgres

package storage_test

import (
	"time"

	"somerepo/test"
	protoUserCommon "somerepo/user_common"
)

func (s *contactSuite) TestContactAddTestify() { //nolint:maintidx
	s.T().Parallel()
	ctx, repository, postgresClient := test.NewRepository(s.T())

	s.T().Cleanup(func() {
		postgresClient.Close()
	})

	tt := []struct { //nolint:dupl
		name                string
		initialUsers        []userModel
		initialData         []contactModel
		userID              string
		contact             string
		contactType         protoUserCommon.ContactType
		wantData            []contactModel
		isCreatedAtRemained bool
		isDuplicateError    bool
	}{
		{
			name: "0. No data",
			initialUsers: []userModel{
				{
					ID:    "123e4567-e89b-12d3-a456-426614174000",
					Email: "email_0",
				},
			},
			initialData: []contactModel{},
			userID:      "123e4567-e89b-12d3-a456-426614174000",
			contact:     "contact_0",
			contactType: protoUserCommon.ContactType_CONTACT_EMAIL,
			wantData: []contactModel{
				{
					UserID:      "123e4567-e89b-12d3-a456-426614174000",
					Contact:     "contact_0",
					ContactType: "email",
				},
			},
			isCreatedAtRemained: false,
			isDuplicateError:    false,
		},
		{
			name: "1. Have data. New contact. New userID",
			initialUsers: []userModel{
				{
					ID:    "123e4567-e89b-12d3-a456-426614174000",
					Email: "email_10",
				},
				{
					ID:    "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Email: "email_11",
				},
			},
			initialData: []contactModel{
				{
					UserID:              "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:             "contact_10",
					ContactType:         "email",
					UseForAuthorization: false,
					DeletedAt:           mustParseTime("0001-01-01 00:00:00 +0000"),
				},
			},
			userID:      "123e4567-e89b-12d3-a456-426614174000",
			contact:     "contact_11",
			contactType: protoUserCommon.ContactType_CONTACT_PHONE,
			wantData: []contactModel{
				{
					UserID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:     "contact_10",
					ContactType: "email",
				},
				{
					UserID:      "123e4567-e89b-12d3-a456-426614174000",
					Contact:     "contact_11",
					ContactType: "phone",
				},
			},
			isCreatedAtRemained: false,
			isDuplicateError:    false,
		},
		{
			name: "2. Have data. New contact. Existing userID. New contact",
			initialUsers: []userModel{
				{
					ID:    "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Email: "email_20",
				},
			},
			initialData: []contactModel{
				{
					UserID:              "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:             "contact_20",
					ContactType:         "email",
					UseForAuthorization: false,
					DeletedAt:           mustParseTime("0001-01-01 00:00:00 +0000"),
				},
			},
			userID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
			contact:     "contact_21",
			contactType: protoUserCommon.ContactType_CONTACT_PHONE,
			wantData: []contactModel{
				{
					UserID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:     "contact_20",
					ContactType: "email",
				},
				{
					UserID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:     "contact_21",
					ContactType: "phone",
				},
			},
			isCreatedAtRemained: false,
			isDuplicateError:    false,
		},
		{
			name: "3. Have data. New contact. Existing userID. Existing contact. New contact_type",
			initialUsers: []userModel{
				{
					ID:    "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Email: "email_30",
				},
			},
			initialData: []contactModel{
				{
					UserID:              "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:             "contact_30",
					ContactType:         "email",
					UseForAuthorization: false,
					DeletedAt:           mustParseTime("0001-01-01 00:00:00 +0000"),
				},
			},
			userID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
			contact:     "contact_30",
			contactType: protoUserCommon.ContactType_CONTACT_PHONE,
			wantData: []contactModel{
				{
					UserID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:     "contact_30",
					ContactType: "email",
				},
				{
					UserID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:     "contact_30",
					ContactType: "phone",
				},
			},
			isCreatedAtRemained: false,
			isDuplicateError:    false,
		},
		{
			name: "4. Have data. Old contact has not zero deleed_at",
			initialUsers: []userModel{
				{
					ID:    "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Email: "email_40",
				},
			},
			initialData: []contactModel{
				{
					UserID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:     "contact_40",
					ContactType: "telegram",
					DeletedAt:   mustParseTime("0001-01-01 00:00:00 +0000"),
				},
			},
			userID:           "abcdefab-cdef-abcd-efab-cdefabcdefab",
			contact:          "contact_40",
			contactType:      protoUserCommon.ContactType_CONTACT_TELEGRAM,
			isDuplicateError: true,
		},
		{
			name: "5. Have data. Old contact has zero deleed_at",
			initialUsers: []userModel{
				{
					ID:    "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Email: "email_50",
				},
			},
			initialData: []contactModel{
				{
					UserID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:     "contact_40",
					ContactType: "telegram",
					DeletedAt:   mustParseTime("0001-01-01 00:00:00 +0000"),
				},
			},
			userID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
			contact:     "contact_40",
			contactType: protoUserCommon.ContactType_CONTACT_TELEGRAM,
			wantData: []contactModel{
				{
					UserID:      "abcdefab-cdef-abcd-efab-cdefabcdefab",
					Contact:     "contact_40",
					ContactType: "telegram",
					DeletedAt:   mustParseTime("0001-01-01 00:00:00 +0000"),
				},
			},
			isCreatedAtRemained: true,
			isDuplicateError:    true,
		},
	}

	for _, v := range tt {
		s.Run(v.name, func() {
			// Parallel
			s.T().Parallel()

			// Clear data
			_, err := postgresClient.ExecContext(ctx, truncateUsersQuery)
			if err != nil {
				s.FailNow(err.Error())
			}

			// Insert initial data
			for _, user := range v.initialUsers {
				_, err = postgresClient.ExecContext(
					ctx,
					insertUserQuery,
					user.ID,
					user.Email,
				)
				if err != nil {
					s.FailNow(err.Error())
				}
			}

			for _, data := range v.initialData {
				_, err = postgresClient.ExecContext(ctx,
					insertContactQuery,
					data.UserID,
					data.Contact,
					data.ContactType,
					data.UseForAuthorization,
					data.DeletedAt)
				if err != nil {
					s.FailNow(err.Error())
				}
			}

			// Check createdAt
			var createdAt time.Time
			if v.isCreatedAtRemained {
				err = postgresClient.GetContext(ctx,
					&createdAt,
					`SELECT created_at FROM user_contacts WHERE contact = $1`,
					v.contact)
				if err != nil {
					s.FailNow(err.Error())
				}
			}

			// ContactAdd
			err = repository.ContactAdd(ctx, v.userID, v.contact, v.contactType)
			if v.isDuplicateError {
				s.Error(err)

				return
			}

			// Get data
			gotData := make([]contactModel, 0)

			err = postgresClient.SelectContext(ctx, &gotData, selectContactQuery)
			if err != nil {
				s.FailNow(err.Error())
			}

			// Check result
			s.Len(gotData, len(v.wantData))

			for gotIndex, gotValue := range gotData {
				s.Equal(v.wantData[gotIndex].UserID, gotValue.UserID)
				s.Equal(v.wantData[gotIndex].Contact, gotValue.Contact)
				s.Equal(v.wantData[gotIndex].ContactType, gotValue.ContactType)
				s.Equal(v.wantData[gotIndex].UseForAuthorization, gotValue.UseForAuthorization)
				s.Equal(v.wantData[gotIndex].IsPrimary, gotValue.IsPrimary)

				if v.isCreatedAtRemained {
					s.LessOrEqual(gotValue.CreatedAt.Sub(createdAt), time.Second)
				}
			}
		})
	}
}
