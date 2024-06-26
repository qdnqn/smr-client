package commands

import (
	"fmt"
	"github.com/qdnqn/smr-client/pkg/commands/resource"
	"github.com/qdnqn/smr-client/pkg/manager"
	"os"
)

const HELP_RESOURCE string = "Eg: smr resource [describe, list]"

func ResourceCommand() {
	Commands = append(Commands, Command{
		name: "resource",
		condition: func(*manager.Manager) bool {
			return true
		},
		functions: []func(*manager.Manager, []string){
			func(mgr *manager.Manager, args []string) {
				if len(os.Args) > 2 {
					switch os.Args[2] {
					case "describe":
						resource.Describe(mgr.Context)
						break
					case "list":
						resource.List(mgr.Context)
						break
					case "get":
						if len(os.Args) > 4 {
							resource.Get(mgr.Context, os.Args[3], os.Args[4])
						} else {
							fmt.Println(HELP_RESOURCE)
						}
						break
					case "edit":
						if len(os.Args) > 4 {
							resource.Edit(mgr.Context, os.Args[3], os.Args[4])
						} else {
							fmt.Println(HELP_RESOURCE)
						}
						break
					default:
						fmt.Println(HELP_RESOURCE)
					}
				} else {
					fmt.Println(HELP_RESOURCE)
				}
			},
		},
		depends_on: []func(*manager.Manager, []string){
			func(mgr *manager.Manager, args []string) {
			},
		},
	})
}
