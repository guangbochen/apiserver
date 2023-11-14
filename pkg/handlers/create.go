package handlers

import (
	"github.com/oneblock-ai/apiserver/v2/pkg/apierror"
	"github.com/oneblock-ai/apiserver/v2/pkg/parse"
	"github.com/oneblock-ai/apiserver/v2/pkg/types"
	"github.com/rancher/wrangler/v2/pkg/schemas/validation"
)

func CreateHandler(apiOp *types.APIRequest) (types.APIObject, error) {
	var err error

	if err := apiOp.AccessControl.CanCreate(apiOp, apiOp.Schema); err != nil {
		return types.APIObject{}, err
	}

	data, err := parse.Body(apiOp.Request)
	if err != nil {
		return types.APIObject{}, err
	}

	store := apiOp.Schema.Store
	if store == nil {
		return types.APIObject{}, apierror.NewAPIError(validation.NotFound, "no store found")
	}

	data, err = store.Create(apiOp, apiOp.Schema, data)
	if err != nil {
		return types.APIObject{}, err
	}

	return data, nil
}
