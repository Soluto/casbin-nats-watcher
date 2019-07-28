package natswatcher

import (
	"fmt"
	"runtime"

	"github.com/casbin/casbin/persist"
	"github.com/nats-io/go-nats"
)

// Watcher implements persist.Watcher interface
type Watcher struct {
	endpoint             string
	options              []nats.Option
	connection           *nats.Conn
	subscription         *nats.Subscription
	policyUpdatedSubject string
	callback             func(string)
}

// NewWatcher creates new Nats watcher.
// Parameters:
// - endpoint
//		Endpoint of Nats server
// - policyUpdatedSubject
//      Nats subject that sends message when policy was updated externally. It leads to call of callback
// - options
//      Options to connect Nats like user, password, etc.
func NewWatcher(endpoint string, policyUpdatedSubject string, options ...nats.Option) (persist.Watcher, error) {
	nw := &Watcher{
		endpoint:             endpoint,
		options:              options,
		policyUpdatedSubject: policyUpdatedSubject,
	}

	// Connecting to Nats
	err := nw.connect()
	if err != nil {
		return nil, err
	}

	// Subscribe to updates
	sub, err := nw.subcribeToUpdates()
	if err != nil {
		return nil, err
	}
	nw.subscription = sub

	runtime.SetFinalizer(nw, finalizer)

	return nw, nil
}

// SetUpdateCallback sets the callback function that the watcher will call
// when the policy in DB has been changed by other instances.
// A classic callback is Enforcer.LoadPolicy().
func (w *Watcher) SetUpdateCallback(callback func(string)) error {
	w.callback = callback
	return nil
}

// Update calls the update callback of other instances to synchronize their policy.
// It is usually called after changing the policy in DB, like Enforcer.SavePolicy(),
// Enforcer.AddPolicy(), Enforcer.RemovePolicy(), etc.
func (w *Watcher) Update() error {
	if w.connection != nil && w.connection.Status() == nats.CONNECTED {
		w.connection.Publish(w.policyUpdatedSubject, []byte(""))
		return nil
	}
	return fmt.Errorf("Connection is nil or in invalid state")
}

func (w *Watcher) connect() error {
	nc, err := nats.Connect(w.endpoint, w.options...)
	if err != nil {
		return err
	}
	w.connection = nc
	return nil
}

// Close stops and releases the watcher, the callback function will not be called any more.
func (w *Watcher) Close() {
	finalizer(w)
}

func (w *Watcher) subcribeToUpdates() (*nats.Subscription, error) {
	sub, err := w.connection.Subscribe(w.policyUpdatedSubject, func(msg *nats.Msg) {
		if w.callback != nil {
			w.callback(string(msg.Data[:]))
		}
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func finalizer(w *Watcher) {
	if w.subscription != nil && w.subscription.IsValid() {
		w.subscription.Unsubscribe()
	}
	w.subscription = nil

	if w.connection != nil && !w.connection.IsClosed() {
		w.connection.Close()
	}
	w.connection = nil

	w.callback = nil
	return
}
