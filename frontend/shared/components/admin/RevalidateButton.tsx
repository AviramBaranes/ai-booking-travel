"use client";

import { useState } from "react";
import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { revalidateHomepage } from "@/shared/actions/revalidate";
import clsx from "clsx";

export function RevalidateButton() {
  const [loading, setLoading] = useState(false);

  const handleClick = async () => {
    setLoading(true);
    try {
      const result = await revalidateHomepage();
      if (!result.ok) {
        console.error("Revalidation failed:", result.error);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={handleClick}
      disabled={loading}
      className="flex items-center gap-1.5 text-[13px]"
    >
      <RefreshCw
        size={14}
        className={clsx(loading && "animate-spin")}
      />
      ניקוי מטמון אתר
    </Button>
  );
}
