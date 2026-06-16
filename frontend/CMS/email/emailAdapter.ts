import { getBaseURL } from "@/shared/api/_api";
import Client from "@/shared/client";
import type { EmailAdapter } from "payload";

export function emailAdapter(): EmailAdapter {
  return () => ({
    name: "email-adapter",
    defaultFromAddress: "",
    defaultFromName: "",

    sendEmail: async (message) => {
      const client = new Client(getBaseURL(), {
        requestInit: {
          headers: {
            "X-Service-Client-Token": process.env.PAYLOAD_EMAIL_API_KEY || "",
          },
        },
      });

      try {
        await client.notifications.SendCMSEmail({
          to: message.to,
          subject: message.subject,
          content: message.html,
        });

        return;
      } catch (error) {
        console.error("Error sending email:", error);
        throw error;
      }
    },
  });
}
