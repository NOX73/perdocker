package perd

import "encoding/json"

type Result interface {
	Bytes() []byte
}

func NewResult(out, err []byte, code int) Result {
	return &result{string(out), string(err), code}
}

type result struct {
	StdOut     string `json:"std_out"`
	StdErr     string `json:"std_err"`
	StatusCode int    `json:"code"`
}

func (r *result) Bytes() []byte {
	j, _ := json.Marshal(r)
	return j
}
