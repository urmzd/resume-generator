import { useState, type ReactNode } from "react";
import { ChevronRight } from "lucide-react";
import {
  Collapsible,
  CollapsibleTrigger,
  CollapsibleContent,
} from "@/components/ui/collapsible";
import { cn } from "@/lib/utils";

interface EditorSectionProps {
  title: string;
  defaultOpen?: boolean;
  children: ReactNode;
}

export default function EditorSection({
  title,
  defaultOpen = false,
  children,
}: EditorSectionProps) {
  const [open, setOpen] = useState(defaultOpen);

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <CollapsibleTrigger className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm font-semibold hover:bg-accent">
        <ChevronRight
          className={cn("h-4 w-4 shrink-0 transition-transform", open && "rotate-90")}
        />
        {title}
      </CollapsibleTrigger>
      <CollapsibleContent>
        <div className="space-y-3 px-2 pb-3 pt-1">{children}</div>
      </CollapsibleContent>
    </Collapsible>
  );
}
