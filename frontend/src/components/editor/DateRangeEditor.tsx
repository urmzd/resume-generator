import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { DateRange } from "@/types";

interface DateRangeEditorProps {
  value: DateRange;
  onChange: (dr: DateRange) => void;
}

function toMonthInput(iso: string | undefined): string {
  if (!iso) return "";
  // Handle ISO 8601 date strings — extract YYYY-MM
  return iso.slice(0, 7);
}

function fromMonthInput(monthStr: string): string {
  if (!monthStr) return "";
  // month input gives YYYY-MM, we need a full ISO date
  return monthStr + "-01T00:00:00Z";
}

export default function DateRangeEditor({ value, onChange }: DateRangeEditorProps) {
  const isPresent = !value.end;

  return (
    <div className="space-y-1.5">
      <div className="grid grid-cols-2 gap-2">
        <div>
          <Label>Start</Label>
          <Input
            type="month"
            value={toMonthInput(value.start)}
            onChange={(e) =>
              onChange({ ...value, start: fromMonthInput(e.target.value) })
            }
          />
        </div>
        <div>
          <Label>End</Label>
          <Input
            type="month"
            value={toMonthInput(value.end)}
            onChange={(e) =>
              onChange({
                ...value,
                end: e.target.value ? fromMonthInput(e.target.value) : undefined,
              })
            }
            disabled={isPresent}
          />
        </div>
      </div>
      <label className="flex items-center gap-2 text-sm">
        <input
          type="checkbox"
          checked={isPresent}
          onChange={(e) =>
            onChange({
              ...value,
              end: e.target.checked ? undefined : fromMonthInput(""),
            })
          }
          className="rounded border-input"
        />
        Present (no end date)
      </label>
    </div>
  );
}
