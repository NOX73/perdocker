package perd

type Runner interface {
  RunWorker()
  RunWorkers(int)
  Eval(string) Result
}

type runner struct {}

