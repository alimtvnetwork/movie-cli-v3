import { useState } from "react";
import { Star } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import type { MediaItem } from "@/types/media";
import { ratingColorClass } from "@/lib/media-utils";

interface MediaCardProps {
  media: MediaItem;
  onClick?: () => void;
}

function PosterFallback({ title }: { title: string }) {
  return (
    <div className="flex h-full w-full items-center justify-center bg-gradient-to-br from-muted to-muted-foreground/20">
      <span className="text-4xl font-bold text-muted-foreground/60">
        {title.charAt(0).toUpperCase()}
      </span>
    </div>
  );
}

export { PosterFallback };

export function MediaCard({ media, onClick }: MediaCardProps) {
  const [imgError, setImgError] = useState(false);
  const rating = media.imdbRating || media.tmdbRating;
  const displayGenres = (media.genres ?? []).slice(0, 3);
  const displayCast = (media.cast ?? []).slice(0, 3);

  return (
    <Card className="cursor-pointer overflow-hidden transition-shadow hover:shadow-lg" onClick={onClick}>
      <AspectRatio ratio={2 / 3}>
        {media.thumbnailUrl && !imgError ? (
          <img
            src={media.thumbnailUrl}
            alt={media.title}
            className="h-full w-full object-cover"
            loading="lazy"
            onError={() => setImgError(true)}
          />
        ) : (
          <PosterFallback title={media.title} />
        )}
      </AspectRatio>

      <div className="space-y-2 p-4">
        <div className="flex items-start justify-between gap-2">
          <h3 className="line-clamp-1 text-sm font-semibold text-card-foreground">
            {media.title}
          </h3>
          <Badge
            variant="secondary"
            className={
              media.type === "movie"
                ? "shrink-0 bg-badge-movie/15 text-badge-movie"
                : "shrink-0 bg-badge-tv/15 text-badge-tv"
            }
          >
            {media.type === "movie" ? "Movie" : "TV"}
          </Badge>
        </div>

        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <span>{media.year}</span>
          {rating > 0 && (
            <>
              <span>•</span>
              <span className={`flex items-center gap-0.5 font-medium ${ratingColorClass(rating)}`}>
                <Star className="h-3 w-3 fill-current" />
                {rating.toFixed(1)}
              </span>
            </>
          )}
        </div>

        {displayGenres.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {displayGenres.map((genre) => (
              <Badge key={genre} variant="outline" className="text-[10px] px-1.5 py-0">
                {genre}
              </Badge>
            ))}
          </div>
        )}

        {displayCast.length > 0 && (
          <p className="line-clamp-1 text-xs text-muted-foreground">
            {displayCast.join(", ")}
          </p>
        )}
      </div>
    </Card>
  );
}
