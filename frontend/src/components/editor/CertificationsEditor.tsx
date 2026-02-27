import { Plus, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import EditorSection from "./EditorSection";
import type { Certifications, Certification } from "@/types";

interface CertificationsEditorProps {
  value: Certifications;
  onChange: (c: Certifications) => void;
}

export default function CertificationsEditor({
  value,
  onChange,
}: CertificationsEditorProps) {
  const items = value.items || [];

  const updateItem = (i: number, cert: Certification) => {
    const next = [...items];
    next[i] = cert;
    onChange({ ...value, items: next });
  };

  const addItem = () =>
    onChange({ ...value, items: [...items, { name: "" }] });

  const removeItem = (i: number) =>
    onChange({ ...value, items: items.filter((_, idx) => idx !== i) });

  return (
    <EditorSection title="Certifications">
      {items.map((cert, i) => (
        <div key={i} className="space-y-1.5 rounded-md border p-2">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-muted-foreground">
              #{i + 1}
            </span>
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={() => removeItem(i)}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
          <div>
            <Label>Name</Label>
            <Input
              value={cert.name}
              onChange={(e) => updateItem(i, { ...cert, name: e.target.value })}
            />
          </div>
          <div>
            <Label>Issuer</Label>
            <Input
              value={cert.issuer || ""}
              onChange={(e) => updateItem(i, { ...cert, issuer: e.target.value })}
            />
          </div>
          <div>
            <Label>Notes</Label>
            <Input
              value={cert.notes || ""}
              onChange={(e) => updateItem(i, { ...cert, notes: e.target.value })}
            />
          </div>
        </div>
      ))}
      <Button type="button" variant="outline" size="sm" className="w-full" onClick={addItem}>
        <Plus className="h-3.5 w-3.5 mr-1" />
        Add Certification
      </Button>
    </EditorSection>
  );
}
