package natswatcher

import (
	"runtime"

	"github.com/casbin/casbin/persist"
	"github.com/nats-io/go-nats"
)

// Watcher implements persist.Watcher interface
type Watcher struct {
	endpoint                     string
	connection                   *nats.Conn
	subscription                 *nats.Subscription
	policyUpdatedSubject         string
	policyUpdatedByCasbinSubject string
	callback                     func(string)
}

// NewWatcher creates new Nats watcher.
// Parameters:
// - endpoint
//		Endpoint of Nats server
// - policyUpdatedSubject
//      Nats subject that sends message when policy was updated externally. It leads to call of callback
// - policyUpdatedByCasbinSubject
//      Nats subject that sends message when policy was updated by casbin. (like SavePolicy, AddPolicy, RemovePolicy)
func NewWatcher(endpoint string, policyUpdatedSubject string, policyUpdatedByCasbinSubject string) (persist.Watcher, error) {
	nw := &Watcher{
		endpoint:                     endpoint,
		policyUpdatedSubject:         policyUpdatedSubject,
		policyUpdatedByCasbinSubject: policyUpdatedByCasbinSubject,
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

	runtime.SetFinalizer(nw, nw.unsubscribe)

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
	w.connection.Publish(w.policyUpdatedByCasbinSubject, []byte(""))
	return nil
}

func (w *Watcher) connect() error {
	nc, err := nats.Connect(w.endpoint)
	if err != nil {
		return err
	}
	w.connection = nc
	return nil
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

func (w *Watcher) unsubscribe() error {
	return w.subscription.Unsubscribe()
}
