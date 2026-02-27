import { Plus, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import EditorSection from "./EditorSection";
import DateRangeEditor from "./DateRangeEditor";
import LocationEditor from "./LocationEditor";
import StringListEditor from "./StringListEditor";
import type { ExperienceList, Experience } from "@/types";

interface ExperienceEditorProps {
  value: ExperienceList;
  onChange: (e: ExperienceList) => void;
}

export default function ExperienceEditor({
  value,
  onChange,
}: ExperienceEditorProps) {
  const positions = value.positions || [];

  const updatePos = (i: number, pos: Experience) => {
    const next = [...positions];
    next[i] = pos;
    onChange({ ...value, positions: next });
  };

  const addPos = () =>
    onChange({
      ...value,
      positions: [
        ...positions,
        {
          company: "",
          title: "",
          dates: { start: "" },
        },
      ],
    });

  const removePos = (i: number) =>
    onChange({ ...value, positions: positions.filter((_, idx) => idx !== i) });

  return (
    <EditorSection title="Experience">
      {positions.map((pos, i) => (
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
              onClick={() => removePos(i)}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
          <div>
            <Label>Title</Label>
            <Input
              value={pos.title}
              onChange={(e) => updatePos(i, { ...pos, title: e.target.value })}
            />
          </div>
          <div>
            <Label>Company</Label>
            <Input
              value={pos.company}
              onChange={(e) => updatePos(i, { ...pos, company: e.target.value })}
            />
          </div>
          <DateRangeEditor
            value={pos.dates}
            onChange={(dates) => updatePos(i, { ...pos, dates })}
          />
          <div className="space-y-1.5">
            <Label>Location</Label>
            <LocationEditor
              value={pos.location || { city: "" }}
              onChange={(location) => updatePos(i, { ...pos, location })}
            />
          </div>
          <StringListEditor
            label="Technologies"
            items={pos.technologies || []}
            onChange={(technologies) => updatePos(i, { ...pos, technologies })}
            placeholder="Technology…"
          />
          <StringListEditor
            label="Highlights"
            items={pos.highlights || []}
            onChange={(highlights) => updatePos(i, { ...pos, highlights })}
            placeholder="Highlight…"
          />
        </div>
      ))}
      <Button type="button" variant="outline" size="sm" className="w-full" onClick={addPos}>
        <Plus className="h-3.5 w-3.5 mr-1" />
        Add Position
      </Button>
    </EditorSection>
  );
}
