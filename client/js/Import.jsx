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

		// if tyepof this.props.messages != 'undefined' {
		// 	const items = this.props.messages.map(function(message, i) {
		// 		return (
		// 			<p key={i} className="console__line">message</p>
		// 		)
		// 	})
		// }
	
		return (
			<section className="row">
				<div className="col-xs-12 bottom-spacer-half">
					<button className="btn btn-default" onClick={handleImportMovies}>Import</button>
				</div>
				<div className="col-xs-12">
					<div className="console">
						<div className="row" data-ng-if="vm.showConsole">
							<p className="col-xs-12 console__line" key="1">Never Surrender</p>
							<p className="col-xs-12 console__line" key="2">Dont dream its over</p>
						</div>						
					</div>
				</div>
			</section>
		)
	}
})

							// <div className="col-xs-12 console__lines" data-unb-scroll-bottom="vm.lines">
							// 	{items}
							// </div>