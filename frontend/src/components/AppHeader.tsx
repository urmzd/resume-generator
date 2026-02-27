import { Download, FileText, FolderOpen, Pencil } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import ThemeToggle from "./ThemeToggle";
import type { TemplateInfo } from "@/types";

interface AppHeaderProps {
  selectedTemplate: TemplateInfo | null;
  onSavePdf: () => void;
  onSaveNative: () => void;
  onChangeFile: () => void;
  editorOpen?: boolean;
  onToggleEditor?: () => void;
}

export default function AppHeader({
  selectedTemplate,
  onSavePdf,
  onSaveNative,
  onChangeFile,
  editorOpen,
  onToggleEditor,
}: AppHeaderProps) {
  return (
    <header className="flex items-center justify-between border-b px-4 py-2">
      <div className="flex items-center gap-3 min-w-0">
        <p className="text-sm font-medium whitespace-nowrap">
          Resume Generator{" "}
          <span className="text-muted-foreground font-normal">by Urmzd (@urmzd.com)</span>
        </p>
      </div>

      <div className="flex items-center gap-1">
        <ThemeToggle />

        {onToggleEditor && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant={editorOpen ? "default" : "outline"}
                size="sm"
                onClick={onToggleEditor}
              >
                <Pencil className="h-4 w-4 mr-1" />
                Edit
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              {editorOpen ? "Close editor" : "Open resume editor"}
            </TooltipContent>
          </Tooltip>
        )}

        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="outline" size="sm" onClick={onSavePdf}>
              <Download className="h-4 w-4 mr-1" />
              PDF
            </Button>
          </TooltipTrigger>
          <TooltipContent>Save as PDF</TooltipContent>
        </Tooltip>

        {selectedTemplate?.format === "docx" && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline" size="sm" onClick={onSaveNative}>
                <FileText className="h-4 w-4 mr-1" />
                DOCX
              </Button>
            </TooltipTrigger>
            <TooltipContent>Save as DOCX</TooltipContent>
          </Tooltip>
        )}

        <Separator orientation="vertical" className="h-6 mx-1" />

        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="ghost" size="sm" onClick={onChangeFile}>
              <FolderOpen className="h-4 w-4 mr-1" />
              Change File
            </Button>
          </TooltipTrigger>
          <TooltipContent>Load a different resume file</TooltipContent>
        </Tooltip>
      </div>
    </header>
  );
}
