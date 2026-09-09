import { Button } from "@/presentation/ui/button";

// PaginationControls — prev/next buttons with page indicator and total count.
export function PaginationControls({
  page,
  totalPages,
  total,
  onChangePage,
}: {
  page: number;
  totalPages: number;
  total: number;
  onChangePage: (page: number) => void;
}) {
  return (
    <div className="flex items-center justify-between">
      <p className="text-sm text-muted-foreground">
        Page {page} of {totalPages} · {total} vehicles
      </p>
      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          disabled={page <= 1}
          onClick={() => onChangePage(page - 1)}
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          disabled={page >= totalPages}
          onClick={() => onChangePage(page + 1)}
        >
          Next
        </Button>
      </div>
    </div>
  );
}
