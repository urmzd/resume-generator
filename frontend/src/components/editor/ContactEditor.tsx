import { Plus, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import EditorSection from "./EditorSection";
import LocationEditor from "./LocationEditor";
import LinkEditor from "./LinkEditor";
import type { Contact, Link, Location } from "@/types";

interface ContactEditorProps {
  value: Contact;
  onChange: (c: Contact) => void;
}

export default function ContactEditor({ value, onChange }: ContactEditorProps) {
  const set = (field: keyof Contact, v: unknown) =>
    onChange({ ...value, [field]: v });

  const links = value.links || [];

  const updateLink = (i: number, link: Link) => {
    const next = [...links];
    next[i] = link;
    set("links", next);
  };

  const addLink = () => set("links", [...links, { uri: "", label: "" }]);
  const removeLink = (i: number) => set("links", links.filter((_, idx) => idx !== i));

  return (
    <EditorSection title="Contact" defaultOpen>
      <div>
        <Label>Name</Label>
        <Input value={value.name} onChange={(e) => set("name", e.target.value)} />
      </div>
      <div>
        <Label>Email</Label>
        <Input value={value.email} onChange={(e) => set("email", e.target.value)} />
      </div>
      <div>
        <Label>Phone</Label>
        <Input
          value={value.phone || ""}
          onChange={(e) => set("phone", e.target.value)}
        />
      </div>

      <div className="space-y-1.5">
        <Label>Location</Label>
        <LocationEditor
          value={value.location || { city: "" }}
          onChange={(loc: Location) => set("location", loc)}
        />
      </div>

      <div className="space-y-1.5">
        <Label>Links</Label>
        {links.map((link, i) => (
          <div key={i} className="flex items-end gap-1">
            <div className="flex-1">
              <LinkEditor value={link} onChange={(l) => updateLink(i, l)} />
            </div>
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="shrink-0 h-9 w-9"
              onClick={() => removeLink(i)}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
        ))}
        <Button type="button" variant="outline" size="sm" className="w-full" onClick={addLink}>
          <Plus className="h-3.5 w-3.5 mr-1" />
          Add Link
        </Button>
      </div>
    </EditorSection>
  );
}
