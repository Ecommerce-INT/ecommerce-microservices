import { existsSync, readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';

const root = join(import.meta.dir, '..');
const appsDir = join(root, 'apps');

let failed = false;

for (const entry of readdirSync(appsDir, { withFileTypes: true })) {
	if (!entry.isDirectory()) continue;
	const messagesDir = join(appsDir, entry.name, 'messages');
	if (!existsSync(messagesDir)) continue;

	const catalogs = readdirSync(messagesDir)
		.filter((file) => file.endsWith('.json'))
		.sort();
	if (catalogs.length < 2) continue;

	const baseFile = catalogs.includes('vi.json') ? 'vi.json' : catalogs[0];
	const baseKeys = new Set(
		Object.keys(JSON.parse(readFileSync(join(messagesDir, baseFile), 'utf8')))
	);

	for (const file of catalogs) {
		if (file === baseFile) continue;
		const keys = new Set(Object.keys(JSON.parse(readFileSync(join(messagesDir, file), 'utf8'))));
		const missing = [...baseKeys].filter((key) => !keys.has(key));
		const extra = [...keys].filter((key) => !baseKeys.has(key));

		if (missing.length || extra.length) {
			failed = true;
			console.error(`${entry.name}: ${file} differs from ${baseFile}`);
			if (missing.length) console.error(`  missing: ${missing.join(', ')}`);
			if (extra.length) console.error(`  extra:   ${extra.join(', ')}`);
		} else {
			console.log(`${entry.name}: ${baseFile} / ${file} in sync (${keys.size} keys)`);
		}
	}
}

process.exit(failed ? 1 : 0);
