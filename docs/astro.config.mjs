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
          items: ["packages/database", "packages/httpserver", "packages/log"],
        },
        {
          slug: "cli",
        },
      ],
    }),
  ],
});
