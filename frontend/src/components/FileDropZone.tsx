import { useCallback, useRef, useState } from "react";
import { Upload } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface FileDropZoneProps {
  onFile: (data: Uint8Array, format: string) => void;
  onOpenNative: () => void;
  error?: string | null;
}

const ACCEPTED_EXTENSIONS = ["yml", "yaml", "json", "toml"];

function getExtension(filename: string): string {
  const dot = filename.lastIndexOf(".");
  return dot >= 0 ? filename.slice(dot + 1).toLowerCase() : "";
}

export default function FileDropZone({ onFile, onOpenNative, error }: FileDropZoneProps) {
  const [dragOver, setDragOver] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFile = useCallback(
    (file: File) => {
      const ext = getExtension(file.name);
      if (!ACCEPTED_EXTENSIONS.includes(ext)) {
        onFile(new Uint8Array(), `__invalid:${ext}`);
        return;
      }
      const reader = new FileReader();
      reader.onload = () => {
        onFile(new Uint8Array(reader.result as ArrayBuffer), ext);
      };
      reader.readAsArrayBuffer(file);
    },
    [onFile]
  );

  const onDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setDragOver(false);
      if (e.dataTransfer.files.length) handleFile(e.dataTransfer.files[0]);
    },
    [handleFile]
  );

  return (
    <div className="space-y-4">
      <div className="text-center space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">Resume Generator</h1>
        <p className="text-muted-foreground">Drop your resume file or use the native file picker</p>
      </div>

      <Card
        className={`cursor-pointer transition-colors border-dashed border-2 ${
          dragOver ? "border-primary bg-accent" : "hover:border-primary/50"
        }`}
        onClick={onOpenNative}
        onDragOver={(e) => {
          e.preventDefault();
          setDragOver(true);
        }}
        onDragLeave={() => setDragOver(false)}
        onDrop={onDrop}
        role="button"
        tabIndex={0}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            onOpenNative();
          }
        }}
      >
        <CardContent className="flex flex-col items-center justify-center py-12 space-y-3">
          <Upload className="h-10 w-10 text-muted-foreground/50" />
          <div className="text-center">
            <p className="text-sm text-muted-foreground">
              Drag & drop <strong>.yml</strong>, <strong>.yaml</strong>,{" "}
              <strong>.json</strong>, or <strong>.toml</strong>
            </p>
            <p className="text-xs text-muted-foreground/70 mt-1">or click to open file picker</p>
          </div>
          <input
            ref={inputRef}
            type="file"
            accept=".yml,.yaml,.json,.toml"
            className="hidden"
            onChange={(e) => {
              if (e.target.files?.length) handleFile(e.target.files[0]);
            }}
          />
        </CardContent>
      </Card>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <p className="text-center text-xs text-muted-foreground/60">
        Got an issue? Reach out to{" "}
        <a href="mailto:hello@urmzd.com" className="underline hover:text-muted-foreground">hello@urmzd.com</a>
        {" "}or{" "}
        <a href="https://github.com/urmzd/resume-generator/issues" target="_blank" rel="noopener noreferrer" className="underline hover:text-muted-foreground">open an issue on GitHub</a>.
      </p>
    </div>
  );
}
