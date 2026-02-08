import { useState, useCallback } from "react";
import { TooltipProvider } from "@/components/ui/tooltip";
import DropPage from "@/pages/DropPage";
import GalleryPage from "@/pages/GalleryPage";
import {
  OpenFile,
  GetTemplates,
} from "../wailsjs/go/main/App";
import type { AppScreen, ParseResult, TemplateInfo } from "@/types";

export default function App() {
  const [screen, setScreen] = useState<AppScreen>("drop");
  const [resumeInfo, setResumeInfo] = useState<ParseResult | null>(null);
  const [templates, setTemplates] = useState<TemplateInfo[]>([]);
  const [error, setError] = useState<string | null>(null);

  const handleOpenNative = useCallback(async () => {
    setError(null);
    try {
      const result = await OpenFile();
      const tmpls = await GetTemplates();
      setResumeInfo(result);
      setTemplates(tmpls);
      setScreen("gallery");
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e);
      if (!msg.includes("no file selected")) {
        setError(msg);
      }
    }
  }, []);

  const handleFile = useCallback(
    (_data: Uint8Array, format: string) => {
      if (format.startsWith("__invalid:")) {
        const ext = format.replace("__invalid:", "");
        setError(`Unsupported format ".${ext}". Use .yml, .yaml, .json, or .toml.`);
        return;
      }
      handleOpenNative();
    },
    [handleOpenNative]
  );

  const handleReset = useCallback(() => {
    setResumeInfo(null);
    setTemplates([]);
    setError(null);
    setScreen("drop");
  }, []);

  return (
    <TooltipProvider>
      {screen === "drop" && (
        <DropPage
          onFile={handleFile}
          onOpenNative={handleOpenNative}
          error={error}
        />
      )}
      {screen === "gallery" && resumeInfo && (
        <GalleryPage
          resumeInfo={resumeInfo}
          templates={templates}
          error={error}
          onError={setError}
          onReset={handleReset}
        />
      )}
    </TooltipProvider>
  );
}
