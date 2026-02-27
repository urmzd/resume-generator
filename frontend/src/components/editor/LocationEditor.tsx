import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { Location } from "@/types";

interface LocationEditorProps {
  value: Location;
  onChange: (loc: Location) => void;
}

export default function LocationEditor({ value, onChange }: LocationEditorProps) {
  const set = (field: keyof Location, v: string | boolean) =>
    onChange({ ...value, [field]: v });

  return (
    <div className="space-y-1.5">
      <div className="grid grid-cols-2 gap-2">
        <div>
          <Label>City</Label>
          <Input value={value.city || ""} onChange={(e) => set("city", e.target.value)} />
        </div>
        <div>
          <Label>State</Label>
          <Input value={value.state || ""} onChange={(e) => set("state", e.target.value)} />
        </div>
      </div>
      <div className="grid grid-cols-2 gap-2">
        <div>
          <Label>Province</Label>
          <Input value={value.province || ""} onChange={(e) => set("province", e.target.value)} />
        </div>
        <div>
          <Label>Country</Label>
          <Input value={value.country || ""} onChange={(e) => set("country", e.target.value)} />
        </div>
      </div>
      <label className="flex items-center gap-2 text-sm">
        <input
          type="checkbox"
          checked={value.remote || false}
          onChange={(e) => set("remote", e.target.checked)}
          className="rounded border-input"
        />
        Remote
      </label>
    </div>
  );
}
