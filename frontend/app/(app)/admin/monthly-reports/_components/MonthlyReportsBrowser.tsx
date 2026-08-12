"use client";

import { useMemo, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { ChevronRight, Download, Folder, FileSpreadsheet } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  getMonthlyReportUrl,
  listMonthlyReports,
} from "@/shared/api/monthly-reports-api";
import type { billing } from "@/shared/client";

type ReportFolder = billing.MonthlyReportFolder;

const folderKey = (folder: ReportFolder) =>
  `${folder.entityType}/${folder.entityId}`;

export function MonthlyReportsBrowser() {
  const [openFolderKey, setOpenFolderKey] = useState<string | null>(null);

  const foldersQuery = useQuery({
    queryKey: ["monthly-reports"],
    queryFn: listMonthlyReports,
  });

  const folders = useMemo(
    () => foldersQuery.data?.folders ?? [],
    [foldersQuery.data],
  );
  const openFolder =
    folders.find((folder) => folderKey(folder) === openFolderKey) ?? null;

  const download = useMutation({
    mutationFn: async ({
      folder,
      period,
    }: {
      folder: ReportFolder;
      period: string;
    }) => {
      const { url } = await getMonthlyReportUrl(
        folder.entityType,
        folder.entityId,
        period,
      );
      return url;
    },
    onSuccess: (url) => {
      // The signed link is short lived, so it is opened the moment it arrives.
      window.location.assign(url);
    },
  });

  if (foldersQuery.isPending) {
    return <EmptyState title="טוען דוחות..." />;
  }

  if (foldersQuery.isError) {
    return (
      <EmptyState
        title="לא הצלחנו לטעון את הדוחות"
        description="נסה לרענן את העמוד."
      />
    );
  }

  if (folders.length === 0) {
    return (
      <EmptyState
        title="אין עדיין דוחות"
        description="דוחות נשמרים אוטומטית בכל הרצה של החיוב החודשי."
      />
    );
  }

  if (!openFolder) {
    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {folders.map((folder) => (
          <button
            key={folderKey(folder)}
            type="button"
            onClick={() => setOpenFolderKey(folderKey(folder))}
            className="bg-card shadow-card hover:shadow-card-hover flex cursor-pointer items-center gap-3 rounded-2xl p-5 text-right transition"
          >
            <Folder className="text-navy size-8 shrink-0" />
            <span className="flex flex-col">
              <span className="type-h6 text-navy">{folder.name}</span>
              <span className="type-paragraph text-text-secondary">
                {folder.files.length} דוחות
              </span>
            </span>
          </button>
        ))}
      </div>
    );
  }

  return (
    <div className="bg-card shadow-card flex flex-col gap-4 rounded-2xl p-6">
      <div className="text-text-secondary type-paragraph flex items-center gap-1">
        <button
          type="button"
          onClick={() => setOpenFolderKey(null)}
          className="hover:text-navy cursor-pointer underline-offset-4 hover:underline"
        >
          כל התיקיות
        </button>
        <ChevronRight className="size-4 rotate-180" />
        <span className="text-navy">{openFolder.name}</span>
      </div>

      <ul className="flex flex-col divide-y">
        {openFolder.files.map((file) => (
          <li
            key={file.period}
            className="flex items-center justify-between gap-4 py-3"
          >
            <span className="flex items-center gap-3">
              <FileSpreadsheet className="text-navy size-6 shrink-0" />
              <span className="flex flex-col">
                <span className="type-paragraph text-navy">
                  {formatPeriod(file.period)}
                </span>
                <span className="type-paragraph text-text-secondary">
                  {formatFileSize(file.size)}
                </span>
              </span>
            </span>

            <Button
              variant="outline"
              onClick={() =>
                download.mutate({ folder: openFolder, period: file.period })
              }
              disabled={download.isPending}
            >
              <Download className="size-4" />
              הורדה
            </Button>
          </li>
        ))}
      </ul>

      {download.isError && (
        <p className="type-paragraph text-destructive">
          לא הצלחנו להוריד את הדוח. נסה שוב.
        </p>
      )}
    </div>
  );
}

function EmptyState({
  title,
  description,
}: {
  title: string;
  description?: string;
}) {
  return (
    <div className="bg-card shadow-card flex flex-col items-center gap-2 rounded-2xl p-12 text-center">
      <h3 className="type-h6 text-navy">{title}</h3>
      {description && (
        <p className="type-paragraph text-text-secondary">{description}</p>
      )}
    </div>
  );
}

// formatPeriod turns the YYYY-MM the report is filed under into a readable Hebrew month.
function formatPeriod(period: string) {
  const [year, month] = period.split("-").map(Number);
  if (!year || !month) return period;

  return new Date(year, month - 1).toLocaleDateString("he-IL", {
    month: "long",
    year: "numeric",
  });
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${Math.round(bytes / 1024)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
