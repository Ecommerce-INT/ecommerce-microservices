import { locales, localizeHref } from '#lib/paraglide/runtime.js';
import type { RequestHandler } from './$types';

/** Public, indexable storefront routes (private flows stay out — see robots.txt). */
const PAGES = ['/', '/products'];

export const GET: RequestHandler = ({ url }) => {
	const entries = PAGES.flatMap((path) =>
		locales.map((locale) => new URL(localizeHref(path, { locale }), url.origin).href)
	);

	const body = [
		'<?xml version="1.0" encoding="UTF-8"?>',
		'<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
		...entries.map((loc) => `\t<url><loc>${loc}</loc></url>`),
		'</urlset>',
		''
	].join('\n');

	return new Response(body, { headers: { 'content-type': 'application/xml' } });
};
