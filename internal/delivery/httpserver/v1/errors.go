package v1

import (
	"errors"
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
	"net/http"
)

var (
	forbiddenError   = errors.New("forbidden")
	emptyParameterID = errors.New("missed identifier param")
)

const (
	ErrServiceInternalText    = "internal server error, contact support with your request id"
	ErrServiceUnavailableText = "temporary server error, try later"
	ErrAuthAPIText            = "invalid API authorization"
	ErrCredentialsText        = "invalid username or password"
	ErrInputJSONText          = "bad input json"
	ErrValidationText         = "validation error"
	ErrNotFoundText           = "resource is not found"
)

type errorHandlingDetails struct {
	ClientError string
	DebugError  string
	Code        int
}

func handleDomainError(e error) errorHandlingDetails {
	var resErr errorHandlingDetails

	if errors.Is(e, forbiddenError) {
		resErr.Code = http.StatusForbidden
		resErr.ClientError = e.Error()
		resErr.DebugError = e.Error()
		return resErr
	}
	if errors.Is(e, emptyParameterID) {
		resErr.Code = http.StatusBadRequest
		resErr.ClientError = e.Error()
		resErr.DebugError = e.Error()
		return resErr
	}

	if customErr, ok := e.(errs.CustomErrorWrapper); ok {
		digErr := customErr.Dig()

		resErr.ClientError = customErr.Message

		switch digErr.Code {
		case errs.Other:
			resErr.Code = http.StatusInternalServerError
		case errs.InvalidOperation:
			resErr.Code = http.StatusMethodNotAllowed
		case errs.InvalidArgument:
			resErr.Code = http.StatusBadRequest
		case errs.MalformedRequest:
			resErr.Code = http.StatusBadRequest
		case errs.IO:
			resErr.Code = http.StatusServiceUnavailable
		case errs.Logic:
			resErr.Code = http.StatusConflict
		case errs.Exist:
			resErr.Code = http.StatusCreated
		case errs.NotExist:
			resErr.Code = http.StatusNotFound
			resErr.ClientError = ErrNotFoundText
		case errs.APIAuthorization:
			resErr.Code = http.StatusUnauthorized
			resErr.ClientError = ErrAuthAPIText
		case errs.UserCredentials:
			resErr.Code = http.StatusOK
			resErr.ClientError = ErrCredentialsText
		case errs.NotPermitted:
			resErr.Code = http.StatusForbidden
		case errs.Private:
			resErr.Code = http.StatusForbidden
		case errs.Internal:
			resErr.Code = http.StatusInternalServerError
		case errs.BrokenLink:
			resErr.Code = http.StatusBadRequest
		case errs.Database:
			resErr.Code = http.StatusInternalServerError
		case errs.DatabaseConnection:
			resErr.Code = http.StatusServiceUnavailable
		case errs.RemoteConnection:
			resErr.Code = http.StatusServiceUnavailable
		case errs.Validation:
			resErr.Code = http.StatusUnprocessableEntity
			resErr.ClientError = fmt.Sprintf("%s: %s", ErrValidationText, digErr.Error())
		case errs.Unanticipated:
			resErr.Code = http.StatusInternalServerError
		}

		resErr.DebugError = digErr.Error()
	}

	switch resErr.Code {
	case http.StatusInternalServerError:
		resErr.ClientError = ErrServiceInternalText
	case http.StatusServiceUnavailable:
		resErr.ClientError = ErrServiceUnavailableText
	}

	return resErr
}

func newJSONBindingErrorWrapper(e error) error {
	return errs.NewErrorWrapper(errs.MalformedRequest, e, ErrInputJSONText)
}
