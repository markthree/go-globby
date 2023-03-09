import { globby } from 'globby'
import { globbyBin } from './npm/index.mjs'

console.time('globbyBin')
const result = await globbyBin(['**/*.json'], {
	cwd: 'node_modules'
})
console.timeEnd('globbyBin')

console.time('globby')
const result2 = await globby(['**/*.json'], {
	cwd: 'node_modules',
	dot: true
})
console.timeEnd('globby')
