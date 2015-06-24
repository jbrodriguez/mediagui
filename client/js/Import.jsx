const 	React 	= require('react'),
		movies 	= require('./movies')

const Console = React.createClass({
	componentDidUpdate: function() {
		var node = this.getDOMNode();
		node.scrollTop = node.scrollHeight;
	},	

	render: function() {
		const items = this.props.messages.map(function(message, i) {
			return (
				<p key={i} className="col-xs-12 console__line">{message.payload}</p>
			)
		})

		return (
			<div className="console">
				<div className="row">
					<div className="col-xs-12 console__lines">
						{items}
					</div>
				</div>						
			</div>
		)		
	}
})


module.exports = React.createClass({
	handleImportMovies: function() {
		movies.importMovies()
	},

	render: function() {
		return (
			<section className="row">
				<div className="col-xs-12 bottom-spacer-half">
					<button className="btn btn-default" onClick={this.handleImportMovies}>Import</button>
				</div>
				<div className="col-xs-12">
					<Console messages={this.props.messages} />
				</div>
			</section>
		)
	}
})

