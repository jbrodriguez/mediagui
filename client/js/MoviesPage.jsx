const 	React 			= require('react'),
		Link			= require('react-router').Link,
		RouteHandler 	= require('react-router').RouteHandler,
		Pager 			= require('react-paginate'),
		moment 			= require('moment'),
		optionsBO 		= require('./options.js')

module.exports = React.createClass({
	render: function() {
		const movies = this.props.movies
		const options = this.props.options

		var pagination;

		const handlePageClick = function(data) {
			const selected = data.selected;
			const offset = Math.ceil(selected * options.limit);
			optionsBO.setOffset(offset)
		}

		console.log('where is the love: total('+movies.total+')>limit('+options.limit+')')

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

		// const styles = {height: "17em"}
		// const styleo = {overflow: "hidden", maxHeight: "17em"}

		var items = movies.items.map(function(movie, i) {
			return (
				<article key={i}>
					<div className="col-xs-12">
						<h2>{movie.title} ({movie.year})</h2>
					</div>
					<div className="col-xs-12">
						<div className="row moviep-images">
							<div className="col-xs-12 col-sm-2">
								<img src={"/img/p" + movie.cover} />
							</div>
							<div className="col-xs-12 col-sm-10">
								<img src={"/img/b" + movie.backdrop} />
							</div>
						</div>
					</div>
					<div className="col-xs-12">
						<div className="row between-xs">
							<span className="col-xs-12 col-sm-6 director">{movie.director}</span>
							<span className="col-xs-12 col-sm-6 end-sm">{movie.production_countries}</span>
						</div>
					</div>
					<div className="col-xs-12">
						<div className="row between-xs">
							<span className="col-xs-12 col-sm-6 ">{movie.actors}</span>
							<span className="col-xs-12 col-sm-6 end-sm">{movie.genres}</span>
						</div>
					</div>
					<div className="col-xs-12 bottom-spacer">
						<div className="row between-xs">
							<div className="col-xs-12 col-sm-10">
								<span className="label">{movie.resolution}</span>
								<span className="label secondary spacer">{movie.location}</span>
							</div>
							<div className="col-xs-12 col-sm-2 end-sm">
								<span className="label"><i className="icon-plus"></i>&nbsp;{moment(movie.added).format('MMM DD, YYYY H:mm')}</span>
							</div>
						</div>
					</div>
					<div className="col-xs-12">
						<span>{movie.overview}</span>
					</div>
					<div className="col-xs-12 bottom-spacer-large">
						<div className="row between-xs">
							<div className="col-xs-12 col-sm-10 top-xs">
								<input type="text"></input>
								<button className="btn btn-default">Fix</button>
							</div>
							<div className="col-xs-12 col-sm-2 end-sm top-xs">
								<span className="label"><i className="icon-watched"></i></span>
							</div>
						</div>
					</div>										
				</article>				
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