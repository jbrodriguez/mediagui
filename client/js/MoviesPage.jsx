const 	React 			= require('react'),
		Link			= require('react-router').Link,
		RouteHandler 	= require('react-router').RouteHandler

module.exports = React.createClass({
	render: function() {
		const movies = this.props.movies.items

		const styles = {height: "17em"}
		const styleo = {overflow: "hidden", maxHeight: "17em"}

		var items = movies.map(function(movie, i) {
			return (
				<div key={i}>
					<div className="col-xs-12">
						<h2>{movie.title} ({movie.year})</h2>
					</div>
					<div className="col-xs-12">
						<div className="row" style={styleo}>
							<div className="col-xs-12 col-lg-2">
								<img src={"/img/p" + movie.cover} />
							</div>
							<div className="col-xs-12 col-lg-10">
								<img src={"/img/b" + movie.backdrop} />
							</div>
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
	}
})