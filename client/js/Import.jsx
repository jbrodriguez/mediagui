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
				<p key={i} className="console__line">message</p>
			)
		})
	
		return (
			<section className="row">
				<div className="col-xs-12">
					<button className="btn btn-default" onClick={handleImportMovies}>import</button>
				</div>
				<div className="col-xs-12">
					<div>
						<div className="row console" data-ng-if="vm.showConsole">
							<div className="col-xs-12 console__lines" data-unb-scroll-bottom="vm.lines">
								<p ng-repeat="line in vm.lines track by $index" class="console__line">{{line}}</p>
							</div>
						</div>						
					</div>
				</div>
			</section>
		)
	}
})

