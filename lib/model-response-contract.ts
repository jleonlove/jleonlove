/** Provider-neutral response contracts for the Atlas Cognitive Execution Router.
 * Adapters must be explicitly qualified before enabling external model calls.
 */
export type AtlasMessage = Readonly<{
  role: "system" | "developer" | "user" | "assistant" | "tool";
  content: string;
  toolCallId?: string;
}>;
export type AtlasToolCall = Readonly<{
  id: string;
  name: string;
  argumentsJson: string;
}>;
export type AtlasResponse = Readonly<{
  provider: string;
  model: string;
  text: string;
  toolCalls: readonly AtlasToolCall[];
  finishReason: "stop" | "tool_calls" | "length" | "refusal" | "error";
  inputTokens?: number;
  outputTokens?: number;
}>;
export type AtlasRequest = Readonly<{
  tenantId: string;
  traceId: string;
  model: string;
  messages: readonly AtlasMessage[];
  allowedTools: readonly string[];
}>;
export interface AtlasModelAdapter {
  readonly provider: string;
  generate(request: AtlasRequest, signal: AbortSignal): Promise<AtlasResponse>;
}
export function validateModelResponse(
  request: AtlasRequest,
  response: AtlasResponse,
): AtlasResponse {
  if (!request.tenantId || !request.traceId || !response.provider ||
      !response.model || !["stop", "tool_calls", "length", "refusal", "error"].includes(response.finishReason)) {
    throw new Error("Invalid model response envelope");
  }
  const seen = new Set<string>();
  for (const call of response.toolCalls) {
    if (!call.id || !call.name || seen.has(call.id) ||
        !request.allowedTools.includes(call.name)) {
      throw new Error("Unauthorized or duplicate model tool call");
    }
    seen.add(call.id);
    const parsed: unknown = JSON.parse(call.argumentsJson);
    if (parsed === null || typeof parsed !== "object" || Array.isArray(parsed)) {
      throw new Error("Tool arguments must be a JSON object");
    }
  }
  if (response.toolCalls.length > 0 && response.finishReason !== "tool_calls") {
    throw new Error("Tool calls require tool_calls finish reason");
  }
  return response;
}
