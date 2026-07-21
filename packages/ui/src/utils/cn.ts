import { clsx } from "clsx";

export const cn = (...inputs: (string | undefined | null | false | Record<string, boolean | undefined | null>)[]) =>
  clsx(inputs);

export const pcn = (base: string, props?: Record<string, boolean | undefined | null>) =>
  clsx(base, props);
