package cli

import "github.com/mishankov/platforma/log"

func Run(args []string) {
	command := args[1]
	switch command {
	case "generate":
		generateCommand(args[2:])
	default:
		log.Error("unknown command", "command", command)
	}
}
