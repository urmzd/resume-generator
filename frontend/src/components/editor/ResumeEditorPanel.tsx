import { Save, Check, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import ContactEditor from "./ContactEditor";
import SummaryEditor from "./SummaryEditor";
import CertificationsEditor from "./CertificationsEditor";
import SkillsEditor from "./SkillsEditor";
import ExperienceEditor from "./ExperienceEditor";
import EducationEditor from "./EducationEditor";
import ProjectsEditor from "./ProjectsEditor";
import LanguagesEditor from "./LanguagesEditor";
import LayoutEditor from "./LayoutEditor";
import type { Resume, ValidationError } from "@/types";

interface ResumeEditorPanelProps {
  draft: Resume;
  dirty: boolean;
  loading: boolean;
  saving: boolean;
  errors: ValidationError[];
  onDraftChange: (r: Resume) => void;
  onApply: () => void;
  onSave: () => void;
}

export default function ResumeEditorPanel({
  draft,
  dirty,
  loading,
  saving,
  errors,
  onDraftChange,
  onApply,
  onSave,
}: ResumeEditorPanelProps) {
  return (
    <div className="flex h-full flex-col border-r">
      <div className="flex items-center gap-2 border-b px-3 py-2">
        <Button
          size="sm"
          onClick={onApply}
          disabled={!dirty || loading}
        >
          <Check className="h-3.5 w-3.5 mr-1" />
          Apply
        </Button>
        <Button
          size="sm"
          variant="outline"
          onClick={onSave}
          disabled={saving}
        >
          <Save className="h-3.5 w-3.5 mr-1" />
          Save to File
        </Button>
      </div>

      {errors.length > 0 && (
        <div className="px-3 pt-2">
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>
              {errors.map((e, i) => (
                <div key={i}>
                  <strong>{e.field}</strong>: {e.message}
                </div>
              ))}
            </AlertDescription>
          </Alert>
        </div>
      )}

      <div className="flex-1 overflow-y-auto px-2 py-2 space-y-1">
        <ContactEditor
          value={draft.contact}
          onChange={(contact) => onDraftChange({ ...draft, contact })}
        />
        <SummaryEditor
          value={draft.summary || ""}
          onChange={(summary) => onDraftChange({ ...draft, summary })}
        />
        <ExperienceEditor
          value={draft.experience}
          onChange={(experience) => onDraftChange({ ...draft, experience })}
        />
        <EducationEditor
          value={draft.education}
          onChange={(education) => onDraftChange({ ...draft, education })}
        />
        <SkillsEditor
          value={draft.skills}
          onChange={(skills) => onDraftChange({ ...draft, skills })}
        />
        <ProjectsEditor
          value={draft.projects || { projects: [] }}
          onChange={(projects) => onDraftChange({ ...draft, projects })}
        />
        <CertificationsEditor
          value={draft.certifications || { items: [] }}
          onChange={(certifications) =>
            onDraftChange({ ...draft, certifications })
          }
        />
        <LanguagesEditor
          value={draft.languages || { languages: [] }}
          onChange={(languages) => onDraftChange({ ...draft, languages })}
        />
        <LayoutEditor
          value={draft.layout || {}}
          onChange={(layout) => onDraftChange({ ...draft, layout })}
        />
      </div>
    </div>
  );
}
