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
					<li key={i}>
						<div className="cover-container" key={i}>
							<img src={"/img/p" + movie.cover} />
							<span className="crimson">{movie.title} </span><br />
							{movie.year} | 
							<span className="bright">{movie.imdb_rating}</span> |
							<span className="label">{movie.runtime}</span>
							<div className="ribbon" data-ng-if="movie.count_watched > 0">
								<span>watched</span>
							</div>	
						</div>
					</li>				
				)
			})

			return (
				<section className="container row covers">
					<ul className="small-block-grid-2 medium-block-grid-4 large-block-grid-6">
						{items}
					</ul>
				</section>
			)

		} else {
			return null
		}

	}
})

