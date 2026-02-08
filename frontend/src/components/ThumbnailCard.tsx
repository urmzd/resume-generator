import { FileText } from "lucide-react";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import type { TemplateInfo } from "@/types";

interface ThumbnailCardProps {
  template: TemplateInfo;
  isActive: boolean;
  previewUrl: string | null;
  onClick: () => void;
}

export default function ThumbnailCard({ template, isActive, previewUrl, onClick }: ThumbnailCardProps) {
  return (
    <Card
      role="tab"
      aria-selected={isActive}
      tabIndex={0}
      onClick={onClick}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          onClick();
        }
      }}
      className={cn(
        "flex flex-col items-center gap-1 p-1.5 cursor-pointer transition-all shrink-0",
        "hover:bg-accent hover:text-accent-foreground",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
        isActive && "ring-2 ring-primary"
      )}
    >
      {/* Mini preview */}
      <div className="w-20 h-24 rounded border bg-muted/50 overflow-hidden relative">
        {previewUrl ? (
          <iframe
            src={`${previewUrl}#toolbar=0&navpanes=0`}
            title={`${template.displayName} preview`}
            className="pointer-events-none absolute top-0 left-0"
            style={{
              width: "612px",
              height: "792px",
              transform: "scale(0.13)",
              transformOrigin: "top left",
            }}
            tabIndex={-1}
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center">
            <FileText className="h-6 w-6 text-muted-foreground/50" />
          </div>
        )}
      </div>

      <div className="flex items-center gap-1">
        <span className="text-[10px] font-medium whitespace-nowrap">{template.displayName}</span>
        <Badge
          variant={
            template.format === "html"
              ? "default"
              : template.format === "docx"
                ? "secondary"
                : "outline"
          }
          className="text-[9px] px-1 py-0 leading-tight"
        >
          {template.format.toUpperCase()}
        </Badge>
      </div>
    </Card>
  );
}
