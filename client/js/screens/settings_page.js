import React from 'react'
import { isNotValid } from '../lib/utils'

export default class SettingsPage extends React.Component {
	constructor() {
		super()

		this.addFolder = this.addFolder.bind(this)

		this.state = {
			folder: ""
		}
	}

	// componentWillMount() {
	// 	this.props.actions.settings.getConfig()
	// }

	render() {
		if (isNotValid(this.props.state.settings)) {
			return null
		}

		const settings = this.props.state.settings

		var noFolders;
		if (settings.mediaFolders.length === 0) {
			noFolders = (
				<div className="col-xs-12 bottom-spacer-half">
					<p className="bg-warning">There are no folders selected for importing. Please enter the media folders where you store your movies, to scan them</p>
				</div>				
			)
		}

		const folders = settings.mediaFolders.map(function(item, i) {
			return (
				<tr key={i}>
					<td><i className="icon-prune"></i></td>
					<td>{item}</td>
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
								<input className="addon-field" type="text" onKeyDown={this.addFolder}></input>
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
									{folders}
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

	addFolder(e) {
		if (e.key !== "Enter") {
			return
		}

		e.preventDefault()

		// console.log("settingsPage.addFolder: ", e.target.value)
		this.props.actions.settings.addMediaFolder(e.target.value)
	}	
}