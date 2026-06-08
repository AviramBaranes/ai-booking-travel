"use client";

import { Button } from "@/components/ui/button";
import { approvePriceOffer } from "@/shared/api/price-offers-api";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { toast } from "sonner";

interface ApproveButtonProps {
  id: number;
  text: string;
  successText: string;
}
export function ApproveButton({ id, text, successText }: ApproveButtonProps) {
  const [disable, setDisable] = useState(false);
  const { isPending, mutate } = useMutation({
    mutationFn: () => approvePriceOffer(id),
    onSuccess: () => {
      setDisable(true);
      toast.success(successText, {
        duration: 5000,
        position: "top-center",
      });
      setTimeout(() => {
        window.location.reload();
      }, 5000);
    },
  });

  return (
    <Button
      type="button"
      variant="brand"
      loading={isPending}
      disabled={disable}
      onClick={() => mutate()}
      className="w-full py-6 type-paragraph font-bold"
    >
      {text}
    </Button>
  );
}
