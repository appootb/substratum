package configure

var (
	callbackImpl Callback
)

// CallbackImplementor returns the value updated callback service implementor.
func CallbackImplementor() Callback {
	return callbackImpl
}

// RegisterCallbackImplementor registers the value updated callback service implementor.
func RegisterCallbackImplementor(callback Callback) {
	callbackImpl = callback
}

type UpdateEvent func()

type CallbackFunc struct {
	Value DynamicType
	Event UpdateEvent
}

type Callback interface {
	RegChan() chan<- *CallbackFunc
	EvtChan() chan<- DynamicType
}
