package commands

import "github.com/qdnqn/smr-client/pkg/manager"

type Command struct {
	name       string
	flag       string
	condition  func(*manager.Manager) bool
	functions  []func(*manager.Manager, []string)
	depends_on []func(*manager.Manager, []string)
}
