import amqplib from "amqplib/callback_api";
import { env } from "../../env/server.mjs";
const { RABBIT_MQ_PORT, RABBIT_MQ_USER, RABBIT_MQ_PASSWORD, RABBIT_MQ_HOST } = env;

let conn: amqplib.Connection | undefined;

let channel: amqplib.Channel | undefined;

export const amqpConnection = async () => {
    if(conn) return conn;
    return new Promise<amqplib.Connection>((resolve) => {
        amqplib.connect(`amqp://${RABBIT_MQ_USER}:${RABBIT_MQ_PASSWORD}@${RABBIT_MQ_HOST}:${RABBIT_MQ_PORT}`, (err, conn) => {
            if(err) {
                throw err
            }
            resolve(conn);
        });
    })
}

export const amqpChannel = async () => {
    if(channel) return channel;
    const conn = await amqpConnection();
    return new Promise<amqplib.Channel>((resolve, reject) => {
        conn.createChannel(async (err, channel) => {
            if(err) reject(err)
            resolve(channel);
        })
    })
}