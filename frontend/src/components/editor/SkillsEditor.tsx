import { Plus, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import EditorSection from "./EditorSection";
import type { Skills, SkillCategory } from "@/types";

interface SkillsEditorProps {
  value: Skills;
  onChange: (s: Skills) => void;
}

export default function SkillsEditor({ value, onChange }: SkillsEditorProps) {
  const cats = value.categories || [];

  const updateCat = (i: number, cat: SkillCategory) => {
    const next = [...cats];
    next[i] = cat;
    onChange({ ...value, categories: next });
  };

  const addCat = () =>
    onChange({
      ...value,
      categories: [...cats, { category: "", items: [] }],
    });

  const removeCat = (i: number) =>
    onChange({ ...value, categories: cats.filter((_, idx) => idx !== i) });

  return (
    <EditorSection title="Skills">
      {cats.map((cat, i) => (
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
              onClick={() => removeCat(i)}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
          <div>
            <Label>Category</Label>
            <Input
              value={cat.category}
              onChange={(e) => updateCat(i, { ...cat, category: e.target.value })}
            />
          </div>
          <div>
            <Label>Items (comma-separated)</Label>
            <Input
              value={cat.items.join(", ")}
              onChange={(e) =>
                updateCat(i, {
                  ...cat,
                  items: e.target.value.split(",").map((s) => s.trim()),
                })
              }
              placeholder="React, TypeScript, Go…"
            />
          </div>
        </div>
      ))}
      <Button type="button" variant="outline" size="sm" className="w-full" onClick={addCat}>
        <Plus className="h-3.5 w-3.5 mr-1" />
        Add Category
      </Button>
    </EditorSection>
  );
}
