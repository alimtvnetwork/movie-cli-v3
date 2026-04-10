import { useMemo } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import { BarChart, Bar, XAxis, YAxis, PieChart, Pie, Cell } from "recharts";
import { Film, Tv, Star, HardDrive } from "lucide-react";
import type { MediaItem } from "@/types/media";
import { formatFileSize } from "@/lib/media-utils";

interface StatsPanelProps {
  media: MediaItem[];
}

const PIE_COLORS = [
  "hsl(var(--chart-1))",
  "hsl(var(--chart-2))",
  "hsl(var(--chart-3))",
  "hsl(var(--chart-4))",
  "hsl(var(--chart-5))",
  "hsl(var(--accent))",
  "hsl(var(--muted-foreground))",
];

export function StatsPanel({ media }: StatsPanelProps) {
  const stats = useMemo(() => {
    const movies = media.filter((m) => m.type === "movie").length;
    const tv = media.filter((m) => m.type === "tv").length;
    const avgImdb =
      media.length > 0
        ? media.reduce((s, m) => s + (m.imdbRating ?? 0), 0) / media.length
        : 0;
    const avgTmdb =
      media.length > 0
        ? media.reduce((s, m) => s + (m.tmdbRating ?? 0), 0) / media.length
        : 0;
    const totalSize = media.reduce((s, m) => s + (m.fileSize ?? 0), 0);

    const genreMap = new Map<string, number>();
    media.forEach((m) => (m.genres ?? []).forEach((g) => genreMap.set(g, (genreMap.get(g) || 0) + 1)));
    const genreData = Array.from(genreMap.entries())
      .sort((a, b) => b[1] - a[1])
      .map(([name, count]) => ({ name, count }));

    const ratingBuckets = [
      { range: "0-4", count: 0 },
      { range: "4-6", count: 0 },
      { range: "6-7", count: 0 },
      { range: "7-8", count: 0 },
      { range: "8+", count: 0 },
    ];
    media.forEach((m) => {
      const r = m.imdbRating ?? 0;
      if (r < 4) ratingBuckets[0].count++;
      else if (r < 6) ratingBuckets[1].count++;
      else if (r < 7) ratingBuckets[2].count++;
      else if (r < 8) ratingBuckets[3].count++;
      else ratingBuckets[4].count++;
    });

    return { movies, tv, avgImdb, avgTmdb, totalSize, genreData, ratingBuckets };
  }, [media]);

  const barConfig: ChartConfig = {
    count: { label: "Titles", color: "hsl(var(--primary))" },
  };

  const pieConfig: ChartConfig = Object.fromEntries(
    stats.genreData.map((g, i) => [
      g.name,
      { label: g.name, color: PIE_COLORS[i % PIE_COLORS.length] },
    ])
  );

  if (media.length === 0) {
    return (
      <p className="text-sm text-muted-foreground text-center py-4">
        No data to display. Add media to see statistics.
      </p>
    );
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <SummaryCard icon={Film} label="Movies" value={stats.movies} />
        <SummaryCard icon={Tv} label="TV Shows" value={stats.tv} />
        <SummaryCard icon={Star} label="Avg IMDb" value={stats.avgImdb.toFixed(1)} />
        <SummaryCard icon={HardDrive} label="Total Size" value={formatFileSize(stats.totalSize)} />
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-base">Genre Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer config={pieConfig} className="aspect-square max-h-[260px] w-full">
              <PieChart>
                <ChartTooltip content={<ChartTooltipContent nameKey="name" />} />
                <Pie
                  data={stats.genreData}
                  dataKey="count"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  innerRadius={50}
                  outerRadius={90}
                  paddingAngle={2}
                >
                  {stats.genreData.map((g, i) => (
                    <Cell key={g.name} fill={PIE_COLORS[i % PIE_COLORS.length]} />
                  ))}
                </Pie>
              </PieChart>
            </ChartContainer>
            <div className="flex flex-wrap gap-2 mt-3">
              {stats.genreData.map((g, i) => (
                <span key={g.name} className="flex items-center gap-1.5 text-xs text-muted-foreground">
                  <span
                    className="inline-block h-2 w-2 rounded-full"
                    style={{ backgroundColor: PIE_COLORS[i % PIE_COLORS.length] }}
                  />
                  {g.name}
                </span>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-base">Rating Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer config={barConfig} className="aspect-video max-h-[260px] w-full">
              <BarChart data={stats.ratingBuckets}>
                <XAxis dataKey="range" tickLine={false} axisLine={false} />
                <YAxis allowDecimals={false} tickLine={false} axisLine={false} />
                <ChartTooltip content={<ChartTooltipContent />} />
                <Bar dataKey="count" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ChartContainer>
            <p className="text-xs text-muted-foreground mt-2 text-center">
              Avg TMDb: {stats.avgTmdb.toFixed(1)} · Avg IMDb: {stats.avgImdb.toFixed(1)}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function SummaryCard({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ElementType;
  label: string;
  value: string | number;
}) {
  return (
    <Card>
      <CardContent className="p-4 flex items-center gap-3">
        <div className="rounded-md bg-primary/10 p-2">
          <Icon className="h-5 w-5 text-primary" />
        </div>
        <div>
          <p className="text-2xl font-bold leading-none">{value}</p>
          <p className="text-xs text-muted-foreground mt-1">{label}</p>
        </div>
      </CardContent>
    </Card>
  );
}
