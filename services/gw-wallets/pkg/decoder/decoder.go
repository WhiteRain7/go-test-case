package decoder

import (
	"encoding/json"
	"errors"
	"net/http"

	codederrors "github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
)

func Decode(req *http.Request, out any) error {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&out); err != nil {
		return codederrors.New(400, err, "decode error")
	}

	if decoder.More() {
		return codederrors.New(400, errors.New("too many fields detected"), "too many fields")
	}

	return nil
}
