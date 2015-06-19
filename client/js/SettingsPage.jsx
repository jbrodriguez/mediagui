const 	React 	= require('react')

module.exports = React.createClass({
	// componentWillMount: function() {
	// 	movies.getCover()
	// },
	render: function() {
		var noFolders;

		const settings = this.props.settings

		if (settings.mediaFolders.length === 0) {
			noFolders = (
				<div className="col-xs-12 bottom-spacer-half">
					<p className="bg-warning">There are no folders selected for importing. Please enter the media folders where you store your movies, to scan them</p>
				</div>				
			)
		}
	
		return (
			<section className="row">
				{noFolders}

				<div className="col-xs-12 bottom-spacer-half">
					<form>
					<fieldset>
						<legend>
							Where are your movies stored ?
						</legend>
						<div className="row bottom-spacer-large">
							<div className="col-xs-12 addon">
								<span className="addon-item">Folder</span>
								<input className="addon-field" type="text"></input>
								<button className="btn btn-default">Add</button>
							</div>
						</div>
						<div className="row bottom-spacer-large">
							<div className="col-xs-12">
								<table>
								<thead>
									<th width="50">#</th>
									<th>Folder</th>
								</thead>
								<tbody>
									<tr data-ng-repeat="folder in vm.options.config.mediaFolders">
									<td><i class="fa fa-times"></i></td>
									<td>/Volumes/hal-films</td>
									</tr>
									<tr data-ng-repeat="folder in vm.options.config.mediaFolders">
									<td><i class="fa fa-times"></i></td>
									<td>/Volumes/hal-films</td>
									</tr>		
								</tbody>
								</table>
							</div>
						</div>
					</fieldset>
					</form>
				</div>

			</section>
		)
	}
})

