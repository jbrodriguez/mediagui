import React from 'react'
import Cover from './cover_card'
import { isNotValid } from '../lib/utils'

export default class CoversPage extends React.Component {
	componentWillMount() {
		// console.log('conversPage.componentWillMount')
		// const options = this.props.state.options

		// const proxy = {
		// 	query: options.query,
		// 	filterBy: options.filterBy,
		// 	sortBy: options.sortBy,
		// 	sortOrder: options.sortOrder,
		// 	limit: 60,
		// 	offset: options.offset
		// }

		const proxy = Object.assign({}, this.props.state.options, {limit: 60})

		this.props.actions.movies.getMovies(proxy)

		// this.props.actions.movies.getCover()
	}

	componentWillReceiveProps(nextProps) {
		// console.log('nextProps: ', nextProps)
	}

	render() {
		// console.log('CoversPage.rendering')
		// console.log('this.props: ', this.props)
		// this.props.actions.movies.getCover()
		// console.log('coverspage.this.props: ', this.props)
		if (isNotValid(this.props.state.movies)) {
			return null
		}

		const movies = this.props.state.movies.items

		// console.log('movies: ', movies)

		var items = movies.map(function(movie) {
			return (
				<Cover movie={movie} key={movie.title} />
			)
		})

		return (
			<section className="row">
				{items}
			</section>
		)		
	}
}