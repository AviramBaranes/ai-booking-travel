import LocationMerge from "./_components/LocationMerge";

export default function LocationMergePage() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-700">מיזוג מיקומים</h1>
      <p className="text-sm text-gray-500">
        בחרו שני מיקומים מהחיפוש או מההצעות, קבעו מאיזה מיקום לקחת כל שדה,
        ושמרו. המיקום השני יימחק, וקודי הספקים והשמות החלופיים שלו יעברו
        למיקום הממוזג.
      </p>
      <LocationMerge />
    </div>
  );
}
