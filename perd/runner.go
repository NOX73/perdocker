package perd

type Runner interface {
  RunWorker()
  Eval(string) Result
}

type runner struct {}

