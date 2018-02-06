package natswatcher

import (
	"fmt"
	"testing"
	"time"

	"github.com/AleF83/casbin"
	gnatsd "github.com/nats-io/gnatsd/test"
	"github.com/nats-io/go-nats"
)

func TestWatcher(t *testing.T) {
	// Setup nats server
	s := gnatsd.RunDefaultServer()
	defer s.Shutdown()

	natsEndpoint := fmt.Sprintf("nats://localhost:%d", nats.DefaultPort)
	natsSubject := "casbin-policy-updated-subject"

	updaterCh := make(chan string, 1)
	listenerCh := make(chan string, 1)

	// updater represents the Casbin enforcer instance that changes the policy in DB.
	// Use the endpoint of nats as parameter.
	updater, err := NewWatcher(natsEndpoint, natsSubject)
	if err != nil {
		t.Error("Failed to create updater")
	}
	updater.SetUpdateCallback(func(msg string) {
		go func() {
			updaterCh <- "updater"
			close(updaterCh)
		}()
	})

	// listener represents any other Casbin enforcer instance that watches the change of policy in DB.
	listener, err := NewWatcher(natsEndpoint, natsSubject)
	if err != nil {
		t.Error("Failed to create listener")
	}

	// listener should set a callback that gets called when policy changes.
	err = listener.SetUpdateCallback(func(msg string) {
		go func() {
			listenerCh <- "listener"
			close(listenerCh)
		}()
	})
	if err != nil {
		t.Error("Failed to set listener callback")
	}

	// updater changes the policy, and sends the notifications.
	err = updater.Update()
	if err != nil {
		t.Error("The updater failed to send Update")
	}

	// Validate that listener received message
	select {
	case res := <-listenerCh:
		if res != "listener" {
			t.Errorf("Message from unknown source: %v", res)
		}
	case res := <-updaterCh:
		if res != "updater" {
			t.Errorf("Message from unknown source: %v", res)
		}
	case <-time.After(time.Second * 10):
		close(updaterCh)
		close(listenerCh)
		t.Error("Updater or listener didn't received message in time")
	}

}

func TestWithEnforcer(t *testing.T) {
	// Setup nats server
	s := gnatsd.RunDefaultServer()
	defer s.Shutdown()

	natsEndpoint := fmt.Sprintf("nats://localhost:%d", nats.DefaultPort)
	natsSubject := "casbin-policy-updated-subject"
	cannel := make(chan string, 1)

	// Initialize the watcher.
	// Use the endpoint of etcd cluster as parameter.
	w, err := NewWatcher(natsEndpoint, natsSubject)
	if err != nil {
		t.Error("Failed to create updater")
	}

	// Initialize the enforcer.
	e := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	// Set the watcher for the enforcer.
	e.SetWatcher(w)

	// By default, the watcher's callback is automatically set to the
	// enforcer's LoadPolicy() in the SetWatcher() call.
	// We can change it by explicitly setting a callback.
	w.SetUpdateCallback(func(msg string) {
		cannel <- "enforcerer"
	})

	// Update the policy to test the effect.
	e.SavePolicy()

	// Validate that listener received message
	select {
	case res := <-cannel:
		if res != "enforcerer" {
			t.Error("Got unexpected message")
		}
	case <-time.After(time.Second * 10):
		t.Error("The enforcerer didn't send message in time")
	}
}
