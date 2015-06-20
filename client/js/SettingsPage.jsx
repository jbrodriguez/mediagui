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

		const folders = settings.mediaFolders.map(function(folder, i) {
			return (
				<tr key={i}>
					<td><i class="fa fa-times"></i></td>
					<td>{folder}</td>
				</tr>
			)
		})
	
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
								{folders}
							</div>
						</div>
					</fieldset>
					</form>
				</div>

			</section>
		)
	}
})

