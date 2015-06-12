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
	// getInitialState: function() {
	// 	return {
	// 		date: moment()
	// 	}
	// },

	render: function() {
		const movies = this.props.movies
		const options = this.props.options

		const handlePageClick = function(data) {
			const selected = data.selected;
			const offset = Math.ceil(selected * options.limit);
			optionsBO.setOffset(offset)
		}

		const handleWatched = function(data) {
			return;
		}

		var pagination;
		if (movies.total > options.limit) {
			console.log('moviesPage.total('+movies.total+'>limit('+options.limit)
			pagination = (
		        <Pager previousLabel={<i className="icon-chevron-left"></i>}
		                       nextLabel={<i className="icon-chevron-right"></i>}
		                       breakLabel={<li className="break"><a href="">...</a></li>}
		                       pageNum={Math.ceil(movies.total / options.limit)}
		                       marginPagesDisplayed={3}
		                       pageRangeDisplayed={5}
		                       clickCallback={handlePageClick}
		                       containerClassName={"pagination col-xs-12"}
		                       subContainerClassName={"pages"}
		                       activeClass={"active"} />				
			)
		}

		var that = this;
		var items = movies.items.map(function(movie, i) {
			// var watched;

			// if (movie.last_watched != '') {
			// 	watched = (
			// 		<span className="label success spacer"><i className="icon-watched"></i>&nbsp;{moment(movie.last_watched).format('MMM DD, YYYY')}</span>
			// 	)
			// }

			return (
				<Movie movie={movie} key={i} />

				// <article key={i}>
				// 	<div className="col-xs-12">
				// 		<h2>{movie.title} ({movie.year})</h2>
				// 	</div>
				// 	<div className="col-xs-12">
				// 		<div className="row moviep-images">
				// 			<div className="col-xs-12 col-sm-2">
				// 				<img src={"/img/p" + movie.cover} />
				// 			</div>
				// 			<div className="col-xs-12 col-sm-10">
				// 				<img src={"/img/b" + movie.backdrop} />
				// 			</div>
				// 		</div>
				// 	</div>
				// 	<div className="col-xs-12">
				// 		<div className="row between-xs">
				// 			<span className="col-xs-12 col-sm-6 director">{movie.director}</span>
				// 			<span className="col-xs-12 col-sm-6 end-sm">{movie.production_countries}</span>
				// 		</div>
				// 	</div>
				// 	<div className="col-xs-12">
				// 		<div className="row between-xs">
				// 			<span className="col-xs-12 col-sm-6 ">{movie.actors}</span>
				// 			<span className="col-xs-12 col-sm-6 end-sm">{movie.genres}</span>
				// 		</div>
				// 	</div>
				// 	<div className="col-xs-12 bottom-spacer">
				// 		<div className="row between-xs">
				// 			<div className="col-xs-12 col-sm-9">
				// 				<span className="label">{movie.resolution}</span>
				// 				<span className="label secondary spacer">{movie.location}</span>
				// 			</div>
				// 			<div className="col-xs-12 col-sm-3 end-sm">
				// 				{watched}							
				// 				<span className="label"><i className="icon-plus"></i>&nbsp;{moment(movie.added).format('MMM DD, YYYY H:mm')}</span>
				// 			</div>
				// 		</div>
				// 	</div>
				// 	<div className="col-xs-12">
				// 		<span>{movie.overview}</span>
				// 	</div>
				// 	<div className="col-xs-12 bottom-spacer-large">
				// 		<div className="row between-xs">
				// 			<div className="col-xs-12 col-sm-10 top-xs date-picker">
				// 				<input type="text"></input>
				// 				<button className="btn btn-default">Fix</button>
				// 				<DatePicker
				// 					key="example1"
				// 					selected={moment()}
				// 					onChange={handleWatched}
				// 				/>
				// 			</div>
				// 			<div className="col-xs-12 col-sm-2 end-sm top-xs">
				// 				<span className="label"><i className="icon-watched"></i></span>
				// 			</div>
				// 		</div>
				// 	</div>										
				// </article>				
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