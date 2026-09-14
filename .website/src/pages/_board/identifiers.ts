export function isErrorThread(id: string): boolean {
	return id.startsWith('errors/');
}

export function threadCode(id: string): string {
	const code = id.split('/')[1];
	if (code === undefined) {
		throw new Error(`thread "${id}" carries no code segment`);
	}
	return code;
}
