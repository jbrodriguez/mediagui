const 	React 	= require('react')

module.exports = React.createClass({
	// componentWillMount: function() {
	// 	movies.getCover()
	// },

	render: function() {
		// console.log('MoviesCover.jsx: ' + JSON.stringify(this.props, null, 4))

		const movies = this.props.movies.items

		console.log('movies: ' + movies)

		if (typeof movies != 'undefined') {
			var items = movies.map(function(movie, i) {
				return (
					<div key={i} className="col-xs-12 col-sm-6 col-md-3 col-lg-2">
						<div key={i} className="cover-container" >
							<img src={"/img/p" + movie.cover} />
							<span className="crimson">{movie.title} </span><br />
							{movie.year} | 
							<span className="bright">{movie.imdb_rating}</span> |
							<span className="label">{movie.runtime}</span>
							<div className="ribbon" data-ng-if="movie.count_watched > 0">
								<span>watched</span>
							</div>	
						</div>
					</div>				
				)
			})

			return (
				<section className="row">
					{items}
				</section>
			)

		} else {
			return null
		}

	}
})

