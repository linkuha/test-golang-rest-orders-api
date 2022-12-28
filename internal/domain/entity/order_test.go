package entity_test

import (
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrderValidateOK(t *testing.T) {
	cases := []struct {
		name   string
		in     *entity.Order
		expErr error
	}{
		{
			name: "ok_without_id",
			in: &entity.Order{
				Number: 1,
				UserID: "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
			},
		},
		{
			name: "ok_full",
			in: &entity.Order{
				ID:     "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
				Number: 1,
				UserID: "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
			},
		},
	}

	for _, tCase := range cases {
		err := tCase.in.Validate()
		require.NoError(t, err)
	}
}

func TestOrderValidateError(t *testing.T) {
	cases := []struct {
		name   string
		in     *entity.Order
		expErr error
	}{
		{
			name: "empty_user_id",
			in:   &entity.Order{Number: 1},
		},
		{
			name: "bad_user_id",
			in: &entity.Order{
				ID:     "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
				Number: 1,
				UserID: "c401f9dc-1e68-4b44asd93b09e3fe7",
			},
		},
		{
			name: "bad_id",
			in: &entity.Order{
				ID:     "c401fdafe7",
				Number: 1,
				UserID: "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
			},
		},
		{
			name: "empty_number",
			in:   &entity.Order{Number: 0, UserID: "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"},
		},
	}

	for _, tCase := range cases {
		err := tCase.in.Validate()
		require.Error(t, err)
	}
}
