package commands

import (
	"github.com/qdnqn/smr-client/pkg/manager"
	"os"
)

var Commands []Command

func PreloadCommands() {
	Apply()
	Context()
	Gitops()
	Ps()
}

func Run(mgr *manager.Manager) {
	for _, comm := range Commands {
		for _, arg := range os.Args {
			if comm.name == arg {
				if comm.condition(mgr) {
					for _, fn := range comm.depends_on {
						fn(mgr, os.Args)
					}

					for _, fn := range comm.functions {
						fn(mgr, os.Args)
					}
				}
			}
		}
	}
}
