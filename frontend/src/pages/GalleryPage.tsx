import { Alert, AlertDescription } from "@/components/ui/alert";
import AppHeader from "@/components/AppHeader";
import PdfViewer from "@/components/PdfViewer";
import ThumbnailStrip from "@/components/ThumbnailStrip";
import GalleryContainer from "@/containers/GalleryContainer";
import type { TemplateInfo } from "@/types";

interface GalleryPageProps {
  templates: TemplateInfo[];
  error: string | null;
  onError: (msg: string) => void;
  onReset: () => void;
}

export default function GalleryPage({
  templates,
  error,
  onError,
  onReset,
}: GalleryPageProps) {
  return (
    <GalleryContainer templates={templates} onError={onError}>
      {({
        selectedIndex,
        selectedTemplate,
        pdfUrl,
        isLoading,
        onSelectTemplate,
        onSavePdf,
        onSaveNative,
        getCachedUrl,
      }) => (
        <div className="flex h-screen flex-col">
          <AppHeader
            selectedTemplate={selectedTemplate}
            onSavePdf={onSavePdf}
            onSaveNative={onSaveNative}
            onChangeFile={onReset}
          />

          {error && (
            <div className="px-4 pt-2">
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            </div>
          )}

          <div className="flex-1 min-h-0 p-4 flex flex-col gap-2">
            <div className="flex-1 min-h-0">
              <PdfViewer pdfUrl={pdfUrl} isLoading={isLoading} />
            </div>
            <ThumbnailStrip
              templates={templates}
              selectedIndex={selectedIndex}
              onSelect={onSelectTemplate}
              getCachedUrl={getCachedUrl}
            />
          </div>
        </div>
      )}
    </GalleryContainer>
  );
}
