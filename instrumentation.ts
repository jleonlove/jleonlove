export async function register() {
  if (process.env.NEXT_RUNTIME === "nodejs") {
    const { registerNodeRuntime } = await import("./instrumentation-node");
    registerNodeRuntime();
  }
}

export async function onRequestError(
  error: unknown,
  request: { path?: string; method?: string },
  context: { routeType?: string; routePath?: string },
) {
  if (process.env.NEXT_RUNTIME !== "nodejs") return;
  const { observeNodeRequestError } = await import("./instrumentation-node");
  observeNodeRequestError(error, request, context);
}
