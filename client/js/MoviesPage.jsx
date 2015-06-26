const 	React 			= require('react'),
		Link			= require('react-router').Link,
		RouteHandler 	= require('react-router').RouteHandler,
		Pager 			= require('react-paginate'),
		Movie 			= require('./Movie.jsx'),
		// DatePicker 		= require('react-datepicker'),
		// moment 			= require('moment'),
		optionsBO 		= require('./options.js')


// const DatePickerWrapper = React.createClass({
// 	// getInitialState: function() {
// 	// 	return {
// 	// 		date: moment()
// 	// 	}
// 	// },

// 	render: function() {
// 		<span>
// 			<DatePicker
// 				key="example1"
// 				selected={movie.last_watched}
// 				onChange={handleWatched}
// 			/>

// 		</span>
// 	}
// })

module.exports = React.createClass({
	// componentWillUpdate: function() {
	// 	console.log('this.selected', this.selected)
	// 	console.log('selected', this.state.selected)
	// 	this.forced = 0
	// 	if (this.selected) {
	// 		this.forced = this.selected
	// 	}		
	// },
	handlePageClick: function(data) {
		this.shouldScroll = true
		// this.selected = data.selected
		// const offset = Math.ceil(this.selected * this.props.options.limit);
		const offset = Math.ceil(data.selected * this.props.options.limit);
		optionsBO.setOffset(offset)		
	},

	componentDidUpdate: function() {
		if (this.shouldScroll) {
			window.scrollTo(0, 0)
		}

		this.shouldScroll = false

		// console.log('didUdpate')
	},	

	render: function() {
		const movies = this.props.movies
		const options = this.props.options

		const selected = options.offset / options.limit

		console.log('offset: ' + options.offset + ' limit: ' + options.limit)
		console.log('selected: ' + selected)

		var pagination;
		if (movies.total > options.limit) {
			// console.log('moviesPage.total('+movies.total+')>limit('+options.limit+'); selected='+this.selected)

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

		var items = movies.items.map(function(movie) {
			return (
				<Movie movie={movie} key={movie.title+movie.modified} />
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
})