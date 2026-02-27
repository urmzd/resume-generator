import { useState, useCallback, useEffect } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import AppHeader from "@/components/AppHeader";
import PdfViewer from "@/components/PdfViewer";
import ThumbnailStrip from "@/components/ThumbnailStrip";
import GalleryContainer from "@/containers/GalleryContainer";
import ResumeEditorPanel from "@/components/editor/ResumeEditorPanel";
import { useResumeEditor } from "@/hooks/useResumeEditor";
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
  const [editorOpen, setEditorOpen] = useState(false);
  const editor = useResumeEditor();

  useEffect(() => {
    editor.fetchResume();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleToggleEditor = useCallback(() => {
    setEditorOpen((prev) => !prev);
  }, []);

  return (
    <GalleryContainer templates={templates} onError={onError}>
      {({
        selectedIndex,
        selectedTemplate,
        pdfUrl,
        isLoading,
        pageCount,
        onSelectTemplate,
        onSavePdf,
        onSaveNative,
        getCachedUrl,
        invalidateCache,
      }) => {
        const handleApply = async () => {
          const ok = await editor.applyChanges();
          if (ok) {
            invalidateCache();
          }
        };

        return (
          <div className="flex h-screen flex-col">
            <AppHeader
              selectedTemplate={selectedTemplate}
              onSavePdf={onSavePdf}
              onSaveNative={onSaveNative}
              onChangeFile={onReset}
              editorOpen={editorOpen}
              onToggleEditor={handleToggleEditor}
            />

            {error && (
              <div className="px-4 pt-2">
                <Alert variant="destructive">
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              </div>
            )}

            {pageCount !== null && pageCount > 1 && (
              <div className="px-4 pt-2">
                <Alert>
                  <AlertDescription>
                    Resume is {pageCount} pages. Consider using a compact density or reducing content to fit on one page.
                  </AlertDescription>
                </Alert>
              </div>
            )}

            <div className="flex-1 min-h-0 flex flex-col">
              <div className="flex-1 min-h-0 flex">
                {editorOpen && editor.draft && (
                  <div className="w-[380px] shrink-0 overflow-hidden">
                    <ResumeEditorPanel
                      draft={editor.draft}
                      dirty={editor.dirty}
                      loading={editor.loading}
                      saving={editor.saving}
                      errors={editor.errors}
                      onDraftChange={editor.setDraft}
                      onApply={handleApply}
                      onSave={editor.saveToFile}
                    />
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
            </div>
          </div>
        );
      }}
    </GalleryContainer>
  );
}
