// result monad
// generated by github.com/nanoservice/monad.go/result
// type: *http.Response
package result_response

import ("net/http")

type handler           func(*http.Response) Result
type errorHandler      func(error)
type deferHandler      func()
type boundDeferHandler func(*http.Response)

type Result struct {
        value         **http.Response
        err           error
        deferHandlers []deferHandler
}

func NewResult(value *http.Response, err error) Result {
        return buildResult(&value, err)
}

func Success(value *http.Response) Result {
        return buildResult(&value, nil)
}

func Failure(err error) Result {
        return buildResult(nil, err)
}

func (r Result) Bind(fn handler) Result {
        if r.err != nil {
          return r
        }

        result := fn(*r.value)
        return r.augment(result.value, result.err)
}

func (r Result) Defer(fn boundDeferHandler) Result {
        if r.err != nil {
                return r
        }

        return Result{
                value:         r.value,
                err:           r.err,
                deferHandlers: append(
                        r.deferHandlers,
                        func() { fn(*r.value) },
                ),
        }
}

func (r Result) Err() error {
        for _, fn := range r.deferHandlers {
                fn()
        }
        return r.err
}

func (r Result) Chain(fns... handler) Result {
        for _, fn := range fns {
                r = r.Bind(fn)
        }
        return r
}

func (r Result) OnErrorFn(fn errorHandler) Result {
        if r.err != nil {
                fn(r.err)
        }
        return r
}

func (r Result) augment(value **http.Response, err error) (result Result) {
        result = buildResult(value, err)
        result.deferHandlers = r.deferHandlers
        return
}

func buildResult(value **http.Response, err error) Result {
        return Result{
                value:         value,
                err:           err,
                deferHandlers: []deferHandler{},
        }
}
