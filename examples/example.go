package example

import (
	natswatcher "github.com/Soluto/casbin-nats-watcher"
	"github.com/casbin/casbin"
)

func main() {
	watcher, _ := natswatcher.NewWatcher("http://...", "my-policy-subject")

	enforcerer := casbin.NewSyncedEnforcer("rbac_model.conf", "rbac_policy.csv")
	enforcerer.SetWatcher(watcher)
}
