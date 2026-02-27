import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { Link } from "@/types";

interface LinkEditorProps {
  value: Link;
  onChange: (link: Link) => void;
}

export default function LinkEditor({ value, onChange }: LinkEditorProps) {
  return (
    <div className="grid grid-cols-2 gap-2">
      <div>
        <Label>URI</Label>
        <Input
          value={value.uri || ""}
          onChange={(e) => onChange({ ...value, uri: e.target.value })}
          placeholder="https://…"
        />
      </div>
      <div>
        <Label>Label</Label>
        <Input
          value={value.label || ""}
          onChange={(e) => onChange({ ...value, label: e.target.value })}
          placeholder="Optional label"
        />
      </div>
    </div>
  );
}
