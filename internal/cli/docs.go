package cli

import (
	"context"

	"github.com/mishankov/platforma/docs"
	"github.com/mishankov/platforma/fileserver"
)

func docsCommand(_ []string) {
	server := fileserver.New(docs.Assets(), "/platforma/", "4444")

	server.Run(context.Background())
}
