package commands

import (
	"fmt"
	"github.com/qdnqn/smr-client/pkg/commands/certkey"
	"github.com/qdnqn/smr-client/pkg/manager"
	"os"
)

const HELP_CERTKEY string = "Eg: smr configuration [describe, list]"

func CertKeyCommand() {
	Commands = append(Commands, Command{
		name: "certkey",
		condition: func(*manager.Manager) bool {
			return true
		},
		functions: []func(*manager.Manager, []string){
			func(mgr *manager.Manager, args []string) {
				if len(os.Args) > 2 {
					switch os.Args[2] {
					case "describe":
						certkey.Describe(mgr.Context)
						break
					case "list":
						certkey.List(mgr.Context)
						break
					case "get":
						if len(os.Args) > 4 {
							certkey.Get(mgr.Context, os.Args[3], os.Args[4])
						} else {
							fmt.Println(HELP_CERTKEY)
						}
						break
					case "edit":
						if len(os.Args) > 4 {
							certkey.Edit(mgr.Context, os.Args[3], os.Args[4])
						} else {
							fmt.Println(HELP_CERTKEY)
						}
						break
					default:
						fmt.Println(HELP_CERTKEY)
					}
				} else {
					fmt.Println(HELP_CERTKEY)
				}
			},
		},
		depends_on: []func(*manager.Manager, []string){
			func(mgr *manager.Manager, args []string) {
			},
		},
	})
}
