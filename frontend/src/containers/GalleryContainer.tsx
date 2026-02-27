import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";
import { PdfCache } from "@/lib/pdf-cache";
import { GeneratePDF, SavePDF, SaveNative } from "../../wailsjs/go/main/App";
import type { TemplateInfo } from "@/types";

interface GalleryRenderProps {
  selectedIndex: number;
  selectedTemplate: TemplateInfo;
  pdfUrl: string | null;
  isLoading: boolean;
  cacheVersion: number;
  pageCount: number | null;
  onSelectTemplate: (index: number) => void;
  onPrev: () => void;
  onNext: () => void;
  onSavePdf: () => void;
  onSaveNative: () => void;
  getCachedUrl: (name: string) => string | null;
  invalidateCache: () => void;
}

interface GalleryContainerProps {
  templates: TemplateInfo[];
  onError: (msg: string) => void;
  children: (props: GalleryRenderProps) => ReactNode;
}

export default function GalleryContainer({ templates, onError, children }: GalleryContainerProps) {
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [pdfUrl, setPdfUrl] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [cacheVersion, setCacheVersion] = useState(0);
  const [pageCount, setPageCount] = useState<number | null>(null);
  const cacheRef = useRef(new PdfCache());
  const generationIdRef = useRef(0);
  const bgGenerationStarted = useRef(false);

  const generateForTemplate = useCallback(
    async (index: number) => {
      const template = templates[index];
      if (!template) return;

      const cached = cacheRef.current.get(template.name);
      if (cached) {
        setPdfUrl(cached);
        return;
      }

      setIsLoading(true);
      const genId = ++generationIdRef.current;

      try {
        const result = await GeneratePDF(template.name);
        // Guard against stale responses from fast navigation
        if (genId !== generationIdRef.current) return;
        const url = cacheRef.current.set(template.name, result.data);
        setPdfUrl(url);
        setPageCount(result.pageCount);
      } catch (e) {
        if (genId !== generationIdRef.current) return;
        onError(e instanceof Error ? e.message : String(e));
      } finally {
        if (genId === generationIdRef.current) {
          setIsLoading(false);
        }
      }
    },
    [templates, onError]
  );

  const generateInBackground = useCallback(
    async (index: number) => {
      const template = templates[index];
      if (!template) return;
      if (cacheRef.current.get(template.name)) return;

      try {
        const result = await GeneratePDF(template.name);
        cacheRef.current.set(template.name, result.data);
        setCacheVersion((v) => v + 1);
      } catch {
        // Silently skip failed background generation
      }
    },
    [templates]
  );

  const onSelectTemplate = useCallback(
    (index: number) => {
      const clamped = Math.max(0, Math.min(index, templates.length - 1));
      setSelectedIndex(clamped);
      generateForTemplate(clamped);
    },
    [templates.length, generateForTemplate]
  );

  const onPrev = useCallback(() => {
    if (selectedIndex > 0) onSelectTemplate(selectedIndex - 1);
  }, [selectedIndex, onSelectTemplate]);

  const onNext = useCallback(() => {
    if (selectedIndex < templates.length - 1) onSelectTemplate(selectedIndex + 1);
  }, [selectedIndex, templates.length, onSelectTemplate]);

  // Auto-generate first template on mount
  useEffect(() => {
    if (templates.length > 0) {
      generateForTemplate(0);
    }
  }, [templates, generateForTemplate]);

  // Background generation for all templates after first render
  useEffect(() => {
    if (templates.length === 0 || bgGenerationStarted.current) return;
    bgGenerationStarted.current = true;

    let cancelled = false;
    (async () => {
      // Small delay to let the selected template finish first
      await new Promise((r) => setTimeout(r, 500));
      for (let i = 0; i < templates.length; i++) {
        if (cancelled) break;
        if (cacheRef.current.get(templates[i].name)) continue;
        await generateInBackground(i);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [templates, generateInBackground]);

  // Keyboard navigation
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const tag = (e.target as HTMLElement)?.tagName;
      if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;

      if (e.key === "ArrowLeft") {
        e.preventDefault();
        onPrev();
      } else if (e.key === "ArrowRight") {
        e.preventDefault();
        onNext();
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [onPrev, onNext]);

  // Cleanup blob URLs on unmount
  useEffect(() => {
    const cache = cacheRef.current;
    return () => cache.clear();
  }, []);

  const onSavePdf = useCallback(async () => {
    const template = templates[selectedIndex];
    if (!template) return;
    try {
      await SavePDF(template.name);
    } catch (e) {
      onError(e instanceof Error ? e.message : String(e));
    }
  }, [templates, selectedIndex, onError]);

  const onSaveNative = useCallback(async () => {
    const template = templates[selectedIndex];
    if (!template) return;
    try {
      await SaveNative(template.name);
    } catch (e) {
      onError(e instanceof Error ? e.message : String(e));
    }
  }, [templates, selectedIndex, onError]);

  const selectedTemplate = templates[selectedIndex];
  if (!selectedTemplate) return null;

  const getCachedUrl = useCallback(
    (name: string) => cacheRef.current.get(name) ?? null,
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [cacheVersion]
  );

  const invalidateCache = useCallback(() => {
    cacheRef.current.clear();
    bgGenerationStarted.current = false;
    setPdfUrl(null);
    setCacheVersion((v) => v + 1);
    generateForTemplate(selectedIndex);
  }, [selectedIndex, generateForTemplate]);

  return (
    <>
      {children({
        selectedIndex,
        selectedTemplate,
        pdfUrl,
        isLoading,
        cacheVersion,
        pageCount,
        onSelectTemplate,
        onPrev,
        onNext,
        onSavePdf,
        onSaveNative,
        getCachedUrl,
        invalidateCache,
      })}
    </>
  );
}
