import Image from "next/image";
import { TypedSection, Populated } from "@/shared/types/payload";
import { SectionHeader } from "../SectionHeader";
import { RichText } from "@payloadcms/richtext-lexical/react";
import { PayloadFormRenderer } from "@/shared/components/forms/FormRenderer";

type ContactSectionProps = {
  section: TypedSection<"contact">;
};

export function ContactSection({ section }: ContactSectionProps) {
  const { title, subtitle, contactForm, contactInfo, eyebrow } =
    section.contact;
  const populatedContactForm = contactForm as Populated<typeof contactForm>;

  return (
    <section
      className="relative w-2/3 mx-auto overflow-hidden rounded-[20px] px-7 py-11 bg-cover bg-center bg-no-repeat"
      style={{ backgroundImage: "url('/assets/contact/contact-bg.png')" }}
    >
      <SectionHeader
        pillText={eyebrow}
        title={title ?? ""}
        subtitle={subtitle}
      />
      <div className="flex gap-12 my-8 items-stretch">
        <div className="w-1/2">
          {typeof populatedContactForm === "object" && populatedContactForm && (
            <PayloadFormRenderer
              form={populatedContactForm}
              title={populatedContactForm.title}
            />
          )}
        </div>
        <div className="w-1/2">
          <div className="flex flex-col gap-4 justify-between h-full">
            {contactInfo?.map((info) => (
              <div className="py-8 px-6 shadow-card flex gap-5 bg-white rounded-xl items-center border border-border" key={info.id}>
                <div className="p-3">
                  <Image
                    src={
                      (info.icon as Populated<typeof info.icon>)?.url ??
                      "/assets/contact/default-icon.png"
                    }
                    alt={
                      (info.icon as Populated<typeof info.icon>)?.alt ??
                      "Contact Icon"
                    }
                    width={64}
                    height={64}
                    className="h-16 w-16 object-contain"
                  />
                </div>
                <div className="flex flex-col gap-0 w-full justify-start">
                  <h5 className="type-h5 text-navy">{info.title}</h5>
                    <RichText data={info.content} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
