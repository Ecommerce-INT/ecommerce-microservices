/**
 * Internal origin of the API gateway (Apache APISIX) used by the SvelteKit
 * server. The browser never calls the gateway directly — every request goes
 * through remote functions on the server, so this URL only has to be reachable
 * from the app container/pod.
 *
 *   k8s      → http://apisix-gateway:9080
 *   compose  → http://apisix:9080
 *   local    → http://localhost:9080 (or http://api.ecommerce.local with hosts entries)
 */
export function apiBaseUrl(): string {
	const url = process.env.API_INTERNAL_URL ?? 'http://localhost:9080';
	return url.replace(/\/+$/, '');
}

/**
 * Browser-reachable origin of the gateway. Only needed to hand the browser off
 * to the SSO login endpoint (`/api/v1/auth/login`) — every other call stays
 * server-side. Defaults to the internal URL, which is already correct locally.
 *
 *   k8s/compose → https://api.example.com (ingress host)
 *   local       → http://localhost:9080 (default)
 */
export function apiPublicUrl(): string {
	const url = process.env.API_PUBLIC_URL ?? apiBaseUrl();
	return url.replace(/\/+$/, '');
}
