// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

// https://astro.build/config
export default defineConfig({
  integrations: [
    starlight({
      title: "platforma",
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/mishankov/platforma",
        },
      ],
      sidebar: [
        {
          slug: "getting-started",
        },
        {
          label: "Packages",
          items: ["packages/database", "packages/httpserver"],
        },
      ],
    }),
  ],
});
