package apiprivate_test

import (
	"reflect"
	"testing"

	protoCommon "somerepo/common"
	"somerepo/test"
	protoUserCommon "somerepo/user_common"
	protoUserPrivate "somerepo/user_private"

	"github.com/golang/mock/gomock"
)

func TestContactVerified(t *testing.T) {
	t.Parallel()

	api, sut := test.NewPrivateSut(t)

	t.Cleanup(func() {
		sut.Ctrl.Finish()
	})

	tt := []struct {
		name         string
		request      *protoUserPrivate.ContactVerifiedRequest
		expectations func(test.Sut)
		wantResponse *protoUserPrivate.ContactVerifiedResponse
	}{
		{
			name: "0. Invalid request. Empty user ID",
			request: &protoUserPrivate.ContactVerifiedRequest{
				UserId:  "",
				Contact: "contact_1",
				Type:    protoUserCommon.ContactType_CONTACT_PHONE,
			},
			expectations: func(_ test.Sut) {},
			wantResponse: &protoUserPrivate.ContactVerifiedResponse{
				Success: false,
				Error: &protoCommon.ResponseError{
					Code:        "USER-E010",
					TransId:     "",
					Message:     "UserID is empty",
					Params:      nil,
					ParentError: nil,
				},
			},
		},
		{
			name: "1. Invalid request. Invalid user ID",
			request: &protoUserPrivate.ContactVerifiedRequest{
				UserId:  "550e8400-e29b-41d4-a716-44665544000.",
				Contact: "contact_2",
				Type:    protoUserCommon.ContactType_CONTACT_PHONE,
			},
			expectations: func(_ test.Sut) {},
			wantResponse: &protoUserPrivate.ContactVerifiedResponse{
				Success: false,
				Error: &protoCommon.ResponseError{
					Code:    "USER-E011",
					TransId: "",
					Message: "User ID \"{{orig}}\" format is not UUID",
					Params: map[string]string{
						"orig": "550e8400-e29b-41d4-a716-44665544000.",
					},
					ParentError: nil,
				},
			},
		},
		{
			name: "2. Invalid request. Empty contact",
			request: &protoUserPrivate.ContactVerifiedRequest{
				UserId:  "550e8400-e29b-41d4-a716-446655440000",
				Contact: "",
				Type:    protoUserCommon.ContactType_CONTACT_PHONE,
			},
			expectations: func(_ test.Sut) {},
			wantResponse: &protoUserPrivate.ContactVerifiedResponse{
				Success: false,
				Error: &protoCommon.ResponseError{
					Code:        "USER-E005",
					TransId:     "",
					Message:     "Contact is empty",
					Params:      nil,
					ParentError: nil,
				},
			},
		},
		{
			name: "3. Invalid request. Type Email",
			request: &protoUserPrivate.ContactVerifiedRequest{
				UserId:  "550e8400-e29b-41d4-a716-446655440000",
				Contact: "contact_4",
				Type:    protoUserCommon.ContactType_CONTACT_EMAIL,
			},
			expectations: func(_ test.Sut) {},
			wantResponse: &protoUserPrivate.ContactVerifiedResponse{
				Success: false,
				Error: &protoCommon.ResponseError{
					Code:        "USER-E006",
					TransId:     "",
					Message:     "Email \"{{orig}}\" is incorrect",
					Params:      map[string]string{"orig": "contact_4"},
					ParentError: nil,
				},
			},
		},
		{
			name: "4. Valid request. No reposioty error",
			request: &protoUserPrivate.ContactVerifiedRequest{
				UserId:  "550e8400-e29b-41d4-a716-446655440000",
				Contact: "contact_4",
				Type:    protoUserCommon.ContactType_CONTACT_PHONE,
			},
			expectations: func(sut test.Sut) {
				sut.MockRepository.EXPECT().ContactVerified(
					gomock.Any(),
					"550e8400-e29b-41d4-a716-446655440000",
					"contact_4",
					protoUserCommon.ContactType_CONTACT_PHONE).
					Return(nil).
					Times(1)
			},
			wantResponse: &protoUserPrivate.ContactVerifiedResponse{
				Success: true,
			},
		},
		{
			name: "5. Valid request. Reposioty error",
			request: &protoUserPrivate.ContactVerifiedRequest{
				UserId:  "550e8400-e29b-41d4-a716-446655440000",
				Contact: "contact_4",
				Type:    protoUserCommon.ContactType_CONTACT_PHONE,
			},
			expectations: func(sut test.Sut) {
				sut.MockRepository.EXPECT().ContactVerified(
					gomock.Any(),
					"550e8400-e29b-41d4-a716-446655440000",
					"contact_4",
					protoUserCommon.ContactType_CONTACT_PHONE).
					Return(errAny).
					Times(1)
			},
			wantResponse: &protoUserPrivate.ContactVerifiedResponse{
				Success: false,
				Error: &protoCommon.ResponseError{
					Code:        "USER-E004",
					TransId:     "",
					Message:     "Internal error. Please try again later",
					Params:      nil,
					ParentError: nil,
				},
			},
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()

			// Expectations
			v.expectations(sut)

			// ContactVerified
			gotResponse, _ := api.ContactVerified(t.Context(), v.request)

			// Check response
			if !reflect.DeepEqual(v.wantResponse, gotResponse) {
				t.Errorf("unexpected response: got %v, want %v", gotResponse, v.wantResponse)
			}
		})
	}
}
