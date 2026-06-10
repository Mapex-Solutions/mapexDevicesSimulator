/**
 * Raised when a request payload fails its Zod schema before reaching the sidecar.
 * `messages` holds one `path message` line per issue so a caller can surface them.
 */
export class SchemaError extends Error {
	public messages: string[];

	constructor(messages: string[]) {
		super('SchemaError');
		this.name = 'SchemaError';
		this.messages = messages;
	}
}

/** Flatten a failed safeParse result into readable `path message` lines. */
export function zodValidationError(result: { error: { issues: Array<{ path: Array<string | number>; message: string }> } }): string[] {
	return result.error.issues.map((issue) => {
		const message = String(issue.message).toLowerCase();
		return `${issue.path.join('.')} ${message}`.trim();
	});
}
