export class AppError extends Error {
  public readonly code: string;
  public readonly field: string | null;

  constructor(code: string, field: string | null = null) {
    super(code);
    this.code = code;
    this.field = field;
  }
}

export function isAppError(err: unknown): err is AppError {
  return err instanceof AppError;
}
