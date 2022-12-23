import { z } from "zod";
import { publicProcedure, router } from "../trpc";
import { generateWireGuardKeypair } from "../../../utils/wireguard/gen-key";
export const wireGuardRouter = router({
    createClient: publicProcedure
    .input(z.object({
        for: z.string(),
        expires: z.number().int().min(-1)
    }))
    .query(async ({ input, ctx })  => {
        const keyPair = await generateWireGuardKeypair();
        ctx.amqp.sendToQueue("wireguard-client-create", Buffer.from(JSON.stringify({
            ...input,
            ...keyPair
        })));
    }),
});