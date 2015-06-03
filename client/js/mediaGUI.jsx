const 	React 			= require('react'),
		Link			= require('react-router').Link,
		RouteHandler 	= require('react-router').RouteHandler

module.exports = React.createClass({
	componentWillMount: function() {
		console.log('reminiscing')
		console.log(Array.isArray(this.props.children)); // => true
	},

	render: function() {
		console.log('somebody to love')
		const settings = this.props.settings
		return (
			<div className="nav">
				<Link to="app">Home</Link>
				<Link to="movies">Movies</Link>
				<RouteHandler />
			</div>
		)
	}
})