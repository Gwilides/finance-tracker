package req

import (
	"net/http"

	"github.com/Gwilides/finance-tracker/pkg/res"
)

func HandleBody[T any](w http.ResponseWriter, req *http.Request) (*T, error) {
	body, err := Decode[T](req.Body)
	if err != nil {
		res.Json(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	err = IsValid(body)
	if err != nil {
		res.Json(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &body, nil

}
