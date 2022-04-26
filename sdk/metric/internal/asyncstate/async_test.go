package asyncstate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/data"
	"go.opentelemetry.io/otel/sdk/metric/internal/pipeline"
	"go.opentelemetry.io/otel/sdk/metric/internal/test"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

var (
	testLibrary = instrumentation.Library{
		Name: "test",
	}

	endTime    = time.Unix(100, 0)
	middleTime = endTime.Add(-time.Millisecond)
	startTime  = endTime.Add(-2 * time.Millisecond)

	testSequence = data.Sequence{
		Start: startTime,
		Last:  middleTime,
		Now:   endTime,
	}
)

type testSDK struct {
	compiler *viewstate.Compiler
}

func (tsdk *testSDK) compile(desc sdkinstrument.Descriptor) pipeline.Register[viewstate.Instrument] {
	comp, err := tsdk.compiler.Compile(desc)
	if err != nil {
		panic(err)
	}
	reg := pipeline.NewRegister[viewstate.Instrument](1)
	reg[0] = comp
	return reg
}

func testAsync(name string, opts ...view.Option) *testSDK {
	return &testSDK{
		compiler: viewstate.New(testLibrary, view.New(name, opts...)),
	}
}

func testState() *State {
	return NewState(0)
}

func testObserver[N number.Any, Traits number.Traits[N]](tsdk *testSDK, name string, ik sdkinstrument.Kind, opts ...instrument.Option) observer[N, Traits] {
	var t Traits
	desc := test.Descriptor(name, ik, t.Kind(), opts...)
	impl := NewInstrument(desc, tsdk, tsdk.compile(desc))
	return NewObserver[N, Traits](impl)
}

func TestNewCallbackError(t *testing.T) {
	tsdk := testAsync("test")

	// no instruments error
	cb, err := NewCallback(nil, tsdk, nil)
	require.Error(t, err)
	require.Nil(t, cb)

	// nil callback error
	cntr := testObserver[int64, number.Int64Traits](tsdk, "counter", sdkinstrument.CounterObserverKind)
	cb, err = NewCallback([]instrument.Asynchronous{cntr}, tsdk, nil)
	require.Error(t, err)
	require.Nil(t, cb)
}

func TestNewCallbackProviderMismatch(t *testing.T) {
	test0 := testAsync("test0")
	test1 := testAsync("test1")

	instA0 := testObserver[int64, number.Int64Traits](test0, "A", sdkinstrument.CounterObserverKind)
	instB1 := testObserver[float64, number.Float64Traits](test1, "A", sdkinstrument.CounterObserverKind)

	cb, err := NewCallback([]instrument.Asynchronous{instA0, instB1}, test0, func(context.Context) {})
	require.Error(t, err)
	require.Contains(t, err.Error(), "asynchronous instrument belongs to a different meter")
	require.Nil(t, cb)

	cb, err = NewCallback([]instrument.Asynchronous{instA0, instB1}, test1, func(context.Context) {})
	require.Error(t, err)
	require.Contains(t, err.Error(), "asynchronous instrument belongs to a different meter")
	require.Nil(t, cb)

	cb, err = NewCallback([]instrument.Asynchronous{instA0}, test0, func(context.Context) {})
	require.NoError(t, err)
	require.NotNil(t, cb)

	cb, err = NewCallback([]instrument.Asynchronous{instB1}, test1, func(context.Context) {})
	require.NoError(t, err)
	require.NotNil(t, cb)

	// nil value not of this SDK
	var fake0 instrument.Asynchronous
	cb, err = NewCallback([]instrument.Asynchronous{fake0}, test0, func(context.Context) {})
	require.Error(t, err)
	require.Contains(t, err.Error(), "asynchronous instrument does not belong to this SDK")
	require.Nil(t, cb)

	// non-nil value not of this SDK
	var fake1 instrument.AsynchronousStruct
	cb, err = NewCallback([]instrument.Asynchronous{fake1}, test0, func(context.Context) {})
	require.Error(t, err)
	require.Contains(t, err.Error(), "asynchronous instrument does not belong to this SDK")
	require.Nil(t, cb)
}

func TestCallbackInvalidation(t *testing.T) {
	errors := test.OTelErrors()

	tsdk := testAsync("test")

	var called int64
	var saveCtx context.Context

	cntr := testObserver[int64, number.Int64Traits](tsdk, "counter", sdkinstrument.CounterObserverKind)
	cb, err := NewCallback([]instrument.Asynchronous{cntr}, tsdk, func(ctx context.Context) {
		cntr.Observe(ctx, called)
		saveCtx = ctx
		called++
	})
	require.NoError(t, err)

	state := testState()

	// run the callback once legitimately
	cb.Run(context.Background(), state)

	// simulate use after callback return
	cntr.Observe(saveCtx, 10000000)

	cntr.inst.SnapshotAndProcess(state)

	require.Equal(t, int64(1), called)
	require.Equal(t, 1, len(*errors))
	require.Contains(t, (*errors)[0].Error(), "used after callback return")

	test.RequireEqualMetrics(
		t,
		test.CollectScope(
			t,
			tsdk.compiler.Collectors(),
			testSequence,
		),
		test.Instrument(
			cntr.inst.descriptor,
			test.Point(startTime, endTime, sum.NewMonotonicInt64(0)),
		),
	)
}

func TestCallbackUndeclaredInstrument(t *testing.T) {
	errors := test.OTelErrors()

	tt := testAsync("test")

	var called int64

	cntr1 := testObserver[int64, number.Int64Traits](tt, "counter1", sdkinstrument.CounterObserverKind)
	cntr2 := testObserver[int64, number.Int64Traits](tt, "counter2", sdkinstrument.CounterObserverKind)

	cb, err := NewCallback([]instrument.Asynchronous{cntr1}, tt, func(ctx context.Context) {
		cntr2.Observe(ctx, called)
		called++
	})
	require.NoError(t, err)

	state := testState()

	// run the callback once legitimately
	cb.Run(context.Background(), state)

	cntr1.inst.SnapshotAndProcess(state)
	cntr2.inst.SnapshotAndProcess(state)

	require.Equal(t, int64(1), called)
	require.Equal(t, 1, len(*errors))
	require.Contains(t, (*errors)[0].Error(), "instrument not declared for use in callback")

	test.RequireEqualMetrics(
		t,
		test.CollectScope(
			t,
			tt.compiler.Collectors(),
			testSequence,
		),
		test.Instrument(
			cntr1.inst.descriptor,
		),
		test.Instrument(
			cntr2.inst.descriptor,
		),
	)
}

func TestCallbackDroppedInstrument(t *testing.T) {
	errors := test.OTelErrors()

	tt := testAsync("test",
		view.WithClause(
			view.MatchInstrumentName("drop"),
			view.WithAggregation(aggregation.DropKind),
		),
	)

	cntrDrop := testObserver[float64, number.Float64Traits](tt, "drop", sdkinstrument.CounterObserverKind)
	cntrKeep := testObserver[float64, number.Float64Traits](tt, "keep", sdkinstrument.CounterObserverKind)

	cb, _ := NewCallback([]instrument.Asynchronous{cntrKeep}, tt, func(ctx context.Context) {
		cntrDrop.Observe(ctx, 1000)
		cntrKeep.Observe(ctx, 1000)
	})

	state := testState()

	cb.Run(context.Background(), state)

	cntrKeep.inst.SnapshotAndProcess(state)
	cntrDrop.inst.SnapshotAndProcess(state)

	require.Equal(t, 1, len(*errors))
	require.Contains(t, (*errors)[0].Error(), "instrument not declared for use in callback")

	test.RequireEqualMetrics(
		t,
		test.CollectScope(
			t,
			tt.compiler.Collectors(),
			testSequence,
		),
		test.Instrument(
			cntrKeep.inst.descriptor,
			test.Point(startTime, endTime, sum.NewMonotonicFloat64(1000)),
		),
	)
}

func TestInstrumentUseOutsideCallback(t *testing.T) {
	errors := test.OTelErrors()

	tt := testAsync("test")

	cntr := testObserver[float64, number.Float64Traits](tt, "cntr", sdkinstrument.CounterObserverKind)

	cntr.Observe(context.Background(), 1000)

	state := testState()

	cntr.inst.SnapshotAndProcess(state)

	require.Equal(t, 1, len(*errors))
	require.Contains(t, (*errors)[0].Error(), "async instrument used outside of callback")

	test.RequireEqualMetrics(
		t,
		test.CollectScope(
			t,
			tt.compiler.Collectors(),
			testSequence,
		),
		test.Instrument(
			cntr.inst.descriptor,
		),
	)
}
