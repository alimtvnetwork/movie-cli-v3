/**
 * Shared media utility functions.
 * Consolidates formatting, rating helpers, and constants used across dashboard components.
 */

/** Returns a semantic Tailwind text color class for a numeric rating. */
export function ratingColorClass(rating: number): string {
  if (rating >= 7) return "text-rating-high";
  if (rating >= 5) return "text-rating-mid";
  return "text-rating-low";
}

/** Human-readable file size string. Uses SI units (1 GB = 1,000,000,000 bytes). */
export function formatFileSize(bytes: number): string {
  if (bytes >= 1e9) return `${(bytes / 1e9).toFixed(1)} GB`;
  if (bytes >= 1e6) return `${(bytes / 1e6).toFixed(0)} MB`;
  if (bytes >= 1e3) return `${(bytes / 1e3).toFixed(0)} KB`;
  return `${bytes} B`;
}

/** Safely reads a value from localStorage, returning the fallback on any error. */
export function safeLocalGet(key: string, fallback: string): string {
  try {
    return localStorage.getItem(key) ?? fallback;
  } catch {
    return fallback;
  }
}

/** Safely writes a value to localStorage, silently ignoring errors. */
export function safeLocalSet(key: string, value: string): void {
  try {
    localStorage.setItem(key, value);
  } catch {
    // Ignore — private browsing or quota exceeded
  }
}
