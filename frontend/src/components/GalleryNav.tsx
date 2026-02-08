import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";

interface GalleryNavProps {
  onPrev: () => void;
  onNext: () => void;
  hasPrev: boolean;
  hasNext: boolean;
}

export default function GalleryNav({ onPrev, onNext, hasPrev, hasNext }: GalleryNavProps) {
  return (
    <>
      {hasPrev && (
        <Button
          size="icon"
          variant="secondary"
          onClick={onPrev}
          className="absolute left-6 top-1/2 -translate-y-1/2 z-10 bg-background/80 backdrop-blur-sm rounded-full shadow-md"
          aria-label="Previous template"
        >
          <ChevronLeft className="h-5 w-5" />
        </Button>
      )}
      {hasNext && (
        <Button
          size="icon"
          variant="secondary"
          onClick={onNext}
          className="absolute right-6 top-1/2 -translate-y-1/2 z-10 bg-background/80 backdrop-blur-sm rounded-full shadow-md"
          aria-label="Next template"
        >
          <ChevronRight className="h-5 w-5" />
        </Button>
      )}
    </>
  );
}
