const React = require('react')

module.exports = React.createClass({
	componentWillMount: function() {
		console.log('reminiscing')
		console.log(Array.isArray(this.props.children)); // => true
	},

	render: function() {
		console.log('somebody to love')
		const settings = this.props.settings
		return (
	    	<h1>{settings.version}</h1>
		)
	}
})