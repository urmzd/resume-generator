import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import EditorSection from "./EditorSection";

interface SummaryEditorProps {
  value: string;
  onChange: (s: string) => void;
}

export default function SummaryEditor({ value, onChange }: SummaryEditorProps) {
  return (
    <EditorSection title="Summary">
      <div>
        <Label>Summary</Label>
        <Textarea
          value={value}
          onChange={(e) => onChange(e.target.value)}
          rows={4}
          placeholder="Professional summary…"
        />
      </div>
    </EditorSection>
  );
}
