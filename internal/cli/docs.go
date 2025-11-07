package cli

import (
	"context"

	"github.com/platforma-dev/platforma/docs"
	"github.com/platforma-dev/platforma/httpserver"
	"github.com/platforma-dev/platforma/log"
)

func docsCommand(_ []string) {
	ctx := context.Background()
	server := httpserver.NewFileServer(docs.Assets(), "/platforma", "4444")

	if err := server.Run(ctx); err != nil {
		log.ErrorContext(ctx, "documentation serving ended with error", "error", err)
	}
}
