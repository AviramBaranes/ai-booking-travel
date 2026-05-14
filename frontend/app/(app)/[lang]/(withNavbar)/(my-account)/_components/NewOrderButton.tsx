"use client";

import { Button } from "@/components/ui/button";
import { useState } from "react";
import { SearchForm } from "../../_components/home/SearchForm/SearchForm";
import { SearchDataFormWrapper } from "@/shared/components/booking/SearchDataFormWrapper";

export function NewOrderButton({btnText}: {btnText: string}) {
  const [showForm, setShowForm] = useState(false);

  return (
    <>
      {showForm ? (
        <div className="mb-4">
          <SearchDataFormWrapper onClose={() => setShowForm(false)}>
            <SearchForm className="w-full" />
          </SearchDataFormWrapper>
        </div>
      ) : (
        <div className="text-center my-13">
          <Button
            variant="brand"
            onClick={() => setShowForm(true)}
            className="py-6 w-50 mx-auto my-6 font-semibold"
          >
            {btnText}
          </Button>
        </div>
      )}
    </>
  );
}
