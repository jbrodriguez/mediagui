import React from 'react'
import Cover from './cover_card'
import { isNotValid } from '../lib/utils'

export default class CoversPage extends React.Component {
	componentWillMount() {
		const proxy = Object.assign({}, this.props.state.options, {limit: 60})
		this.props.actions.movies.getMovies(proxy)
	}

	render() {
		if (isNotValid(this.props.state.movies)) {
			return null
		}

		const movies = this.props.state.movies.items

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