package error

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBindHelperAlwaysExecutesProvidedBlock(t *testing.T) {
	executed := false
	Bind(func() error {
		executed = true
		return nil
	})

	assert.Equal(t, true, executed)
}

func TestBindHelperReturnsNonFailedErrorIfNil(t *testing.T) {
	e := Bind(func() error { return nil })
	assert.Equal(t, Return(nil), e)
}

func TestBindHelperReturnsWrappedErrorIfFails(t *testing.T) {
	err := errors.New("Something gone wrong")
	e := Bind(func() error { return err })
	assert.Equal(t, Return(err), e)
}

func TestReturnWrapsNil(t *testing.T) {
	e := Return(nil)
	assert.Equal(t, Error{nil, make([]deferrableFunc, 0)}, e)
}

func TestReturnWrapsError(t *testing.T) {
	err := errors.New("Something else have gone wrong")
	e := Return(err)
	assert.Equal(t, Error{err, make([]deferrableFunc, 0)}, e)
}

func TestBindOnNoErrorExecutesProvidedBlock(t *testing.T) {
	executed := false
	Return(nil).Bind(func() error {
		executed = true
		return nil
	})

	assert.Equal(t, true, executed)
}

func TestBindOnNoErrorReturnsNoErrorIfNotFails(t *testing.T) {
	e := Return(nil).Bind(func() error { return nil })
	assert.Equal(t, Return(nil), e)
}

func TestBindOnNoErrorReturnsWrappedErrorIfFails(t *testing.T) {
	err := errors.New("Unable to parse data")
	e := Return(nil).Bind(func() error { return err })
	assert.Equal(t, Return(err), e)
}

func TestBindOnErrorDoesNotExecuteProvidedBlock(t *testing.T) {
	executed := false
	err := errors.New("Incompatible message version")
	Return(err).Bind(func() error {
		executed = true
		return nil
	})

	assert.Equal(t, false, executed)
}

func TestBindOnErrorReturnsSameError(t *testing.T) {
	err := errors.New("Out of imagination to create new error")
	e := Return(err)
	e2 := e.Bind(func() error { return nil })
	assert.Equal(t, e, e2)
}

func TestDeferOnNoErrorDoesNotExecuteProvidedBlock(t *testing.T) {
	executed := false
	Return(nil).Defer(func() { executed = true })
	assert.Equal(t, false, executed)
}

func TestDeferOnNoErrorReturnsSameValue(t *testing.T) {
	e := Return(nil)
	e2 := e.Defer(func() {})
	assert.Equal(t, e.err, e2.err)
}

func TestDeferOnNoErrorAfterErrExecutesProvidedBlock(t *testing.T) {
	executed := false
	Return(nil).Defer(func() { executed = true }).Err()
	assert.Equal(t, true, executed)
}

func TestDeferOnErrorDoesNotExecuteProvidedBlock(t *testing.T) {
	executed := false
	err := errors.New("Yet Another Error Message (YAEM)")
	Return(err).Defer(func() { executed = true })
	assert.Equal(t, false, executed)
}

func TestDeferOnErrorReturnsSameValue(t *testing.T) {
	err := errors.New("Yet Another YAEM")
	e := Return(err)
	e2 := e.Defer(func() {})
	assert.Equal(t, e, e2)
}

func TestDeferOnErrorAfterErrDoesNotExecuteProvidedBlock(t *testing.T) {
	err := errors.New("Yet Another YAEM")
	executed := false
	Return(err).Defer(func() { executed = true }).Err()
	assert.Equal(t, false, executed)
}

func TestDeferOnNoErrorIsPreservedAfterBind(t *testing.T) {
	executed := false
	Return(nil).Defer(
		func() { executed = true },
	).Bind(
		func() error { return nil },
	).Err()

	assert.Equal(t, true, executed)
}

func TestDeferMultiple(t *testing.T) {
	executed := false
	executed2 := false
	Return(nil).Defer(
		func() { executed = true },
	).Defer(
		func() { executed2 = true },
	).Err()

	assert.Equal(t, true, executed)
	assert.Equal(t, true, executed2)
}

func TestDeferIsPreservedAfterOnErrorOnNoError(t *testing.T) {
	executed := false
	Return(nil).Defer(
		func() { executed = true },
	).Bind(
		func() error { return nil },
	).OnError().Err()

	assert.Equal(t, true, executed)
}

func TestDeferIsPreservedAfterOnErrorOnError(t *testing.T) {
	err := errors.New("Boring error message")
	executed := false
	Return(nil).Defer(
		func() { executed = true },
	).Bind(
		func() error { return err },
	).OnError().Err()

	assert.Equal(t, true, executed)
}

func TestErrOnNoErrorReturnsNil(t *testing.T) {
	assert.Equal(t, nil, Return(nil).Err())
}

func TestErrOnErrorReturnsInnerValue(t *testing.T) {
	err := errors.New("Unable to connect to server")
	assert.Equal(t, err, Return(err).Err())
}

func TestOnErrorOnNoErrorReturnsSpecialError(t *testing.T) {
	e := Return(nil).OnError()
	assert.Equal(t, Return(ErrorWasExpected), e)
}

func TestOnErrorOnErrorReturnsNoError(t *testing.T) {
	err := errors.New("Some Error")
	e := Return(err).OnError()
	assert.Equal(t, Return(nil), e)
}

func TestOnErrorFnOnNoErrorDoesNotExecuteProvidedBlock(t *testing.T) {
	var got error = nil
	executed := false
	Return(nil).OnErrorFn(func(err error) { executed = true; got = err })
	assert.Equal(t, false, executed)
	assert.Equal(t, nil, got)
}

func TestOnErrorFnOnErrorDoesExecuteProvidedBlock(t *testing.T) {
	var got error = nil
	executed := false
	err := errors.New("Error")
	Return(err).OnErrorFn(func(err error) { executed = true; got = err })
	assert.Equal(t, true, executed)
	assert.Equal(t, err, got)
}

func TestChainHelperCallsAllFunctionsWhenNoError(t *testing.T) {
	i := 0
	executed_1 := -1
	executed_2 := -1
	executed_3 := -1

	Chain(
		func() error { executed_1 = i; i++; return nil },
		func() error { executed_2 = i; i++; return nil },
		func() error { executed_3 = i; i++; return nil },
	)

	assert.Equal(t, 0, executed_1)
	assert.Equal(t, 1, executed_2)
	assert.Equal(t, 2, executed_3)
}

func TestChainHelperCallsFunctionsUntilErrorOccurs(t *testing.T) {
	i := 0
	executed_1 := -1
	executed_2 := -1
	executed_3 := -1
	err := errors.New("Very peculiar error")

	Chain(
		func() error { executed_1 = i; i++; return err },
		func() error { executed_2 = i; i++; return nil },
		func() error { executed_3 = i; i++; return nil },
	)

	assert.Equal(t, 0, executed_1)
	assert.Equal(t, -1, executed_2)
	assert.Equal(t, -1, executed_3)
}

func TestChainCallsAllFunctionsWhenNoError(t *testing.T) {
	e := Return(nil)

	i := 0
	executed_1 := -1
	executed_2 := -1
	executed_3 := -1

	e.Chain(
		func() error { executed_1 = i; i++; return nil },
		func() error { executed_2 = i; i++; return nil },
		func() error { executed_3 = i; i++; return nil },
	)

	assert.Equal(t, 0, executed_1)
	assert.Equal(t, 1, executed_2)
	assert.Equal(t, 2, executed_3)
}

func TestChainCallsFunctionsUntilErrorOccurs(t *testing.T) {
	e := Return(nil)

	i := 0
	executed_1 := -1
	executed_2 := -1
	executed_3 := -1
	err := errors.New("Very peculiar error")

	e.Chain(
		func() error { executed_1 = i; i++; return err },
		func() error { executed_2 = i; i++; return nil },
		func() error { executed_3 = i; i++; return nil },
	)

	assert.Equal(t, 0, executed_1)
	assert.Equal(t, -1, executed_2)
	assert.Equal(t, -1, executed_3)
}
