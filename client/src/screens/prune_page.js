import React from 'react'
import Console from './console_card'

export default class PrunePage extends React.Component {
	constructor() {
		super()

		this.handlePruneMovies = this.handlePruneMovies.bind(this)
	}

	render() {
		return (
			<section className="row">
				<div className="col-xs-12 bottom-spacer-half">
					<button className="btn btn-default" onClick={this.handlePruneMovies}>Prune</button>
				</div>
				<div className="col-xs-12">
					<Console messages={this.props.state.messages} />
				</div>
			</section>
		)
	}

	handlePruneMovies() {
		this.props.actions.movies.pruneMovies()
	}
}