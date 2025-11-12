import { z } from "zod";

export type Result<T, E> = { ok: true; data: T } | { ok: false; err: E };

export const ok = <T, E>(data: T): Result<T, E> => ({ ok: true, data });
export const err = <T, E>(err: E): Result<T, E> => ({ ok: false, err });

export const ERROR_SCHEMA = z.object({ detail: z.string() });
