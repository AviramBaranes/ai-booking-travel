package booking_handlers

import (
	"encoding/json"

	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	availability "encore.app/services/booking/handlers/availability"
	"encore.dev/rlog"
)

// FindSupplierInfo locates the supplier info of the given plan inside the snapshot's suppliers_info blob,
// which holds the terms and the pickup/dropoff station details returned by the broker availability call.
//
// A missing or malformed entry is not fatal: this data is informational and is resolved after the car has
// already been booked at the broker, so it must never abort a reservation. The zero value is returned instead.
func FindSupplierInfo(snapshot db.AvailablePlansSnapshot, plan availability.PlanPriceDetails) broker.SupplierInfo {
	var suppliers []broker.SupplierInfo
	if err := json.Unmarshal(snapshot.SuppliersInfo, &suppliers); err != nil {
		rlog.Error("failed to unmarshal suppliers info JSON", "snapshotID", snapshot.ID, "error", err)
		return broker.SupplierInfo{}
	}

	for _, supplier := range suppliers {
		if supplier.Name == plan.SupplierName {
			return supplier
		}
	}

	rlog.Warn("no supplier info found for plan", "snapshotID", snapshot.ID, "supplierName", plan.SupplierName)
	return broker.SupplierInfo{}
}
