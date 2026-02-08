import FileDropZone from "@/components/FileDropZone";

interface DropPageProps {
  onFile: (data: Uint8Array, format: string) => void;
  onOpenNative: () => void;
  error: string | null;
}

export default function DropPage({ onFile, onOpenNative, error }: DropPageProps) {
  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <div className="w-full max-w-md">
        <FileDropZone onFile={onFile} onOpenNative={onOpenNative} error={error} />
      </div>
    </div>
  );
}
