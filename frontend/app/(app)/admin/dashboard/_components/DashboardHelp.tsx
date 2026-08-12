"use client";

import { ReactNode } from "react";
import { HelpCircle } from "lucide-react";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

export function DashboardHelp() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <button
          type="button"
          data-noprint
          className="flex items-center gap-1.5 rounded-md border border-gray-200 bg-white px-3 py-1.5 text-xs text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900"
        >
          <HelpCircle className="size-4" />
          איך עובד לוח הבקרה?
        </button>
      </DialogTrigger>

      <DialogContent
        dir="rtl"
        className="max-h-[85vh] overflow-y-auto text-right sm:max-w-2xl"
      >
        <DialogHeader>
          <DialogTitle className="text-right">איך עובד לוח הבקרה</DialogTitle>
          <DialogDescription className="text-right">
            מדריך קצר לקריאת הנתונים ולשימוש בפילטרים.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-5 text-sm leading-relaxed text-gray-700">
          <Section title="על מה מבוססים המספרים">
            <ul className="list-inside list-disc space-y-1">
              <li>
                כל הנתונים בלוח מבוססים על הזמנות ש<strong>נוצרו</strong> בטווח
                התאריכים שנבחר — לא על תאריך האיסוף ולא על תאריך הכרטוס.
              </li>
              <li>
                כל הסכומים מוצגים בשקלים. הזמנות במטבע אחר מומרות לפי שער ההמרה
                שנשמר על ההזמנה עצמה.
              </li>
              <li>
                חישובי העלות והרווח נעשים בשרת, כך שכל מסך שמציג רווח מסתמך על
                אותה נוסחה בדיוק.
              </li>
            </ul>
          </Section>

          <Section title="ההגדרות הכספיות">
            <dl className="space-y-2">
              <Definition term="עלות">
                מה שאנחנו חייבים לספק: מחיר הרכב בתוספת ה-ERP של הברוקר.
              </Definition>
              <Definition term="הכנסות">
                המחיר שהלקוח או הסוכן משלם בפועל, כולל ה-ERP שאנחנו מוכרים
                ולאחר הנחת קופון.
              </Definition>
              <Definition term="רווח">
                הכנסות פחות עלות.
              </Definition>
            </dl>
            <p className="mt-2 rounded-md border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-900">
              שימו לב: הרווח כאן מנוכה הנחות. בהזמנות עם קופון הוא יהיה נמוך מזה
              שמוצג בשורות של &quot;דוח רווחיות&quot;, שמתעלם מאחוז ההנחה.
              הסיכומים הכוללים של שני המסכים כן מתיישבים זה עם זה.
            </p>
          </Section>

          <Section title="הפילטרים">
            <p className="mb-2">
              שורת הפילטרים בראש העמוד משפיעה על <strong>כל</strong> הכרטיסים
              והגרפים יחד, כך שכל המסך תמיד מתאר את אותה קבוצת הזמנות.
            </p>
            <dl className="space-y-2">
              <Definition term="תאריכים">
                קיצורי דרך מוכנים (היום, 7 ימים אחרונים, החודש הקודם, מתחילת
                השנה ועוד) או בחירה ידנית של טווח בלוח השנה. ברירת המחדל היא
                השבוע הנוכחי.
              </Definition>
              <Definition term="קהל">
                הפרדה בין הזמנות עסקיות (סוכן ששייך למשרד וארגון) לבין הזמנות
                פרטיות.
              </Definition>
              <Definition term="ביטולים">
                האם לכלול הזמנות מבוטלות בחישובים. כברירת מחדל הן אינן נכללות.
              </Definition>
              <Definition term="פילוחים">
                קובע מה גרפי הפילוח מודדים — כמות הזמנות, הכנסות או רווח.
              </Definition>
            </dl>
            <p className="mt-2 text-xs text-gray-500">
              הבחירות נשמרות בכתובת העמוד, כך שאפשר לשלוח קישור לתצוגה מסוימת.
            </p>
          </Section>

          <Section title="ההשוואה לתקופה הקודמת">
            <p>
              החיצים ▲▼ שמופיעים בכרטיסי המספרים משווים את הטווח שנבחר לתקופה
              באותו אורך שקדמה לו מיד.
            </p>
            <ul className="mt-2 list-inside list-disc space-y-1">
              <li>
                לדוגמה: בחירה של 1–7 באוגוסט (7 ימים) תושווה ל-25–31 ביולי.
              </li>
              <li>
                אחוז השינוי מחושב כך: (נוכחי פחות קודם) חלקי הערך הקודם, כפול
                מאה. אם בתקופה הקודמת לא היו נתונים כלל, לא יוצג אחוז — עלייה
                מאפס אינה אחוז בעל משמעות.
              </li>
              <li>
                ב&quot;עלות&quot; וב&quot;שיעור ביטולים&quot; ירידה היא התוצאה
                הרצויה, ולכן הצבע מתהפך. כיוון החץ תמיד מתאר את השינוי בפועל.
              </li>
              <li>
                ההשוואה מבוססת על מספר ימים ולא על לוח השנה: חודש בן 28 יום
                יושווה ל-28 הימים שקדמו לו, ולא לחודש הקודם המלא.
              </li>
            </ul>
          </Section>

          <Section title="קריאת הגרפים">
            <ul className="list-inside list-disc space-y-1">
              <li>
                בכל כרטיס גרף יש כפתור <strong>טבלה</strong> שמציג את אותם
                הנתונים כמספרים — שימושי להעתקה ולקריאה מדויקת.
              </li>
              <li>
                בפילוחים עם הרבה ערכים (מדינות, מותגי ספק, סוגי רכב) מוצגים
                שמונת המובילים והשאר מקובצים תחת &quot;אחר&quot;.
              </li>
              <li>
                גרפי הזמן מוצגים משמאל לימין: התאריך המוקדם ביותר בצד שמאל.
              </li>
              <li>
                צבעי הסטטוסים קבועים ומשמשים אך ורק לסטטוס, ולצידם תמיד מופיעה
                תווית טקסט.
              </li>
            </ul>
          </Section>

          <Section title="כרטיסי התשלומים">
            <p className="mb-2">
              גם המספרים האלה מוגבלים לטווח שנבחר, כלומר הם מתייחסים רק להזמנות
              שנוצרו בתקופה הזו ולא ליתרה ההיסטורית הכוללת.
            </p>
            <dl className="space-y-2">
              <Definition term="פתוח לספקים">
                העלות של הזמנות שטרם שולמו לספק, ללא הזמנות מבוטלות. כולל קנסות
                שטרם שולמו לספק.
              </Definition>
              <Definition term="שולם לספקים">
                העלות של הזמנות שכבר סומנו כמשולמות לספק.
              </Definition>
              <Definition term="פתוח מלקוחות">
                הזמנות שהופק להן שובר וטרם שולמו. זיכויים שממתינים להחזר מקטינים
                את הסכום.
              </Definition>
              <Definition term="נגבה מלקוחות">
                הזמנות שסטטוס התשלום שלהן הוא &quot;שולם&quot;.
              </Definition>
            </dl>
            <p className="mt-2 text-xs text-gray-500">
              גרף גיל החוב מפלח את החוב הפתוח לפי הזמן שחלף מרגע יצירת ההזמנה.
            </p>
          </Section>
        </div>
      </DialogContent>
    </Dialog>
  );
}

function Section({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section>
      <h3 className="mb-1.5 font-semibold text-navy">{title}</h3>
      {children}
    </section>
  );
}

function Definition({ term, children }: { term: string; children: ReactNode }) {
  return (
    <div>
      <dt className="inline font-medium text-gray-900">{term}: </dt>
      <dd className="inline">{children}</dd>
    </div>
  );
}
