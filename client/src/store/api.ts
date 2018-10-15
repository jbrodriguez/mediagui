import to from 'await-to-js'
import { ConfigState, Movie } from '@/types'
import { Config } from '@fortawesome/fontawesome-svg-core'

const encode = (data: any) => {
	const encoded = Object.keys(data).map(key => {
		const value = encodeURIComponent(data[key].toString())
		return `${key}=${value}`
	})
	return encoded.join('&')
}

const checkStatus = (response: Response) => {
	if (!response.ok) {
		throw Error(`${response.status}: ${response.statusText}`)
	}

	return response
}

const retrieve = (params: string, options?: any): Promise<Response> => {
	return fetch(params, options).then(checkStatus)
}

class Api {
	private ep: string =
		document && document.location
			? `${document.location.protocol}//${document.location.hostname}:7623/api/v1`
			: `http://blackbeard.apertoire.org:7623/api/v1`

	public async getConfig(): Promise<ConfigState> {
		let config: ConfigState = {
			version: '1.0.0',
			unraidMode: true,
			unraidHosts: [],
			mediaFolders: [],
		}

		const [err, data] = await to<Response>(retrieve(`${this.ep}/config`))
		if (err) {
			// console.log(`reply.err(${err})`)
		} else {
			if (data) {
				config = await data.json()
				// console.log(`data(${JSON.stringify(config)})`)
			}
		}

		return config
	}

	public async getMovies(options: any): Promise<Movie[]> {
		let movies: Movie[] = []

		const [err, data] = await to<Response>(retrieve(`${this.ep}/movies?${encode(options)}`))
		if (err) {
			// console.log(`reply.err(${err})`)
		} else {
			if (data) {
				movies = await data.json()
				// console.log(`data(${d})`)
			}
		}

		return movies
	}

	// public async loadStatement(filepath: string): Promise<Statement | null> {
	// 	let statement: Statement | null = null

	// 	const query = encode({ filepath })
	// 	const [err, data] = await to<Response>(retrieve(`${this.ep}/loadStatement?${query}`))
	// 	if (err) {
	// 		// console.log(`reply.err(${err})`)
	// 	} else {
	// 		if (data) {
	// 			statement = await data.json()
	// 			// console.log(`data(${d})`)
	// 		}
	// 	}

	// 	return statement
	// }

	// public async saveStatement(statement: Statement | null) {
	// 	const [err, data] = await to<Response>(
	// 		retrieve(`${this.ep}/saveStatement`, {
	// 			method: 'POST',
	// 			headers: new Headers({ 'Content-Type': 'application/json' }),
	// 			body: JSON.stringify(statement),
	// 		}),
	// 	)

	// 	if (err) {
	// 		// console.log(`reply.err(${err})`)
	// 	} else {
	// 		// if (data) {
	// 		// 	statement = await data.json()
	// 		// 	// console.log(`data(${d})`)
	// 		// }
	// 	}

	// 	return statement
	// }

	// public async listStatements(): Promise<Statement[] | null> {
	// 	let statements: Statement[] = []

	// 	const [err, data] = await to<Response>(retrieve(`${this.ep}/listStatements`))
	// 	if (err) {
	// 		// console.log(`reply.err(${err})`)
	// 	} else {
	// 		if (data) {
	// 			statements = await data.json()
	// 			// console.log(`data(${d})`)
	// 		}
	// 	}

	// 	return statements
	// }

	// public async getStmtPlot(id: number): Promise<Plot | null> {
	// 	let plot: Plot | null = null

	// 	const query = encode({ id })
	// 	const [err, data] = await to<Response>(retrieve(`${this.ep}/stmtPlot?${query}`))
	// 	if (err) {
	// 		// console.log(`reply.err(${err})`)
	// 	} else {
	// 		if (data) {
	// 			plot = await data.json()
	// 			// console.log(`data(${d})`)
	// 		}
	// 	}

	// 	return plot
	// }

	// public async getHistory(): Promise<Plot | null> {
	// 	let plot: Plot | null = null

	// 	const [err, data] = await to<Response>(retrieve(`${this.ep}/getHistory`))
	// 	if (err) {
	// 		// console.log(`reply.err(${err})`)
	// 	} else {
	// 		if (data) {
	// 			plot = await data.json()
	// 			// console.log(`data(${d})`)
	// 		}
	// 	}

	// 	return plot
	// }
}

const api = new Api()

export default api

// function tony<T, U = any>(promise: Promise<T>, errorExt?: object): Promise<[U | null, T | undefined]> {
// 	console.log('tonifier')
// 	return promise
// 		.then<[null, T]>((data: T) => {
// 			console.log(`tonydata(${data})`)
// 			return [null, data]
// 		})
// 		.catch<[U, undefined]>(err => {
// 			console.log(`hereamin-(${errorExt})-err(${err})`)
// 			if (errorExt) {
// 				Object.assign(err, errorExt)
// 			}

// 			console.log(`errinside(${err})`)
// 			return [err, undefined]
// 		})
// }
