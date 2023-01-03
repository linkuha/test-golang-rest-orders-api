package entity_test

import (
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProductValidateOK(t *testing.T) {
	cases := []struct {
		name   string
		in     *entity.Product
		expErr error
	}{
		{
			name: "ok_without_id",
			in:   &entity.Product{Name: "qwerty"},
		},
	}

	for _, tCase := range cases {
		err := tCase.in.Validate()
		require.NoError(t, err)
	}
}

func TestProductValidateError(t *testing.T) {
	cases := []struct {
		name   string
		in     *entity.Product
		expErr error
	}{}

	for _, tCase := range cases {
		err := tCase.in.Validate()
		require.Error(t, err)
	}
}
