import { desc, eq, or } from 'drizzle-orm';
import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';
import { db, initDb } from './db';
import { initKafka } from './kafka';
import { sendEmail } from './mailer';
import { notifications, paymentNotifications } from './schema';

const app = new Hono();

app.use('*', logger());
app.use(
  '*',
  cors({
    origin: '*',
    allowHeaders: ['Authorization', 'Content-Type', 'X-Correlation-Id', 'X-Request-Id'],
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS']
  })
);

app.use('*', async (c, next) => {
  const correlationId = c.req.header('X-Correlation-Id') || crypto.randomUUID();
  c.header('X-Correlation-Id', correlationId);
  await next();
});

// Health checks
const healthHandler = (c: any) =>
  c.json({
    status: 'UP',
    components: {
      db: { status: 'UP' }
    }
  });

app.get('/actuator/health', healthHandler);
app.get('/notification/actuator/health', healthHandler);

function registerNotificationRoutes(router: Hono) {
  // Email APIs
  router.post('/api/email/sendSimpleMail', async (c) => {
    try {
      const body = await c.req.json();
      const res = await sendEmail(body);
      return c.text(res);
    } catch (err: any) {
      return c.text(`Error: ${err.message}`, 500);
    }
  });

  router.post('/api/email/sendMailWithAttachment', async (c) => {
    try {
      const body = await c.req.json();
      const res = await sendEmail(body);
      return c.text(res);
    } catch (err: any) {
      return c.text(`Error: ${err.message}`, 500);
    }
  });

  router.post('/api/email/sendMail', async (c) => {
    try {
      const to = c.req.query('to') || '';
      const subject = c.req.query('subject') || 'Notification';
      const bodyText = c.req.query('body') || '';
      const res = await sendEmail({ recipient: to, subject, msgBody: bodyText });
      return c.text(res);
    } catch (err: any) {
      return c.text(`Error: ${err.message}`, 500);
    }
  });

  // Notification CRUD
  router.get('/api/notifications', async (c) => {
    try {
      const rows = await db.select().from(notifications).orderBy(desc(notifications.id));
      return c.json(rows);
    } catch {
      return c.json([], 200);
    }
  });

  router.get('/api/notifications/:id', async (c) => {
    try {
      const id = parseInt(c.req.param('id'), 10);
      const rows = await db.select().from(notifications).where(eq(notifications.id, id));
      if (rows.length === 0) return c.body(null, 404);
      return c.json(rows[0]);
    } catch {
      return c.body(null, 500);
    }
  });

  router.post('/api/notifications', async (c) => {
    try {
      const body = await c.req.json();
      const [row] = await db
        .insert(notifications)
        .values({
          userId: body.userId || '',
          title: body.title || '',
          message: body.message || '',
          type: body.type || 'INFO',
          status: body.status || 'UNREAD'
        })
        .returning();
      return c.json(row, 201);
    } catch (err: any) {
      return c.json({ error: err.message }, 500);
    }
  });

  router.delete('/api/notifications/:id', async (c) => {
    try {
      const id = parseInt(c.req.param('id'), 10);
      await db.delete(notifications).where(eq(notifications.id, id));
      return c.json(true);
    } catch {
      return c.json(false, 500);
    }
  });

  // Payment Notifications CRUD
  router.get('/api/payment-notifications', async (c) => {
    try {
      const rows = await db
        .select()
        .from(paymentNotifications)
        .orderBy(desc(paymentNotifications.id));
      return c.json(rows);
    } catch {
      return c.json([]);
    }
  });

  router.get('/api/payment-notifications/:paymentId', async (c) => {
    try {
      const id = parseInt(c.req.param('paymentId'), 10);
      const rows = await db
        .select()
        .from(paymentNotifications)
        .where(
          or(eq(paymentNotifications.paymentId, id), eq(paymentNotifications.id, id))
        );
      if (rows.length === 0) return c.body(null, 404);
      return c.json(rows[0]);
    } catch {
      return c.body(null, 500);
    }
  });

  router.post('/api/payment-notifications', async (c) => {
    try {
      const body = await c.req.json();
      const [row] = await db
        .insert(paymentNotifications)
        .values({
          paymentId: body.paymentId || null,
          userId: body.userId || null,
          orderId: body.orderId || null,
          amount: String(body.amount || 0),
          isPayed: body.isPayed ?? true,
          paymentStatus: body.paymentStatus || 'SUCCESS'
        })
        .returning();
      return c.json(row);
    } catch (err: any) {
      return c.json({ error: err.message }, 500);
    }
  });

  router.delete('/api/payment-notifications/:paymentId', async (c) => {
    try {
      const id = parseInt(c.req.param('paymentId'), 10);
      await db
        .delete(paymentNotifications)
        .where(
          or(eq(paymentNotifications.paymentId, id), eq(paymentNotifications.id, id))
        );
      return c.body(null, 204);
    } catch {
      return c.body(null, 500);
    }
  });
}

// Register with /notification prefix and direct root
const notificationGroup = new Hono();
registerNotificationRoutes(notificationGroup);
app.route('/notification', notificationGroup);
registerNotificationRoutes(app);

const PORT = parseInt(process.env.PORT || '8090', 10);

console.log(`[Notification Service] Starting Bun server on port ${PORT}...`);
await initDb();
initKafka();

export default {
  port: PORT,
  fetch: app.fetch
};
