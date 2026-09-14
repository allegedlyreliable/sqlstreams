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

export function isDecisionRecordThread(id: string): boolean {
	return id.startsWith('decisions/');
}

export function recordNumber(id: string): string {
	const number = id.split('/')[1];
	if (number === undefined) {
		throw new Error(`thread "${id}" carries no record number segment`);
	}
	return number;
}
