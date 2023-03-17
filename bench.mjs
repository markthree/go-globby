import { globby } from 'globby'
import { globbyBin } from './npm/index.mjs'

console.time('globbyBin')
const result = await globbyBin(['*.json', '!package.json'])
console.timeEnd('globbyBin')

console.time('globby')
const result2 = await globby(['*.json', '!package.json'])
console.timeEnd('globby')

console.log('globbyBin', result)
console.log('globby', result2)
