package request

import (
	"net/http"
	"strconv"

	"github.com/go-playground/form"
)

type TakeRequest struct {
	Amount *float64 `json:"amount" form:"amount"`
}

func (r *TakeRequest) AssignFormValues(req *http.Request) {
	req.ParseMultipartForm(32 << 20)

	decoder := form.NewDecoder()
	decoder.RegisterCustomTypeFunc(func(vals []string) (any, error) {
		if len(vals) == 0 || vals[0] == "" {
			return nil, nil
		}
		parsed, err := strconv.ParseFloat(vals[0], 64)
		if err != nil {
			return nil, err
		}
		return parsed, nil
	}, float64(0))

	decoder.Decode(&r, req.Form)
}
