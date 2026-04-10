import { useState } from "react";
import { MediaCard } from "./MediaCard";
import { MediaDetailModal } from "./MediaDetailModal";
import type { MediaItem } from "@/types/media";
import { toast } from "sonner";

interface MediaGridProps {
  media: MediaItem[];
}

export function MediaGrid({ media }: MediaGridProps) {
  const [selected, setSelected] = useState<MediaItem | null>(null);
  if (media.length === 0) {
    return (
      <section className="rounded-lg border border-border bg-card p-8 text-center text-muted-foreground">
        <p className="text-lg font-medium">No media found</p>
        <p className="text-sm mt-1">
          Run <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">mahin movie scan</code> to get started.
        </p>
      </section>
    );
  }

  const handleSelect = (item: MediaItem) => {
    setSelected(item);
    toast.info(`Viewing details for ${item.title}`);
  };

  const handleClose = (open: boolean) => {
    if (!open) {
      toast(`Closed ${selected?.title ?? "details"}`);
      setSelected(null);
    }
  };

  return (
    <>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        {media.map((item) => (
          <MediaCard key={item.id} media={item} onClick={() => handleSelect(item)} />
        ))}
      </div>
      <MediaDetailModal media={selected} open={!!selected} onOpenChange={handleClose} />
    </>
  );
}
