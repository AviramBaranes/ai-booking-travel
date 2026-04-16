import { create } from "zustand";
import { broker } from "@/shared/client";

interface BookingSessionState {
  selectedPlanIndex: number;
  isErpSelected: boolean;
  selectedAddons: broker.SelectAddOn[];
  setSelectedPlanIndex: (index: number) => void;
  setIsErpSelected: (value: boolean) => void;
  setSelectedAddons: (addons: broker.SelectAddOn[]) => void;
  clearSession: () => void;
}

export const useBookingSessionStore = create<BookingSessionState>((set) => ({
  selectedPlanIndex: 0,
  isErpSelected: false,
  selectedAddons: [],
  setSelectedPlanIndex: (index) => set({ selectedPlanIndex: index }),
  setIsErpSelected: (value) => set({ isErpSelected: value }),
  setSelectedAddons: (addons) => set({ selectedAddons: addons }),
  clearSession: () =>
    set({ selectedPlanIndex: 0, isErpSelected: false, selectedAddons: [] }),
}));
