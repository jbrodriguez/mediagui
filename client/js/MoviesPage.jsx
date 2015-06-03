const 	React 			= require('react'),
		Link			= require('react-router').Link,
		RouteHandler 	= require('react-router').RouteHandler

module.exports = React.createClass({
	render: function() {
		const movies = this.props.movies

		var items = movies.map(function(movie) {
			return (
				<div>{movie.title (movie.year)}</div>
			)
		})

		return (
			{items}
		)
	}
})