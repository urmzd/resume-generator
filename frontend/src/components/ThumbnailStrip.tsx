import { useCallback, useEffect, useRef, useState } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import ThumbnailCard from "./ThumbnailCard";
import type { TemplateInfo } from "@/types";

interface ThumbnailStripProps {
  templates: TemplateInfo[];
  selectedIndex: number;
  onSelect: (index: number) => void;
  getCachedUrl: (name: string) => string | null;
}

export default function ThumbnailStrip({ templates, selectedIndex, onSelect, getCachedUrl }: ThumbnailStripProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(false);

  const updateScrollButtons = useCallback(() => {
    const el = containerRef.current;
    if (!el) return;
    setCanScrollLeft(el.scrollLeft > 0);
    setCanScrollRight(el.scrollLeft + el.clientWidth < el.scrollWidth - 1);
  }, []);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    updateScrollButtons();
    el.addEventListener("scroll", updateScrollButtons, { passive: true });
    const observer = new ResizeObserver(updateScrollButtons);
    observer.observe(el);
    return () => {
      el.removeEventListener("scroll", updateScrollButtons);
      observer.disconnect();
    };
  }, [updateScrollButtons]);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;
    const active = container.children[selectedIndex] as HTMLElement | undefined;
    active?.scrollIntoView({ behavior: "smooth", inline: "center", block: "nearest" });
  }, [selectedIndex]);

  const scroll = useCallback((direction: "left" | "right") => {
    const el = containerRef.current;
    if (!el) return;
    const amount = el.clientWidth * 0.6;
    el.scrollBy({ left: direction === "left" ? -amount : amount, behavior: "smooth" });
  }, []);

  return (
    <div className="relative">
      {canScrollLeft && (
        <Button
          size="icon"
          variant="secondary"
          onClick={() => scroll("left")}
          className="absolute left-1 top-1/2 -translate-y-1/2 z-10 bg-background/80 backdrop-blur-sm rounded-full shadow-md h-8 w-8"
          aria-label="Scroll thumbnails left"
        >
          <ChevronLeft className="h-4 w-4" />
        </Button>
      )}
      {canScrollRight && (
        <Button
          size="icon"
          variant="secondary"
          onClick={() => scroll("right")}
          className="absolute right-1 top-1/2 -translate-y-1/2 z-10 bg-background/80 backdrop-blur-sm rounded-full shadow-md h-8 w-8"
          aria-label="Scroll thumbnails right"
        >
          <ChevronRight className="h-4 w-4" />
        </Button>
      )}
      <div
        ref={containerRef}
        role="tablist"
        className="flex justify-center gap-2 overflow-x-auto scrollbar-thin px-4 py-1.5"
      >
        {templates.map((t, i) => (
          <ThumbnailCard
            key={t.name}
            template={t}
            isActive={i === selectedIndex}
            previewUrl={getCachedUrl(t.name)}
            onClick={() => onSelect(i)}
          />
        ))}
      </div>
    </div>
  );
}
