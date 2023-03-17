# go-globby

fast glob. but use go, so high-speed

<br />

## Usage

### install

```shell
npm i go-globby
```

### program

```ts
import { globbyBin } from 'go-globby'

await globbyBin([
	'**/*.json', // Recursive acquisition
	'!package.json', // Ignore the current directory's package.json
	'!node_modules/**/*' //  Ignore node_modules folder
])
```

<br />

## License

Made with [markthree](https://github.com/markthree)

Published under [MIT License](./LICENSE).
