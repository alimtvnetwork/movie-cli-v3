import { useState } from "react";
import { DashboardLayout } from "@/components/dashboard/DashboardLayout";
import { DashboardHeader } from "@/components/dashboard/DashboardHeader";
import { MediaGrid } from "@/components/dashboard/MediaGrid";
import { SearchBar } from "@/components/dashboard/SearchBar";
import { GenreFilter } from "@/components/dashboard/GenreFilter";
import { TypeFilter } from "@/components/dashboard/TypeFilter";
import { SortSelect } from "@/components/dashboard/SortSelect";
import { StatsPanel } from "@/components/dashboard/StatsPanel";
import { useMediaFilters } from "@/components/dashboard/useMediaFilters";
import { mockMedia } from "@/data/mock-media";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp, BarChart3, RotateCcw } from "lucide-react";
import { safeLocalGet, safeLocalSet } from "@/lib/media-utils";
import { toast } from "sonner";

const Index = () => {
  const [statsOpen, setStatsOpen] = useState(() => safeLocalGet("statsPanel.open", "true") === "true");

  const handleStatsToggle = (open: boolean) => {
    setStatsOpen(open);
    safeLocalSet("statsPanel.open", String(open));
  };

  const {
    search, setSearch,
    genre, setGenre,
    type, setType,
    sort, setSort,
    allGenres,
    filtered,
    hasActiveFilters,
    resetFilters,
  } = useMediaFilters(mockMedia);

  const handleReset = () => {
    resetFilters();
    toast.success("Filters reset");
  };

  return (
    <DashboardLayout>
      <DashboardHeader media={mockMedia} />

      <Collapsible open={statsOpen} onOpenChange={handleStatsToggle}>
        <div className="flex items-center justify-between">
          <CollapsibleTrigger asChild>
            <Button variant="ghost" size="sm" className="gap-2 text-muted-foreground">
              <BarChart3 className="h-4 w-4" />
              Statistics
              {statsOpen ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
            </Button>
          </CollapsibleTrigger>
        </div>
        <CollapsibleContent className="mt-2">
          <StatsPanel media={mockMedia} />
        </CollapsibleContent>
      </Collapsible>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:flex-wrap">
        <div className="flex-1 min-w-[200px]">
          <SearchBar value={search} onChange={setSearch} />
        </div>
        <GenreFilter genres={allGenres} value={genre} onChange={setGenre} />
        <TypeFilter value={type} onChange={setType} />
        <SortSelect value={sort} onChange={setSort} />
        {hasActiveFilters && (
          <Button variant="ghost" size="sm" className="gap-1.5 text-muted-foreground" onClick={handleReset}>
            <RotateCcw className="h-3.5 w-3.5" />
            Reset
          </Button>
        )}
      </div>

      <MediaGrid media={filtered} />
    </DashboardLayout>
  );
};

export default Index;
