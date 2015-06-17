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
	
		return (
			<section className="row">
				<a href="" onClick={handleImportMovies}>import</a>
			</section>
		)
	}
})

