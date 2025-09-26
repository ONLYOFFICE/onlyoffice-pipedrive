package response

import "encoding/json"

type DataResponse struct {
	Data any    `json:"data"`
	Code string `json:"code"`
}

func (d DataResponse) ToJSON() []byte {
	json, err := json.Marshal(d)
	if err != nil {
		return nil
	}

	return json
}
