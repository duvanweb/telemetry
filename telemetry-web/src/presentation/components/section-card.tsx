import type { ReactNode } from "react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/presentation/ui/card";
import { Badge } from "@/presentation/ui/badge";

// SectionCard — wrapper Card with a title and optional "Demo" badge.
export function SectionCard({
  title,
  demo = false,
  children,
}: {
  title: string;
  demo?: boolean;
  children: ReactNode;
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          {title}
          {demo && <Badge variant="secondary">Demo</Badge>}
        </CardTitle>
      </CardHeader>
      <CardContent>{children}</CardContent>
    </Card>
  );
}
