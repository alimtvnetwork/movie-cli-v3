import type { MediaItem } from "@/types/media";

interface DashboardHeaderProps {
  media: MediaItem[];
}

export function DashboardHeader({ media }: DashboardHeaderProps) {
  return (
    <header>
      <h1 className="text-3xl font-bold tracking-tight text-foreground">
        Media Library
      </h1>
      <p className="text-muted-foreground mt-1">
        {media.length} titles in your collection
      </p>
    </header>
  );
}
