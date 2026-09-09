import { useState } from "react";
import type { FormEvent } from "react";
import { useSearchParams } from "react-router-dom";
import { Button } from "@/presentation/ui/button";
import { Input } from "@/presentation/ui/input";

const PLATE_REGEX = /^[A-Z]{3}-[0-9]{3}$/;

// VehicleSearch — input + button to search a vehicle by plate.
// Validates the plate format client-side before updating the URL search params.
export function VehicleSearch() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [value, setValue] = useState(searchParams.get("q") ?? "");
  const [error, setError] = useState<string | null>(null);

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    const trimmed = value.trim();

    if (trimmed === "") {
      setError(null);
      const next = new URLSearchParams(searchParams);
      next.delete("q");
      next.set("page", "1");
      setSearchParams(next);
      return;
    }

    if (!PLATE_REGEX.test(trimmed)) {
      setError("Invalid plate format, expected ABC-123");
      return;
    }

    setError(null);
    const next = new URLSearchParams(searchParams);
    next.set("q", trimmed);
    next.set("page", "1");
    setSearchParams(next);
  }

  return (
    <form onSubmit={handleSubmit} className="flex items-center gap-2">
      <Input
        type="text"
        placeholder="Search by plate (ABC-123)"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        className="max-w-xs"
      />
      <Button type="submit">Search</Button>
      {error && (
        <span className="text-sm text-destructive">{error}</span>
      )}
    </form>
  );
}
