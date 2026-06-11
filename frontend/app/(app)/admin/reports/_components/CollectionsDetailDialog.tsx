"use client";

import { Dialog, DialogTitle, DialogContent } from "@/components/ui/dialog";
import { Suspense } from "react";
import { BillingResults } from "@/app/(app)/accounting/billing/_components/BillingResults";
import { BillingEntity } from "@/app/(app)/accounting/billing/_components/BillingEntityCombobox";
import { BillingShellSkeleton } from "@/app/(app)/accounting/billing/_components/BillingShell";

interface CollectionsDetailDialogProps {
  entity: BillingEntity | null;
  onClose: () => void;
}

export function CollectionsDetailDialog({
  entity,
  onClose,
}: CollectionsDetailDialogProps) {
  return (
    <Dialog open={entity != null} onOpenChange={(open) => !open && onClose()}>
      <DialogTitle>{entity?.name}</DialogTitle>
      <DialogContent
        className="max-h-10/12 w-10/12 max-w-10/12! overflow-y-auto m-auto"
        dir="rtl"
      >
        <Suspense
          key={`${entity?.kind}-${entity?.id}`}
          fallback={<BillingShellSkeleton />}
        >
          {entity && <BillingResults entity={entity} />}
        </Suspense>
      </DialogContent>
    </Dialog>
  );
}
