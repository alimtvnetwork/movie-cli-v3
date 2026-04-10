import { useMemo, useState } from "react";
import type { MediaItem } from "@/types/media";
import type { SortOption } from "./SortSelect";

const VALID_SORTS: SortOption[] = ["title-asc", "rating-desc", "year-desc", "popularity-desc"];

export function useMediaFilters(media: MediaItem[]) {
  const [search, setSearch] = useState("");
  const [genre, setGenre] = useState("all");
  const [type, setType] = useState("all");
  const [sort, setSort] = useState<SortOption>("title-asc");

  const handleSortChange = (value: SortOption) => {
    if (VALID_SORTS.includes(value)) {
      setSort(value);
    }
  };

  const allGenres = useMemo(() => {
    const set = new Set<string>();
    media.forEach((m) => (m.genres ?? []).forEach((g) => set.add(g)));
    return Array.from(set).sort();
  }, [media]);

  const filtered = useMemo(() => {
    let result = media;

    if (search) {
      const q = search.toLowerCase();
      result = result.filter((m) =>
        m.title.toLowerCase().includes(q) ||
        (m.cleanTitle?.toLowerCase().includes(q)) ||
        (m.director?.toLowerCase().includes(q)) ||
        (m.cast ?? []).some((c) => c.toLowerCase().includes(q))
      );
    }

    if (genre !== "all") {
      result = result.filter((m) => (m.genres ?? []).includes(genre));
    }

    if (type !== "all") {
      result = result.filter((m) => m.type === type);
    }

    result = [...result].sort((a, b) => {
      switch (sort) {
        case "title-asc":
          return a.title.localeCompare(b.title);
        case "rating-desc":
          return ((b.imdbRating ?? 0) || (b.tmdbRating ?? 0)) - ((a.imdbRating ?? 0) || (a.tmdbRating ?? 0));
        case "year-desc":
          return (b.year ?? 0) - (a.year ?? 0);
        case "popularity-desc":
          return (b.popularity ?? 0) - (a.popularity ?? 0);
        default:
          return 0;
      }
    });

    return result;
  }, [media, search, genre, type, sort]);

  const hasActiveFilters = search !== "" || genre !== "all" || type !== "all" || sort !== "title-asc";

  const resetFilters = () => {
    setSearch("");
    setGenre("all");
    setType("all");
    setSort("title-asc");
  };

  return {
    search, setSearch,
    genre, setGenre,
    type, setType,
    sort, setSort: handleSortChange,
    allGenres,
    filtered,
    hasActiveFilters,
    resetFilters,
  };
}
