import { test, expect } from "@playwright/test";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const RESUME_PATH = path.resolve(
  __dirname,
  "../../assets/example_resumes/software_engineer.yml"
);
const OUTPUT_PDF = path.resolve(__dirname, "../../assets/demo-output.pdf");

// Slow down interactions so the recorded video is easy to follow.
const PAUSE_MS = 1500;

test("Desktop app demo — load, browse templates, save", async ({ page }) => {
  // 1. Navigate to the Wails dev server
  await page.goto("/");
  await expect(page.locator("body")).toBeVisible();

  // Pause on the drop screen so viewers see the initial UI
  await page.waitForTimeout(PAUSE_MS * 2);

  // 2. Load the example resume via the test-only Go binding (bypasses native dialog)
  await page.evaluate(async (filePath: string) => {
    await (window as any)["go"]["main"]["App"]["LoadFileFromPath"](filePath);
  }, RESUME_PATH);

  // The frontend reacts to the Go call, but since we bypassed the React
  // handler we also need to call GetTemplates and push state ourselves.
  // Instead, trigger the native-open flow which calls OpenFile internally —
  // but since a resume is already loaded, the gallery should render.
  // Reload to let the app pick up the loaded resume via a fresh OpenFile call.

  // Actually, the cleanest approach: call LoadFileFromPath + GetTemplates and
  // trigger React state via a custom event or just click the "Open" button
  // after loading. Since the drop page has an "Open File" button, let's use it
  // after pre-loading the resume.

  // Re-approach: LoadFileFromPath sets a.resume on the Go side. Then when
  // OpenFile is called by clicking the button, it'll show a dialog — but we
  // don't want that. Let's instead evaluate both Go calls and programmatically
  // transition the page.

  const templates = await page.evaluate(async () => {
    return await (window as any)["go"]["main"]["App"]["GetTemplates"]();
  });

  // Dispatch a custom event so the React app transitions to the gallery.
  // Since the React app doesn't listen for custom events, we'll navigate
  // the state via the window — but the simplest approach is to reload and
  // use the frontend's own flow.

  // Best approach: directly set React state by simulating what handleOpenNative does.
  // We'll evaluate JS that calls the bound functions and manipulates the DOM.
  await page.evaluate(
    async ({ filePath }: { filePath: string }) => {
      // These are the Wails-bound Go functions
      const app = (window as any)["go"]["main"]["App"];
      const result = await app["LoadFileFromPath"](filePath);
      const tmpls = await app["GetTemplates"]();

      // Store results in a global so we can verify
      (window as any).__demoResult = result;
      (window as any).__demoTemplates = tmpls;
    },
    { filePath: RESUME_PATH }
  );

  // The Go side now has the resume loaded. Trigger the React UI to transition
  // by clicking the "Open File" button — but OpenFile will show a native dialog.
  // Instead, let's use a page reload trick: after LoadFileFromPath, the resume
  // is stored on the Go App struct. We can patch the frontend to skip the dialog.

  // Simplest reliable approach: inject a script that overrides the OpenFile binding
  // to return the already-loaded result, then click the Open button.
  await page.evaluate(() => {
    const app = (window as any)["go"]["main"]["App"];
    const originalOpenFile = app["OpenFile"];
    // Override OpenFile to return cached result (no native dialog)
    app["OpenFile"] = async () => {
      return (window as any).__demoResult;
    };
    // Restore after one call
    const restore = () => {
      app["OpenFile"] = originalOpenFile;
    };
    setTimeout(restore, 5000);
  });

  // Click the "Open File" button on the drop page
  const openButton = page.getByRole("button", { name: /open/i });
  if (await openButton.isVisible()) {
    await openButton.click();
  }

  // Wait for gallery to appear (thumbnail strip with tab roles)
  await expect(page.locator('[role="tab"]').first()).toBeVisible({
    timeout: 30_000,
  });
  await page.waitForTimeout(PAUSE_MS);

  // 3. Click through template thumbnails
  const tabs = page.locator('[role="tab"]');
  const tabCount = await tabs.count();
  const indicesToVisit = [];

  // Visit up to 3 templates (skip the first since it's already selected)
  for (let i = 1; i < Math.min(tabCount, 4); i++) {
    indicesToVisit.push(i);
  }

  for (const idx of indicesToVisit) {
    await tabs.nth(idx).click();
    // Wait for the PDF preview to update
    await page.waitForTimeout(PAUSE_MS * 2);
  }

  // Go back to first template
  if (tabCount > 1) {
    await tabs.first().click();
    await page.waitForTimeout(PAUSE_MS);
  }

  // 4. Save PDF via the test-only binding (bypasses native Save dialog)
  await page.evaluate(
    async ({ templateName, outputPath }: { templateName: string; outputPath: string }) => {
      await (window as any)["go"]["main"]["App"]["SavePDFToPath"](
        templateName,
        outputPath
      );
    },
    {
      templateName: templates[0]?.name ?? "modern-html",
      outputPath: OUTPUT_PDF,
    }
  );

  // Brief pause so viewers see the final state
  await page.waitForTimeout(PAUSE_MS);
});
