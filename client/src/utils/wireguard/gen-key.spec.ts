import { generateWireGuardKeypair } from "./gen-key";
import z from "zod";
describe("generateWireGuardKeypair", () => {
    test("generates a valid wireguard keypair", async () => {
        z.object({
            publicKey: z.string().min(16).endsWith("="),
            privateKey: z.string().min(16).endsWith("=")
        }).parse(await generateWireGuardKeypair());
    });
})