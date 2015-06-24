const 	React 	= require('react'),
		Cover 	= require('./Cover.jsx')

module.exports = React.createClass({
	render: function() {
		const movies = this.props.movies.items

		console.log('movies: ', movies)

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
})

