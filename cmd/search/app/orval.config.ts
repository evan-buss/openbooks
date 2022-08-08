import { defineConfig } from "orval";

export default defineConfig({
  search: {
    output: {
      mode: "tags-split",
      target: "src/generated/api.ts",
      schemas: "src/generated/models",
      client: "react-query",
      prettier: true,
      mock: true,
      override: {
        useDates: true,
      },
    },

    input: {
      target: "../docs/swagger.yaml",
    },
  },
});
