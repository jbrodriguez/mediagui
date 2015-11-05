import React from 'react'
import Pager from 'react-paginate'

import Movie from './movie_card'
import { isNotValid } from '../lib/utils'


export default class DuplicatesPage extends React.Component {
	constructor() {
		super()

		this.handlePageClick = this.handlePageClick.bind(this)

		this.shouldScroll = false
	}

	componentWillMount() {
		this.props.actions.movies.getDuplicates()
	}

	componentDidUpdate() {
		if (this.shouldScroll) {
			window.scrollTo(0, 0)
		}

		this.shouldScroll = false
	}

	render() {
		if (isNotValid(this.props.state.movies)) {
			return null
		}

		const movies = this.props.state.movies
		const options = this.props.state.options

		const selected = options.offset / options.limit

		var pagination;
		if (movies.total > options.limit) {
			pagination = (
		        <Pager previousLabel={<i className="icon-chevron-left"></i>}
		                       nextLabel={<i className="icon-chevron-right"></i>}
		                       breakLabel={<li className="break"><a href="">...</a></li>}
		                       pageNum={Math.ceil(movies.total / options.limit)}
		                       marginPagesDisplayed={3}
		                       pageRangeDisplayed={5}
		                       forceSelected={selected}
		                       clickCallback={this.handlePageClick}
		                       containerClassName={"pagination col-xs-12"}
		                       subContainerClassName={"pages"}
		                       activeClass={"active"} />				
			)
		}

		var that = this
		var items = movies.items.map(function(movie, i) {
			return (
				<Movie movie={movie} key={movie.title+movie.modified+i} { ...that.props} />
			)
		})

		return (
			<section className="row moviep">
				{pagination}

				{items}

				{pagination}
			</section>
		)		
	}

	handlePageClick(data) {
		this.shouldScroll = true

		const offset = Math.ceil(data.selected * this.props.state.options.limit);
		this.props.actions.options.setOffset(offset)
	}
}