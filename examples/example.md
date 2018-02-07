# CASBIN-NATS-WATCHER Example

```go
import (
    natswatcher "github.com/Soluto/casbin-nats-watcher"
    "github.com/casbin/casbin"
)

func main() {
    watcher, _ := natswatcher.NewWatcher("http://nats-endpoint", "my-policy-subject")

    enforcer := casbin.NewSyncedEnforcer("rbac_model.conf", "rbac_policy.csv")
    enforcer.SetWatcher(watcher)
}
```
