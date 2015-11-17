package args

import (
	"fmt"
	"os"
)

func Parse() (string, []string) {
	if len(os.Args) == 1 {
		PrintHelp()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		return os.Args[1], []string{}
	} else {
		return os.Args[1], os.Args[2:]
	}
}

func PrintHelp() {
	fmt.Fprintf(os.Stdout, `Usage: resty <command> <args>

  The following commands are supported:

  help 	   -  This command will print the current help message
  generate - This command will generate one of the following - model, controller, scaffold.
  
    model - This will generate a model file with the specified fields
      e.g. resty generate model User name:string age:int32

    controller - This will generate a controller file with the specified actions. Actions other than index, show, update, create and destroy, will be stubbed.
      e.g. resty generate controller User index show update create

`)
}
