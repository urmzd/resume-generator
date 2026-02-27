import { Plus, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import EditorSection from "./EditorSection";
import DateRangeEditor from "./DateRangeEditor";
import LinkEditor from "./LinkEditor";
import StringListEditor from "./StringListEditor";
import type { ProjectList, Project } from "@/types";

interface ProjectsEditorProps {
  value: ProjectList;
  onChange: (p: ProjectList) => void;
}

export default function ProjectsEditor({
  value,
  onChange,
}: ProjectsEditorProps) {
  const projects = value.projects || [];

  const updateProj = (i: number, proj: Project) => {
    const next = [...projects];
    next[i] = proj;
    onChange({ ...value, projects: next });
  };

  const addProj = () =>
    onChange({
      ...value,
      projects: [...projects, { name: "" }],
    });

  const removeProj = (i: number) =>
    onChange({ ...value, projects: projects.filter((_, idx) => idx !== i) });

  return (
    <EditorSection title="Projects">
      {projects.map((proj, i) => (
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
              onClick={() => removeProj(i)}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
          <div>
            <Label>Name</Label>
            <Input
              value={proj.name}
              onChange={(e) => updateProj(i, { ...proj, name: e.target.value })}
            />
          </div>
          <LinkEditor
            value={proj.link || { uri: "" }}
            onChange={(link) => updateProj(i, { ...proj, link })}
          />
          {proj.dates !== undefined && (
            <DateRangeEditor
              value={proj.dates}
              onChange={(dates) => updateProj(i, { ...proj, dates })}
            />
          )}
          {proj.dates === undefined && (
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() =>
                updateProj(i, { ...proj, dates: { start: "" } })
              }
            >
              Add Dates
            </Button>
          )}
          <StringListEditor
            label="Technologies"
            items={proj.technologies || []}
            onChange={(technologies) => updateProj(i, { ...proj, technologies })}
            placeholder="Technology…"
          />
          <StringListEditor
            label="Highlights"
            items={proj.highlights || []}
            onChange={(highlights) => updateProj(i, { ...proj, highlights })}
            placeholder="Highlight…"
          />
        </div>
      ))}
      <Button type="button" variant="outline" size="sm" className="w-full" onClick={addProj}>
        <Plus className="h-3.5 w-3.5 mr-1" />
        Add Project
      </Button>
    </EditorSection>
  );
}
