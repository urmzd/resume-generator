import { Plus, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import EditorSection from "./EditorSection";
import DateRangeEditor from "./DateRangeEditor";
import LocationEditor from "./LocationEditor";
import StringListEditor from "./StringListEditor";
import type { EducationList, Education } from "@/types";

interface EducationEditorProps {
  value: EducationList;
  onChange: (e: EducationList) => void;
}

export default function EducationEditor({
  value,
  onChange,
}: EducationEditorProps) {
  const institutions = value.institutions || [];

  const updateInst = (i: number, inst: Education) => {
    const next = [...institutions];
    next[i] = inst;
    onChange({ ...value, institutions: next });
  };

  const addInst = () =>
    onChange({
      ...value,
      institutions: [
        ...institutions,
        {
          institution: "",
          degree: { name: "" },
          dates: { start: "" },
        },
      ],
    });

  const removeInst = (i: number) =>
    onChange({
      ...value,
      institutions: institutions.filter((_, idx) => idx !== i),
    });

  return (
    <EditorSection title="Education">
      {institutions.map((inst, i) => (
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
              onClick={() => removeInst(i)}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
          <div>
            <Label>Institution</Label>
            <Input
              value={inst.institution}
              onChange={(e) =>
                updateInst(i, { ...inst, institution: e.target.value })
              }
            />
          </div>
          <div>
            <Label>Degree</Label>
            <Input
              value={inst.degree.name}
              onChange={(e) =>
                updateInst(i, {
                  ...inst,
                  degree: { ...inst.degree, name: e.target.value },
                })
              }
            />
          </div>
          <DateRangeEditor
            value={inst.dates}
            onChange={(dates) => updateInst(i, { ...inst, dates })}
          />
          <div className="space-y-1.5">
            <Label>Location</Label>
            <LocationEditor
              value={inst.location || { city: "" }}
              onChange={(location) => updateInst(i, { ...inst, location })}
            />
          </div>
          <div className="grid grid-cols-2 gap-2">
            <div>
              <Label>GPA</Label>
              <Input
                value={inst.gpa?.gpa || ""}
                onChange={(e) =>
                  updateInst(i, {
                    ...inst,
                    gpa: { ...inst.gpa, gpa: e.target.value },
                  })
                }
              />
            </div>
            <div>
              <Label>Max GPA</Label>
              <Input
                value={inst.gpa?.max_gpa || ""}
                onChange={(e) =>
                  updateInst(i, {
                    ...inst,
                    gpa: { ...inst.gpa, max_gpa: e.target.value },
                  })
                }
              />
            </div>
          </div>
          <StringListEditor
            label="Descriptions"
            items={inst.degree.descriptions || []}
            onChange={(descriptions) =>
              updateInst(i, {
                ...inst,
                degree: { ...inst.degree, descriptions },
              })
            }
            placeholder="Description…"
          />
        </div>
      ))}
      <Button type="button" variant="outline" size="sm" className="w-full" onClick={addInst}>
        <Plus className="h-3.5 w-3.5 mr-1" />
        Add Institution
      </Button>
    </EditorSection>
  );
}
