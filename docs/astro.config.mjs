// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

// https://astro.build/config
export default defineConfig({
  site: "https://mishankov.github.io",
  base: "/platforma",
  integrations: [
    starlight({
      title: "platforma",
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/platforma-dev/platforma",
        },
      ],
      sidebar: [
        {
          slug: "getting-started",
        },
        {
          label: "Packages",
          items: [
            "packages/application",
            "packages/database",
            "packages/httpserver",
            "packages/log",
            "packages/scheduler",
            "packages/queue",
          ],
        },
        {
          slug: "cli",
        },
      ],
    }),
  ],
});
