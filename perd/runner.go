package perd

type Runner interface {
  Run()
  Eval(string) Result 
}
