const 	React 	= require('react'),
		movies 	= require('./movies')

module.exports = React.createClass({
	// componentWillMount: function() {
	// 	movies.getCover()
	// },
	render: function() {
		const handleImportMovies = function() {
			movies.importMovies()
		}

		const items = this.props.messages.map(function(message, i) {
			return (
				<p key={i} className="col-xs-12 console__line">{message.payload}</p>
			)
		})

		console.log('import.jsx.props.messages', this.props.messages);
	
		return (
			<section className="row">
				<div className="col-xs-12 bottom-spacer-half">
					<button className="btn btn-default" onClick={handleImportMovies}>Import</button>
				</div>
				<div className="col-xs-12">
					<div className="console">
						<div className="row" data-ng-if="vm.showConsole">
							<div className="col-xs-12 console__lines" data-unb-scroll-bottom="vm.lines">
								{items}
							</div>
						</div>						
					</div>
				</div>
			</section>
		)
	}
})

