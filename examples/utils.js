import { check } from "k6";
export function logError(err) {
  const e = new Error(err);
  console.error(e.stack);
}

export function validateResult(res, msg) {
  const successful = res.err == undefined;
  check(res, {
    [msg]: () => successful,
  });
  if (!successful) {
    logError(res.err);
  }
  return successful;
}
