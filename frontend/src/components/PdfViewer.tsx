import { Loader2 } from "lucide-react";

interface PdfViewerProps {
  pdfUrl: string | null;
  isLoading: boolean;
}

export default function PdfViewer({ pdfUrl, isLoading }: PdfViewerProps) {
  if (pdfUrl) {
    return (
      <div className="relative h-full w-full rounded-lg border bg-card overflow-hidden shadow-lg">
        <iframe src={`${pdfUrl}#page=1&view=FitH&toolbar=0&navpanes=0`} className="h-full w-full" title="PDF Preview" />
        {isLoading && (
          <div className="absolute inset-0 flex items-center justify-center bg-background/80 backdrop-blur-sm">
            <div className="flex flex-col items-center gap-2">
              <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
              <p className="text-sm text-muted-foreground">Generating PDF...</p>
            </div>
          </div>
        )}
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="flex h-full w-full items-center justify-center rounded-lg border bg-muted/50 shadow-lg">
        <div className="flex flex-col items-center gap-2">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          <p className="text-sm text-muted-foreground">Generating PDF...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full w-full items-center justify-center rounded-lg border border-dashed shadow-lg">
      <p className="text-sm text-muted-foreground">Select a template to preview</p>
    </div>
  );
}
