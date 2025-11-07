package cli

import (
	"context"

	"github.com/mishankov/platforma/docs"
	"github.com/mishankov/platforma/httpserver"
	"github.com/mishankov/platforma/log"
)

func docsCommand(_ []string) {
	ctx := context.Background()
	server := httpserver.NewFileServer(docs.Assets(), "/platforma", "4444")

	if err := server.Run(ctx); err != nil {
		log.ErrorContext(ctx, "documentation serving ended with error", "error", err)
	}
}
