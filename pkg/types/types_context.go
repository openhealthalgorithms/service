package types

// KeyCtx represents typ to store values in a context.WithValue.
type KeyCtx int

// KeyValuesCtx represents the key to store an instance of ValuesCtx in a context.WithValue.
const KeyValuesCtx KeyCtx = 1

// ValuesCtx represents a value which may be safely used to store/pass a date through the context.
type ValuesCtx struct {
	Params *SafeMap
}

// NewValuesCtx returns a new instance of ValueCtx ready to use.
func NewValuesCtx() ValuesCtx {
	return ValuesCtx{
		Params: NewSafeMap(),
	}
}
