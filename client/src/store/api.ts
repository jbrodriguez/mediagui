import to from 'await-to-js'
import { ConfigState, Movie, Movies } from '@/types'
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

	public async setMovieScore(movie: Movie): Promise<Movie> {
		let changed: Movie = { ...movie }

		const [err, data] = await to<Response>(
			retrieve(`${this.ep}/movies/${movie.id}/score`, {
				method: 'PUT',
				headers: new Headers({ 'Content-Type': 'application/json' }),
				body: JSON.stringify(movie),
			}),
		)
		if (err) {
			// console.log(`reply.err(${err})`)
		} else {
			if (data) {
				changed = await data.json()
				// console.log(`data(${d})`)
			}
		}

		return changed
	}

	public async setMovieWatched(movie: Movie): Promise<Movie> {
		let changed: Movie = { ...movie }

		const [err, data] = await to<Response>(
			retrieve(`${this.ep}/movies/${movie.id}/watched`, {
				method: 'PUT',
				headers: new Headers({ 'Content-Type': 'application/json' }),
				body: JSON.stringify(movie),
			}),
		)
		if (err) {
			// console.log(`reply.err(${err})`)
		} else {
			if (data) {
				changed = await data.json()
				// console.log(`data(${d})`)
			}
		}

		return changed
	}

	public async fixMovie(movie: Movie): Promise<Movie> {
		let changed: Movie = { ...movie }

		const [err, data] = await to<Response>(
			retrieve(`${this.ep}/movies/${movie.id}/fix`, {
				method: 'PUT',
				headers: new Headers({ 'Content-Type': 'application/json' }),
				body: JSON.stringify(movie),
			}),
		)
		if (err) {
			// console.log(`reply.err(${err})`)
		} else {
			if (data) {
				changed = await data.json()
				// console.log(`data(${d})`)
			}
		}

		return changed
	}

	public async setMovieDuplicate(movie: Movie): Promise<Movie> {
		let changed: Movie = { ...movie }

		const [err, data] = await to<Response>(
			retrieve(`${this.ep}/movies/${movie.id}/duplicate`, {
				method: 'PUT',
				headers: new Headers({ 'Content-Type': 'application/json' }),
				body: JSON.stringify(movie),
			}),
		)
		if (err) {
			// console.log(`reply.err(${err})`)
		} else {
			if (data) {
				changed = await data.json()
				// console.log(`data(${d})`)
			}
		}

		return changed
	}

	public async importMovies() {
		const [err, data] = await to<Response>(retrieve(`${this.ep}/import`, { method: 'POST' }))
		if (err) {
			// console.log(`reply.err(${err})`)
		} else {
			if (data) {
				// changed = await data.json()
				// console.log(`data(${d})`)
			}
		}
	}

	public async addMovie(title: string) {
		let changed: Movie = { title } // here i should initialize an empty Movie

		const [err, data] = await to<Response>(
			retrieve(`${this.ep}/add`, {
				method: 'POST',
				headers: new Headers({ 'Content-Type': 'application/json' }),
				body: JSON.stringify(changed),
			}),
		)
		if (err) {
			// console.log(`reply.err(${err})`)
		} else {
			if (data) {
				changed = await data.json()
				// console.log(`data(${d})`)
			}
		}

		return changed
	}

	public async addFolder(folder: string) {
		const [err, data] = await to<Response>(
			retrieve(`${this.ep}/config/folder`, {
				method: 'PUT',
				headers: new Headers({ 'Content-Type': 'application/json' }),
				body: JSON.stringify({ topic: '', payload: folder }),
			}),
		)
		if (err) {
			// console.log(`reply.err(${err})`)
		} else {
			if (data) {
				// changed = await data.json()
				// console.log(`data(${d})`)
			}
		}
	}

	public async getDuplicates(): Promise<Movies> {
		let movies: Movies = { total: 0, items: [] }

		const [err, data] = await to<Response>(retrieve(`${this.ep}/movies/duplicates`))
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
}

const api = new Api()

export default api
