import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export type SortOption = "title-asc" | "rating-desc" | "year-desc" | "popularity-desc";

interface SortSelectProps {
  value: SortOption;
  onChange: (value: SortOption) => void;
}

export function SortSelect({ value, onChange }: SortSelectProps) {
  return (
    <Select value={value} onValueChange={(v) => onChange(v as SortOption)}>
      <SelectTrigger className="w-[170px]">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="title-asc">Title A–Z</SelectItem>
        <SelectItem value="rating-desc">Rating High–Low</SelectItem>
        <SelectItem value="year-desc">Year New–Old</SelectItem>
        <SelectItem value="popularity-desc">Popularity</SelectItem>
      </SelectContent>
    </Select>
  );
}
