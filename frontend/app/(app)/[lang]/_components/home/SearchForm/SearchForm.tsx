"use client";

import { useState } from "react";
import { LocationCombobox } from "./LocationCombobox";
import { SearchFormOptions } from "./SearchFormOptions";
import { Button } from "@/components/ui/button";
import { CalendarInput } from "./CalendarInput";
import { TimeSelect } from "./TimeSelect";

export function SearchForm() {
  const [isReturnDifferentLoc, setIsReturnDifferentLoc] = useState(false);

  return (
    <form className="flex flex-col w-10/12 mx-auto mt-4">
      <SearchFormOptions
        isReturnDifferentLoc={isReturnDifferentLoc}
        setIsReturnDifferentLoc={setIsReturnDifferentLoc}
      />
      <div className="bg-white/95 w-full rounded-l-xl rounded-br-xl flex items-center gap-2 px-5">
        <div className="flex gap-2 flex-1 my-5 *:flex-1">
          <LocationCombobox placeholder="מהיכן תאספו את הרכב?" />
          {isReturnDifferentLoc && (
            <LocationCombobox placeholder="לאן תחזירו את הרכב?" />
          )}
        </div>
        <div className={isReturnDifferentLoc ? "w-1/10" : "w-1/6"}>
          <CalendarInput placeholder="מתי אוספים?" />
        </div>
        <div className={isReturnDifferentLoc ? "w-1/12" : "w-1/10"}>
          <TimeSelect placeholder="בחרו שעה" />
        </div>
        <div className={isReturnDifferentLoc ? "w-1/10" : "w-1/6"}>
          <CalendarInput placeholder="מתי מחזירים?" />
        </div>
        <div className={isReturnDifferentLoc ? "w-1/12" : "w-1/10"}>
          <TimeSelect placeholder="בחרו שעה" />
        </div>
        <div className="w-1/9">
          <Button
            variant="brand"
            className="w-full py-6.5 type-paragraph font-bold"
          >
            יאללה, בוא נזוז
          </Button>
        </div>
      </div>
    </form>
  );
}
