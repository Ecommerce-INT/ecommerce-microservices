/** Only in-app absolute paths are accepted as post-login destinations. */
export function safeRedirect(target: string | null | undefined): string {
	if (target && target.startsWith('/') && !target.startsWith('//')) return target;
	return '/';
}
