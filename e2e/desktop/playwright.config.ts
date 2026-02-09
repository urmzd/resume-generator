import { defineConfig } from "@playwright/test";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const assetsDir = path.resolve(__dirname, "../../assets");

export default defineConfig({
  testDir: ".",
  testMatch: "demo.spec.ts",
  timeout: 120_000,
  use: {
    baseURL: "http://localhost:34115",
    headless: false,
    viewport: { width: 1280, height: 800 },
    video: {
      mode: "on",
      size: { width: 1280, height: 800 },
    },
  },
  outputDir: path.join(assetsDir, "playwright-results"),
  projects: [
    {
      name: "desktop-demo",
      use: { browserName: "chromium" },
    },
  ],
});
