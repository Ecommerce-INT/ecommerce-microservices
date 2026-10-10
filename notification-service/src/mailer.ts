import nodemailer from "nodemailer";
import type { Transporter } from "nodemailer";

const host = process.env.MAIL_HOST || "smtp.gmail.com";
const port = parseInt(process.env.MAIL_PORT || "587", 10);
const user = process.env.MAIL_USERNAME || "";
const pass = process.env.MAIL_PASSWORD || "";

let transporter: Transporter | null = null;

if (user && pass) {
  transporter = nodemailer.createTransport({
    host,
    port,
    secure: port === 465,
    auth: { user, pass },
  });
}

export interface EmailDetails {
  recipient: string;
  msgBody: string;
  subject: string;
  attachment?: string;
}

export async function sendEmail(details: EmailDetails): Promise<string> {
  console.log(`[Notification Service] Sending email to: ${details.recipient}, subject: ${details.subject}`);

  if (!transporter) {
    console.log(`[Notification Service] SMTP not configured. Mock email dispatched: ${details.subject}`);
    return "Mail Sent Successfully (Mock)";
  }

  try {
    await transporter.sendMail({
      from: user || "no-reply@ecommerce.com",
      to: details.recipient,
      subject: details.subject,
      text: details.msgBody,
    });
    return "Mail Sent Successfully...";
  } catch (err: any) {
    console.error("[Notification Service] Error while sending email:", err);
    return `Error while Sending Mail: ${err.message}`;
  }
}
