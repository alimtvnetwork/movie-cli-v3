import { useState } from "react";
import { Star, Film, Tv, HardDrive, FileText, Tag, Copy } from "lucide-react";
import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import { Separator } from "@/components/ui/separator";
import type { MediaItem } from "@/types/media";
import { ratingColorClass, formatFileSize } from "@/lib/media-utils";
import { PosterFallback } from "./MediaCard";

interface MediaDetailModalProps {
  media: MediaItem | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const STATUS_MAP = {
  full: { label: "Full", className: "bg-status-full/15 text-status-full" },
  partial: { label: "Partial", className: "bg-status-partial/15 text-status-partial" },
  "filename-only": { label: "Filename Only", className: "bg-status-filename/15 text-status-filename" },
} as const;

function MetadataStatusBadge({ status }: { status: MediaItem["metadataStatus"] }) {
  const { label, className } = STATUS_MAP[status];
  return <Badge variant="secondary" className={className}>{label}</Badge>;
}

export function MediaDetailModal({ media, open, onOpenChange }: MediaDetailModalProps) {
  const [imgError, setImgError] = useState(false);

  // Reset img error state when media changes
  if (!media) return null;

  return (
    <Dialog open={open} onOpenChange={(o) => { if (!o) setImgError(false); onOpenChange(o); }}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl">
            {media.type === "movie"
              ? <Film className="h-5 w-5 text-badge-movie" />
              : <Tv className="h-5 w-5 text-badge-tv" />}
            {media.title}
            <span className="text-sm font-normal text-muted-foreground">({media.year})</span>
          </DialogTitle>
          <DialogDescription className="sr-only">Details for {media.title}</DialogDescription>
        </DialogHeader>

        <div className="grid gap-4 sm:grid-cols-[180px_1fr]">
          <div className="mx-auto w-[180px]">
            <AspectRatio ratio={2 / 3}>
              {media.thumbnailUrl && !imgError ? (
                <img
                  src={media.thumbnailUrl}
                  alt={media.title}
                  className="h-full w-full rounded-md object-cover"
                  onError={() => setImgError(true)}
                />
              ) : (
                <PosterFallback title={media.title} />
              )}
            </AspectRatio>
          </div>

          <div className="space-y-3">
            <p className="text-sm text-muted-foreground">{media.description}</p>

            <Separator />

            <div className="flex flex-wrap gap-4 text-sm">
              {media.imdbRating > 0 && (
                <span className={`flex items-center gap-1 font-medium ${ratingColorClass(media.imdbRating)}`}>
                  <Star className="h-3.5 w-3.5 fill-current" /> IMDb {media.imdbRating.toFixed(1)}
                </span>
              )}
              {media.tmdbRating > 0 && (
                <span className={`flex items-center gap-1 font-medium ${ratingColorClass(media.tmdbRating)}`}>
                  <Star className="h-3.5 w-3.5 fill-current" /> TMDb {media.tmdbRating.toFixed(1)}
                </span>
              )}
              <span className="text-muted-foreground">Popularity: {media.popularity.toFixed(0)}</span>
            </div>

            {(media.genres?.length ?? 0) > 0 && (
              <div className="flex flex-wrap gap-1.5">
                {media.genres.map((g) => (
                  <Badge key={g} variant="outline" className="text-xs">{g}</Badge>
                ))}
              </div>
            )}

            <div className="space-y-1 text-sm">
              <p><span className="font-medium text-foreground">Director:</span> <span className="text-muted-foreground">{media.director}</span></p>
              <p><span className="font-medium text-foreground">Cast:</span> <span className="text-muted-foreground">{(media.cast ?? []).join(", ")}</span></p>
            </div>

            <Separator />

            <div className="space-y-1 text-xs text-muted-foreground">
              <p className="flex items-center gap-1.5"><HardDrive className="h-3.5 w-3.5" /> {formatFileSize(media.fileSize)} &middot; {media.fileExtension}</p>
              <button
                type="button"
                className="flex items-center gap-1.5 hover:text-foreground transition-colors group"
                onClick={() => {
                  navigator.clipboard.writeText(media.currentFilePath).then(
                    () => toast.success("Copied to clipboard"),
                    () => toast.error("Failed to copy path")
                  );
                }}
                title="Click to copy file path"
              >
                <FileText className="h-3.5 w-3.5" />
                <span className="truncate max-w-[280px]">{media.currentFilePath}</span>
                <Copy className="h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity" />
              </button>
              <div className="flex items-center gap-1.5">
                <Tag className="h-3.5 w-3.5" />
                <MetadataStatusBadge status={media.metadataStatus} />
                {(media.tags ?? []).length > 0 && media.tags.map((t) => (
                  <Badge key={t} variant="secondary" className="text-[10px] px-1.5 py-0">{t}</Badge>
                ))}
              </div>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
