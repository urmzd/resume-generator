import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { SelectNative } from "@/components/ui/select-native";
import EditorSection from "./EditorSection";
import StringListEditor from "./StringListEditor";
import type { Layout } from "@/types";

interface LayoutEditorProps {
  value: Layout;
  onChange: (l: Layout) => void;
}

export default function LayoutEditor({ value, onChange }: LayoutEditorProps) {
  const set = (field: keyof Layout, v: unknown) =>
    onChange({ ...value, [field]: v });

  return (
    <EditorSection title="Layout">
      <div>
        <Label>Density</Label>
        <SelectNative
          value={value.density || ""}
          onChange={(e) => set("density", e.target.value)}
        >
          <option value="">Default</option>
          <option value="compact">Compact</option>
          <option value="normal">Normal</option>
          <option value="comfortable">Comfortable</option>
        </SelectNative>
      </div>
      <div>
        <Label>Typography</Label>
        <SelectNative
          value={value.typography || ""}
          onChange={(e) => set("typography", e.target.value)}
        >
          <option value="">Default</option>
          <option value="sans">Sans</option>
          <option value="serif">Serif</option>
          <option value="mono">Mono</option>
        </SelectNative>
      </div>
      <div>
        <Label>Header</Label>
        <SelectNative
          value={value.header || ""}
          onChange={(e) => set("header", e.target.value)}
        >
          <option value="">Default</option>
          <option value="left">Left</option>
          <option value="center">Center</option>
          <option value="compact">Compact</option>
        </SelectNative>
      </div>
      <div>
        <Label>Skill Columns</Label>
        <Input
          type="number"
          min={1}
          max={4}
          value={value.skill_columns || ""}
          onChange={(e) => set("skill_columns", parseInt(e.target.value) || 0)}
        />
      </div>
      <StringListEditor
        label="Section Order"
        items={value.sections || []}
        onChange={(sections) => set("sections", sections)}
        placeholder="Section name…"
      />
    </EditorSection>
  );
}
