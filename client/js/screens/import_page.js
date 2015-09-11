import React from 'react'
import Console from './console_card'

export default class ImportPage extends React.Component {
	constructor() {
		super()

		this.handleImportMovies = this.handleImportMovies.bind(this)
	}

	render() {
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

	handleImportMovies() {
		this.props.actions.movies.importMovies()
	}
}