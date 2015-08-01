//go:generate nanotemplate -T string --input=result.go.t
//go:generate nanotemplate -T int --input=result.go.t
package result

import (
	"errors"
	"github.com/nanoservice/monad.go/result/result_int"
	"github.com/nanoservice/monad.go/result/result_string"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringExample(t *testing.T) {
	helloFn := func(name string) result_string.Result {
		return result_string.Success("hello, " + name)
	}

	success := result_string.Success("world").Bind(helloFn)
	assert.Equal(t, result_string.Success("hello, world"), success)

	err := errors.New("The error")
	failure := result_string.Failure(err).Bind(helloFn)
	assert.Equal(t, result_string.Failure(err), failure)
}

func TestIntExample(t *testing.T) {
	addTwo := func(x int) result_int.Result {
		return result_int.Success(2 + x)
	}

	success := result_int.Success(7).Bind(addTwo)
	assert.Equal(t, result_int.Success(9), success)

	err := errors.New("The error")
	failure := result_int.Failure(err).Bind(addTwo)
	assert.Equal(t, result_int.Failure(err), failure)
}

func TestOnErrorFn(t *testing.T) {
	var called bool
	var got error
	var r result_string.Result

	called = false
	r = result_string.Success("yep!").
		OnErrorFn(func(e error) { called = true })
	assert.Equal(t, false, called)
	assert.Equal(t, result_string.Success("yep!"), r)

	called = false
	err := errors.New("The error")
	r = result_string.Failure(err).
		OnErrorFn(func(e error) { called = true; got = e })
	assert.Equal(t, true, called)
	assert.Equal(t, err, got)
	assert.Equal(t, result_string.Failure(err), r)
}