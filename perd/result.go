package perd

type Result interface {
  Bytes() []byte
}

func NewResult(out, err string, code int) Result {
  return &result{out, err, code}
}

type result struct {
  stdOut      string
  stdErr      string
  statusCode  int
}

func (r *result) Bytes () []byte {
  return []byte(r.stdOut)
}
