export type ApiError = { code: string; message: string };
export type ApiResponse<T> = {
  success: boolean;
  data?: T;
  error?: ApiError;
  requestId: string;
};

export function requestId(request?: Request) {
  return request?.headers.get("x-request-id") ?? crypto.randomUUID();
}

export function ok<T>(data: T, request?: Request, init?: ResponseInit) {
  return Response.json(
    { success: true, data, requestId: requestId(request) } satisfies ApiResponse<T>,
    init,
  );
}

export function fail(code: string, message: string, status: number, request?: Request) {
  return Response.json(
    { success: false, error: { code, message }, requestId: requestId(request) } satisfies ApiResponse<never>,
    { status },
  );
}
