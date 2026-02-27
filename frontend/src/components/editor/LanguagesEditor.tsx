import { Plus, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import EditorSection from "./EditorSection";
import type { LanguageList, Language } from "@/types";

interface LanguagesEditorProps {
  value: LanguageList;
  onChange: (l: LanguageList) => void;
}

export default function LanguagesEditor({
  value,
  onChange,
}: LanguagesEditorProps) {
  const langs = value.languages || [];

  const updateLang = (i: number, lang: Language) => {
    const next = [...langs];
    next[i] = lang;
    onChange({ ...value, languages: next });
  };

  const addLang = () =>
    onChange({ ...value, languages: [...langs, { name: "" }] });

  const removeLang = (i: number) =>
    onChange({ ...value, languages: langs.filter((_, idx) => idx !== i) });

  return (
    <EditorSection title="Languages">
      {langs.map((lang, i) => (
        <div key={i} className="flex items-end gap-2">
          <div className="flex-1">
            <Label>Name</Label>
            <Input
              value={lang.name}
              onChange={(e) =>
                updateLang(i, { ...lang, name: e.target.value })
              }
            />
          </div>
          <div className="flex-1">
            <Label>Proficiency</Label>
            <Input
              value={lang.proficiency || ""}
              onChange={(e) =>
                updateLang(i, { ...lang, proficiency: e.target.value })
              }
            />
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className="shrink-0 h-9 w-9"
            onClick={() => removeLang(i)}
          >
            <X className="h-3.5 w-3.5" />
          </Button>
        </div>
      ))}
      <Button type="button" variant="outline" size="sm" className="w-full" onClick={addLang}>
        <Plus className="h-3.5 w-3.5 mr-1" />
        Add Language
      </Button>
    </EditorSection>
  );
}
