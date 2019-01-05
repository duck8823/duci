package executor_test

import (
	"context"
	"github.com/duck8823/duci/application/service/executor"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"reflect"
	"testing"
)

func TestDefaultExecutorBuilder(t *testing.T) {
	// given
	want := &executor.Builder{}
	defer want.SetStartFunc(executor.NothingToDoStart)()
	defer want.SetLogFunc(runner.NothingToDo)()
	defer want.SetEndFunc(executor.NothingToDoEnd)()

	// when
	got, err := executor.DefaultExecutorBuilder()

	// then
	if err != nil {
		t.Errorf("error must be nil, but got %+v", err)
	}

	// and
	opts := cmp.Options{
		cmp.AllowUnexported(executor.Builder{}),
		cmp.Transformer("startFuncToPointer", func(f func(context.Context)) uintptr {
			return reflect.ValueOf(f).Pointer()
		}),
		cmp.Transformer("logFuncToPointer", func(f runner.LogFunc) uintptr {
			return reflect.ValueOf(f).Pointer()
		}),
		cmp.Transformer("endFuncToPointer", func(f func(context.Context, error)) uintptr {
			return reflect.ValueOf(f).Pointer()
		}),
		cmpopts.IgnoreInterfaces(struct{ docker.Docker }{}),
	}
	if !cmp.Equal(got, want, opts) {
		t.Errorf("must be equal. but: %+v", cmp.Diff(got, want, opts))
	}
}

func TestBuilder_StartFunc(t *testing.T) {
	// given
	startFunc := func(context.Context) {}

	// and
	want := &executor.Builder{}
	defer want.SetStartFunc(startFunc)()

	// and
	sut := &executor.Builder{}

	// when
	got := sut.StartFunc(startFunc)

	// then
	opts := cmp.Options{
		cmp.AllowUnexported(executor.Builder{}),
		cmp.Transformer("startFuncToPointer", func(f func(context.Context)) uintptr {
			return reflect.ValueOf(f).Pointer()
		}),
		cmpopts.IgnoreInterfaces(struct{ docker.Docker }{}),
	}
	if !cmp.Equal(got, want, opts) {
		t.Errorf("must be equal. but: %+v", cmp.Diff(got, want, opts))
	}
}

func TestBuilder_LogFunc(t *testing.T) {
	// given
	logFunc := func(context.Context, job.Log) {}

	// and
	want := &executor.Builder{}
	defer want.SetLogFunc(logFunc)()

	// and
	sut := &executor.Builder{}

	// when
	got := sut.LogFunc(logFunc)

	// then
	opts := cmp.Options{
		cmp.AllowUnexported(executor.Builder{}),
		cmp.Transformer("logFuncToPointer", func(f runner.LogFunc) uintptr {
			return reflect.ValueOf(f).Pointer()
		}),
		cmpopts.IgnoreInterfaces(struct{ docker.Docker }{}),
	}
	if !cmp.Equal(got, want, opts) {
		t.Errorf("must be equal. but: %+v", cmp.Diff(got, want, opts))
	}
}

func TestBuilder_EndFunc(t *testing.T) {
	// given
	endFunc := func(context.Context, error) {}

	// and
	want := &executor.Builder{}
	defer want.SetEndFunc(endFunc)()

	// and
	sut := &executor.Builder{}

	// when
	got := sut.EndFunc(endFunc)

	// then
	opts := cmp.Options{
		cmp.AllowUnexported(executor.Builder{}),
		cmp.Transformer("endFuncToPointer", func(f func(context.Context, error)) uintptr {
			return reflect.ValueOf(f).Pointer()
		}),
		cmpopts.IgnoreInterfaces(struct{ docker.Docker }{}),
	}
	if !cmp.Equal(got, want, opts) {
		t.Errorf("must be equal. but: %+v", cmp.Diff(got, want, opts))
	}
}
