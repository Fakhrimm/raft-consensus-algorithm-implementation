import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: true, // Bind to all IP addresses
    port: 60005, // Optionally specify a different port
  },
});
