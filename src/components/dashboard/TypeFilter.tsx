import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";

interface TypeFilterProps {
  value: string;
  onChange: (value: string) => void;
}

export function TypeFilter({ value, onChange }: TypeFilterProps) {
  return (
    <ToggleGroup
      type="single"
      value={value}
      onValueChange={(v) => onChange(v || "all")}
      className="border rounded-md"
    >
      <ToggleGroupItem value="all" className="text-xs px-3">All</ToggleGroupItem>
      <ToggleGroupItem value="movie" className="text-xs px-3">Movies</ToggleGroupItem>
      <ToggleGroupItem value="tv" className="text-xs px-3">TV</ToggleGroupItem>
    </ToggleGroup>
  );
}
