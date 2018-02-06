# casbin-nats-watcher

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/Soluto/casbin-nats-watcher.svg?branch=master)](https://travis-ci.org/Soluto/casbin-nats-watcher)
[![Coverage Status](https://coveralls.io/repos/github/Soluto/casbin-nats-watcher/badge.svg?branch=master)](https://coveralls.io/github/Soluto/casbin-nats-watcher?branch=master)
[![Godoc](https://godoc.org/github.com/Soluto/casbin-nats-watcher?status.svg)](https://godoc.org/github.com/Soluto/casbin-nats-watcher)

[Casbin](https://github.com/casbin/casbin) watcher implementation with Nats.io

## Installation

    go get github.com/Soluto/casbin-nats-watcher

## Usage
```go
import (
	natswatcher "github.com/Soluto/casbin-nats-watcher"
	"github.com/casbin/casbin"
)

func main() {
	watcher, _ := natswatcher.NewWatcher("http://...", "my-policy-subject")

	enforcerer := casbin.NewSyncedEnforcer("model.conf", "policy.csv")
	enforcerer.SetWatcher(watcher)
}
```

## Connected pojects
- [Casbin](https://github.com/casbin/casbin)
- [Nats.io](https://github.com/nats-io/go-nats)


## Additional Usage Examples

For real-world example visit [Tweek](https://github.com/Soluto/tweek).

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
