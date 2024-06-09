import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "0.0.0.0", // Bind to all IP addresses
    port: 60005, // Optionally specify a different port
  },
});
