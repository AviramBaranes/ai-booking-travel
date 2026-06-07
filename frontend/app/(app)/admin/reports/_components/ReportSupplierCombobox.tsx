"use client";

import {
  Combobox,
  ComboboxInput,
  ComboboxContent,
  ComboboxList,
  ComboboxItem,
  ComboboxEmpty,
} from "@/components/ui/combobox";

interface ReportSupplierComboboxProps {
  value: string;
  onChange: (value: string) => void;
}

const supplierOptions = [
  { label: "ACA Rent", value: "CT,C5" },
  { label: "Ace", value: "AW" },
  { label: "Ada", value: "DA" },
  { label: "Add Car", value: "AR,AY" },
  { label: "Adobe", value: "AB,A2,A3" },
  { label: "Air Auto", value: "AI,A1" },
  { label: "Alamo", value: "AL,A4" },
  { label: "Alamo Max", value: "A6" },
  { label: "Arnold Clark", value: "CK" },
  { label: "AutAtlantis", value: "AA,AX" },
  { label: "Auto Menorca", value: "AU,A7,A8" },
  { label: "Auto Via", value: "AV" },
  { label: "Autos Union", value: "IO,I6" },
  { label: "AVIS", value: "VI,V7" },
  { label: "Avis Canaries", value: "AC,A5" },
  { label: "B-Rent", value: "B1,B2" },
  { label: "Bravacar Rent A Car", value: "BC" },
  { label: "Buchbinde", value: "BU" },
  { label: "Budget", value: "BG" },
  { label: "Budget Canaries", value: "BS,B3" },
  { label: "Budget Ireland", value: "BI" },
  { label: "Budget Romania", value: "BR" },
  { label: "Canarias.com", value: "CN,C1" },
  { label: "Car Rental Company", value: "CA,C3" },
  { label: "Centauro", value: "CE,CS,CP" },
  { label: "Cicar", value: "CC" },
  { label: "Cooltra", value: "CB" },
  { label: "Cuba on the Road", value: "CR" },
  { label: "Discount", value: "DI" },
  { label: "Dollar", value: "DS,D7" },
  { label: "Dollar Portugal", value: "D1,D2,D4" },
  { label: "Dollar Spain", value: "ZR,Z1,Z3" },
  { label: "Drivalia France", value: "VF,V4" },
  { label: "Drivalia Italy", value: "WR,W1,W2" },
  { label: "Drivalia Portugal", value: "VP,V5" },
  { label: "Drivalia Spain", value: "VS,V2" },
  { label: "Drivalia UK & IR", value: "ET,EP" },
  { label: "Drive A Matic", value: "DM" },
  { label: "Drive on Holidays", value: "DH,D6" },
  { label: "Eastcoast", value: "EC,E8,E7" },
  { label: "EasyCarHire", value: "EM,EE,EF" },
  { label: "Eco Via", value: "EV,ED" },
  { label: "eecar Balearics", value: "EJ,EH" },
  { label: "Enterprise", value: "EN,E9" },
  { label: "EuropCar", value: "ER,E5,E6" },
  { label: "Europcar Chanel Islands", value: "JR" },
  { label: "Europcar Cyprus", value: "EU" },
  { label: "Europcar Greece", value: "GE" },
  { label: "Europcar South Africa", value: "ES,E3,E2" },
  { label: "Ezi-rent", value: "EZ,E4,E1" },
  { label: "Felirent Italy", value: "FR,F5" },
  { label: "FireFly", value: "AD" },
  { label: "First Car Hire", value: "FC,F2,F1" },
  { label: "First car Hire Malta", value: "FH,F4" },
  { label: "Flizzr", value: "FZ,F3" },
  { label: "GoldCar", value: "GC,G1,G2" },
  { label: "Green motion Car Hire", value: "GM,G4,G5" },
  { label: "Guerin", value: "GU,G3,GP" },
  { label: "Helle Hollis", value: "HH,H9" },
  { label: "Hertz", value: "HT,H5,H6" },
  { label: "Hertz AutoHellas", value: "HG" },
  { label: "Hertz Balearics", value: "HZ,H1,Z2" },
  { label: "Hertz Canaries", value: "HC,H3,C2" },
  { label: "Hertz Cape Verde", value: "HV" },
  { label: "Hertz Panama", value: "HN" },
  { label: "Hertz Portugal", value: "HP,H7,H8" },
  { label: "Hiper", value: "HE,H4,H2" },
  { label: "Ilha Verde Rent A Car", value: "IV,I4" },
  { label: "Island Car Rentals", value: "IS" },
  { label: "Italy Car Rent", value: "IC,I3" },
  { label: "Keddy", value: "KD,K1,K2" },
  { label: "LocAuto", value: "LA,L1,L2" },
  { label: "Madeira Rent", value: "MR,M1" },
  { label: "Maggiore Italy", value: "GI,G6" },
  { label: "Master Kings", value: "MK" },
  { label: "Megadrive", value: "MG" },
  { label: "Micauto", value: "ZM,Z5" },
  { label: "Movida", value: "MV" },
  { label: "National", value: "ZL,Z4" },
  { label: "Nokta", value: "NT" },
  { label: "Noleggiare", value: "NG,N1,N2" },
  { label: "OK Rent a Car", value: "OK,O1,O2" },
  { label: "Orlando Rent A Car", value: "OR,OS,OP" },
  { label: "Owners Car Rental", value: "OW,O4,O5" },
  { label: "Record", value: "RE,R4,R5" },
  { label: "Right Car Vehicle Rental", value: "RC,R3" },
  { label: "Routes", value: "RZ,R6" },
  { label: "SIXT", value: "SX,S2,X1" },
  { label: "Sicily By Car", value: "SC,S1,SZ" },
  { label: "Stoutes Car Rental", value: "ST" },
  { label: "SurPrice", value: "SP,S3" },
  { label: "Tempest Car Hire", value: "TE,T1,T2" },
  { label: "Thrifty", value: "TY,T9" },
  { label: "Thrifty Iceland", value: "TI" },
  { label: "Top Car Canaries", value: "TX,T8" },
  { label: "Tourent", value: "TO,T6,T7" },
  { label: "Viaggiare Rent", value: "VR,V3" },
  { label: "Volta Greece", value: "VG,V1,V6" },
  { label: "WIber", value: "WI,W3,W4" },
  { label: "Wheego", value: "W5" },
  { label: "Wheego Premium", value: "WG" }
];

export function ReportSupplierCombobox({
  value,
  onChange,
}: ReportSupplierComboboxProps) {
  const selectedOption = supplierOptions.find(opt => opt.value === value);

  return (
    <Combobox
      items={supplierOptions.map(opt => opt.label)}
      value={selectedOption ? selectedOption.label : (value || null)}
      onValueChange={(nextLabel: string | null) => {
        const option = supplierOptions.find(opt => opt.label === nextLabel);
        onChange(option ? option.value : (nextLabel ?? ""));
      }}
    >
      <ComboboxInput
        placeholder="כל הספקים"
        showClear
        className="w-full"
        value={selectedOption ? selectedOption.label : value}
        onChange={(event) => onChange(event.target.value)}
      />
      <ComboboxContent>
        <ComboboxEmpty>אין ספקים להצגה כרגע</ComboboxEmpty>
        <ComboboxList>
          {(supplierLabel) => (
            <ComboboxItem key={supplierLabel} value={supplierLabel}>
              {supplierLabel}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
