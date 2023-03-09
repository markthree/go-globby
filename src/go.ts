import { execa } from 'execa'
import {
	arch as _arch,
	platform as _platform
} from 'node:os'

import { fileURLToPath } from 'node:url'
import { cwd as _cwd } from 'node:process'
import { dirname, resolve } from 'node:path'
import { createInterface } from 'node:readline'

const _dirname =
	typeof __dirname !== 'undefined'
		? __dirname
		: dirname(fileURLToPath(import.meta.url))

let defaultBinPath = ''

export function inferVersion() {
	const platform = _platform()
	if (!/win32|linux|darwin/.test(platform)) {
		throw new Error(`${platform} is not support`)
	}

	const arch = _arch()

	if (!/amd64_v1|arm64|386|x64/.test(arch)) {
		throw new Error(`${arch} is not support`)
	}

	return `${platform === 'win32' ? 'windows' : platform}_${
		arch === 'x64' ? 'amd64_v1' : arch
	}`
}

export function detectBinName(version = inferVersion()) {
	return `go-globby${
		version.startsWith('windows') ? '.exe' : ''
	}`
}

export function detectDefaultBinPath() {
	if (defaultBinPath) {
		return defaultBinPath
	}

	const version = inferVersion()
	const name = detectBinName(version)
	defaultBinPath = resolve(
		_dirname,
		`../dist/go-globby_${version}/${name}`
	)
	return defaultBinPath
}

interface Options {
	cwd?: string
	binPath?: string
}

export function globbyBin(
	globs: Array<string>,
	options: Options = {}
) {
	const { cwd = _cwd(), binPath = detectDefaultBinPath() } =
		options

	return new Promise<string[]>(async (resolve, reject) => {
		const go = execa(binPath, {
			cwd
		})
		const readline = createInterface({
			input: go.stdout
		})
		const paths: string[] = []

		function close() {
			if (!go.killed) {
				readline.close()
				go.cancel()
			}
		}

		readline.on('line', (path: string) => {
			if (path === 'EOF') {
				resolve(paths)
				close()
				return
			}
			paths.push(path)
		})

		go.stderr.on('data', (error: Buffer) => {
			reject(String(error))
			close()
		})

		process.once('exit', () => {
			close()
		})

		go.stdin.write(`${globs}\n`)
	})
}
