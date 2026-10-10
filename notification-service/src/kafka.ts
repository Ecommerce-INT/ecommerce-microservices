import { Kafka } from "kafkajs";
import { sendEmail } from "./mailer";
import { sql } from "./db";

const kafkaBrokers = (process.env.KAFKA_SERVERS || "localhost:9092").split(",");
const clientId = "notification-service";
const groupId = process.env.PAYMENT_KAFKA_CONSUMER_GROUP_ID || "notification-group";

export async function initKafka() {
  try {
    const kafka = new Kafka({
      clientId,
      brokers: kafkaBrokers,
      retry: {
        retries: 2,
        initialRetryTime: 300,
      },
    });

    const consumer = kafka.consumer({ groupId });
    await consumer.connect();
    await consumer.subscribe({ topics: ["profile-onboarding-topic", "status-payment-successful"], fromBeginning: false });

    console.log("[Notification Service] Kafka consumer connected and subscribed to topics.");

    consumer.run({
      eachMessage: async ({ topic, message }) => {
        try {
          const raw = message.value?.toString() || "{}";
          console.log(`[Notification Service] Received event on topic [${topic}]:`, raw);

          if (topic === "profile-onboarding-topic") {
            const data = JSON.parse(raw);
            await sendEmail({
              recipient: data.recipient || data.email || "customer@ecommerce.local",
              subject: data.subject || "Welcome to Ecommerce!",
              msgBody: data.msgBody || "Thank you for onboarding!",
            });
          } else if (topic === "status-payment-successful") {
            const payment = JSON.parse(raw);
            await sql`
              INSERT INTO payment_notifications (
                payment_id, user_id, order_id, amount, is_payed, payment_status
              ) VALUES (
                ${payment.paymentId || null}, ${payment.userId || null}, ${payment.orderId || null},
                ${payment.amount || 0}, ${payment.isPayed ?? true}, ${payment.paymentStatus || "SUCCESS"}
              )
            `;

            await sendEmail({
              recipient: "hoangtien2k3dev@gmail.com",
              subject: `Payment Successful in Order with userId: ${payment.userId}`,
              msgBody: `Payment in order product cart successfully.\nIsPayed: ${payment.isPayed}\nStatus: ${payment.paymentStatus}\nTime: ${new Date().toISOString()}`,
            });
          }
        } catch (err: any) {
          console.error(`[Notification Service] Error handling Kafka message:`, err);
        }
      },
    });
  } catch (err: any) {
    console.warn("[Notification Service] Notice: Kafka offline or unavailable during startup (continuing without Kafka):", err.message);
  }
}
