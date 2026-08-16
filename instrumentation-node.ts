import packageJson from "./package.json";
import { createRuntimeObservation, emitRuntimeObservation } from "./lib/runtime-observability";

const release = packageJson.atlasRelease;

export function registerNodeRuntime() {
  emitRuntimeObservation(createRuntimeObservation({
    level: "info",
    event: "runtime.ready",
    release,
    environment: process.env.VERCEL_ENV ?? process.env.NODE_ENV ?? "unknown",
    attributes: { runtime: process.env.NEXT_RUNTIME ?? "nodejs", node: process.version },
  }));
}

export function observeNodeRequestError(
  error: unknown,
  request: { path?: string; method?: string },
  context: { routeType?: string; routePath?: string },
) {
  const cause = error instanceof Error
    ? { name: error.name, message: error.message }
    : { name: "UnknownError", message: String(error) };

  emitRuntimeObservation(createRuntimeObservation({
    level: "error",
    event: "request.error",
    release,
    environment: process.env.VERCEL_ENV ?? process.env.NODE_ENV ?? "unknown",
    attributes: {
      ...cause,
      method: request.method,
      path: request.path,
      routeType: context.routeType,
      routePath: context.routePath,
    },
  }));
}
