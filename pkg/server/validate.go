package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/oneblock-ai/apiserver/v2/pkg/apierror"
	"github.com/oneblock-ai/apiserver/v2/pkg/parse"
	"github.com/oneblock-ai/apiserver/v2/pkg/types"

	"github.com/rancher/wrangler/v2/pkg/schemas"
	"github.com/rancher/wrangler/v2/pkg/schemas/validation"
)

const (
	csrfCookie = "CSRF"
	csrfHeader = "X-API-CSRF"
)

func ValidateAction(request *types.APIRequest) (*schemas.Action, error) {
	if request.Action == "" || request.Link != "" || request.Method != http.MethodPost {
		return nil, nil
	}

	if err := request.AccessControl.CanAction(request, request.Schema, request.Action); err != nil {
		return nil, err
	}

	actions := request.Schema.CollectionActions
	if request.Name != "" {
		actions = request.Schema.ResourceActions
	}

	action, ok := actions[request.Action]
	if !ok {
		return nil, apierror.NewAPIError(validation.InvalidAction, fmt.Sprintf("Invalid action: %s", request.Action))
	}

	return &action, nil
}

func CheckCSRF(apiOp *types.APIRequest) error {
	if !parse.IsBrowser(apiOp.Request, false) {
		return nil
	}

	cookie, err := apiOp.Request.Cookie(csrfCookie)
	if errors.Is(err, http.ErrNoCookie) {
		// 16 bytes = 32 Hex Char = 128 bit entropy
		bytes := make([]byte, 16)
		_, err := rand.Read(bytes)
		if err != nil {
			return apierror.WrapAPIError(err, validation.ServerError, "Failed in CSRF processing")
		}

		cookie = &http.Cookie{
			Name:   csrfCookie,
			Value:  hex.EncodeToString(bytes),
			Path:   "/",
			Secure: true,
		}

		http.SetCookie(apiOp.Response, cookie)
	} else if err != nil {
		return apierror.NewAPIError(validation.InvalidCSRFToken, "Failed to parse cookies")
	} else if apiOp.Method != http.MethodGet {
		/*
		 * Very important to use apiOp.Method and not apiOp.Request.Method. The client can override the HTTP method with _method
		 */
		if cookie.Value == apiOp.Request.Header.Get(csrfHeader) {
			// Good
		} else if cookie.Value == apiOp.Request.URL.Query().Get(csrfCookie) {
			// Good
		} else {
			return apierror.NewAPIError(validation.InvalidCSRFToken, "Invalid CSRF token")
		}
	}

	return nil
}
