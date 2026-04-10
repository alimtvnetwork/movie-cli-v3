export interface MediaItem {
  id: number;
  title: string;
  cleanTitle: string;
  year: number;
  type: "movie" | "tv";
  tmdbId: number;
  imdbId: string;
  description: string;
  imdbRating: number;
  tmdbRating: number;
  popularity: number;
  genres: string[];
  director: string;
  cast: string[];
  thumbnailUrl: string;
  currentFilePath: string;
  fileSize: number;
  fileExtension: string;
  tags: string[];
  metadataStatus: "full" | "partial" | "filename-only";
}
