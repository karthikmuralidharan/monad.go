// result monad
// generated by github.com/nanoservice/monad.go/result
// type: string
package result_string

import ()

type Result struct {
	value *string
	err   error
}

func NewResult(value string, err error) Result {
  return Result{value: &value, err: err}
}

func Success(value string) Result {
	return Result{value: &value, err: nil}
}

func Failure(err error) Result {
	return Result{value: nil, err: err}
}

func (r Result) Bind(fn func(string) Result) Result {
	if r.err != nil {
		return r
	}
	return fn(*r.value)
}

func (r Result) OnErrorFn(fn func(error)) Result {
  if r.err != nil {
    fn(r.err)
  }
  return r
}