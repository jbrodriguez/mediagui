import React from 'react'

export default class Console extends React.Component {
	componentDidUpdate() {
		var node = this.getDOMNode();
		node.scrollTop = node.scrollHeight;
	}

	render() {
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
}
