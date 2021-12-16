package metric

// callbacksLock sync.Mutex
// callbacks     []*callback

// callback struct {
// 	function    func(context.Context) error
// 	instruments []sdkapi.Instrument
// }

// func (m *Accumulator) NewCallback(instruments []sdkapi.Instrument, function func(context.Context) error) (sdkapi.Callback, error) {
// 	cb := &callback{
// 		function:    function,
// 		instruments: instruments,
// 	}

// 	m.callbacksLock.Lock()
// 	defer m.callbacksLock.Unlock()
// 	m.callbacks = append(m.callbacks, cb)
// 	return cb, nil
// }

// func (cb *callback) Instruments() []sdkapi.Instrument {
// 	return cb.instruments
// }
// m.runCallbacks(ctx)

// func (m *Accumulator) runCallbacks(ctx context.Context) {
// 	m.callbacksLock.Lock()
// 	callbacks := m.callbacks
// 	m.callbacksLock.Unlock()

// 	for _, cb := range callbacks {
// 		cb.function(ctx)
// 	}
// }
