import { globby } from 'globby'
import { globbyBin } from './npm/index.mjs'

const binStartTime = Date.now()
const binResult = await globbyBin([
	'**/*.json',
	'!package.json',
	'!node_modules/**/*'
])
const binDuration = Date.now() - binStartTime

const nodeStartTime = Date.now()
const nodeResult = await globby([
	'**/*.json',
	'!package.json',
	'!node_modules/**/*'
])
const nodeDuration = Date.now() - nodeStartTime

console.log(
	`globby    - duration: ${
		nodeDuration / 1000
	}s; result: ${nodeResult}`
)
console.log(
	`go-globby - duration: ${
		binDuration / 1000
	}s; result: ${binResult}`
)

console.log('\n')

console.log(
	'go-globby vs globby -',
	(nodeDuration / binDuration).toFixed(2) + ' ↑'
)
