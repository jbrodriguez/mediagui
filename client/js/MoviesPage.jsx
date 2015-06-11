const 	React 			= require('react'),
		Link			= require('react-router').Link,
		RouteHandler 	= require('react-router').RouteHandler,
		pager 			= require('react-paginate')

module.exports = React.createClass({
	render: function() {
		const movies = this.props.movies.items
		const options = this.props.options

		var pagination;

		if (movies.total > options.limit) {
			pagination = 
		}

		// const styles = {height: "17em"}
		// const styleo = {overflow: "hidden", maxHeight: "17em"}

		var items = movies.map(function(movie, i) {
			return (
				<article key={i} className="moviep">
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
				</article>				
			)
		})

		return (
			<section className="row">
				{pagination}

				{items}

				{pagination}
			</section>
		)
	}
})