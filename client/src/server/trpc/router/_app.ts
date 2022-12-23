import { router } from "../trpc";
import { authRouter } from "./auth";
import { wireGuardRouter } from "./wireguard";

export const appRouter = router({
  auth: authRouter,
  wireguard: wireGuardRouter
});

// export type definition of API
export type AppRouter = typeof appRouter;
