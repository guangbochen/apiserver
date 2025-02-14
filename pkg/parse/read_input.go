package parse

import (
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/rancher/wrangler/v2/pkg/data/convert"
	"github.com/rancher/wrangler/v2/pkg/schemas/validation"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/oneblock-ai/apiserver/v2/pkg/apierror"
	"github.com/oneblock-ai/apiserver/v2/pkg/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const reqMaxSize = (2 * 1 << 20) + 1

var bodyMethods = map[string]bool{
	http.MethodPut:  true,
	http.MethodPost: true,
}

type Decode func(interface{}) error

func ReadBody(req *http.Request) (types.APIObject, error) {
	if !bodyMethods[req.Method] {
		return types.APIObject{}, nil
	}

	decode := getDecoder(req, io.LimitReader(req.Body, maxFormSize))

	data := map[string]interface{}{}
	if err := decode(&data); err != nil {
		return types.APIObject{}, apierror.NewAPIError(validation.InvalidBodyContent,
			fmt.Sprintf("Failed to parse body: %v", err))
	}

	return toAPI(data), nil
}

func toAPI(data map[string]interface{}) types.APIObject {
	return types.APIObject{
		Type:   convert.ToString(data["type"]),
		ID:     convert.ToString(data["id"]),
		Object: data,
	}
}

func getDecoder(req *http.Request, reader io.Reader) Decode {
	if req.Header.Get("Content-type") == "application/yaml" {
		return yaml.NewYAMLToJSONDecoder(reader).Decode
	}
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	return decoder.Decode
}
