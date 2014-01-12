package perd

type Result interface {
  Bytes() []byte
}

func NewResult(res string) Result {
  return &result{res, "", 0}
}

type result struct {
  stdOut      string
  stdErr      string
  statusCode  int
}

func (r *result) Bytes () []byte {
  return []byte(r.stdOut)
}
