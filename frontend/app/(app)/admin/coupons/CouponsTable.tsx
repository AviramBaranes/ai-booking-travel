"use client";

import { z } from "zod";
import { booking } from "@/shared/client";
import { CrudTable } from "@/app/(app)/admin/components/crud-table/CrudTable";
import { ColumnDef } from "@/app/(app)/admin/components/crud-table/types";
import {
  listCoupons,
  createCoupon,
  updateCoupon,
  deleteCoupon,
} from "@/shared/api/coupons-api";

const columns: ColumnDef<booking.CouponResponse>[] = [
  { key: "id", label: "מזהה", type: "number", editable: false },
  { key: "name", label: "שם", type: "text" },
  { key: "code", label: "קוד", type: "text" },
  { key: "discount", label: "הנחה %", type: "number" },
  { key: "isEnabled", label: "פעיל", type: "checkbox" },
];

const couponSchema = z.object({
  name: z.string().min(1, "שדה חובה"),
  code: z.string().min(1, "שדה חובה"),
  discount: z
    .number({ error: "מספר נדרש" })
    .min(0, "ערך מינימלי 0")
    .max(100, "ערך מקסימלי 100"),
  isEnabled: z.boolean(),
});

export default function CouponsTable() {
  return (
    <CrudTable<
      booking.CouponResponse,
      booking.CreateCouponRequest,
      booking.UpdateCouponRequest
    >
      columns={columns}
      queryKey="coupons"
      getId={(r) => r.id}
      listFn={() => listCoupons()}
      extractList={(r) =>
        (r as booking.ListCouponsResponse | undefined)?.coupons ?? []
      }
      createFn={createCoupon}
      updateFn={updateCoupon}
      deleteFn={deleteCoupon}
      createSchema={couponSchema}
      updateSchema={couponSchema}
    />
  );
}
