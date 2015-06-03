const 	React 			= require('react'),
		Router 			= require('react-router'),
		Link			= Router.Link,
		RouteHandler 	= Router.RouteHandler

module.exports = React.createClass({
	componentWillMount: function() {
		console.log('reminiscing')
		console.log(Array.isArray(this.props.children)); // => true
	},

	render: function() {
		console.log('somebody to love')
		const settings = this.props.settings
		return (
			<div className="body">
				<header>
					<nav className="row">
						<ul>
							<Link to="app">Home</Link>
						</ul>
						<ul>
							<Link to="movies">Movies</Link>
							<select data-ng-model="home.options.filterBy" data-ng-options="option.value as option.display for option in home.options.filterByOptions">
							</select>
							<input type="search" placeholder="Enter search string" data-ng-model="home.options.searchTerm" ng-model-options="{ debounce: 750 }" />
							<select data-ng-model="home.options.sortBy" data-ng-options="option.value as option.display for option in home.options.sortByOptions">
							</select>
							<a href="#" data-ng-click="home.sortOrder()"><i class="fa" data-ng-class="home.options.sortOrder === 'desc' ? 'fa-chevron-circle-down' : 'fa-chevron-circle-up'"></i></a>
						</ul>
					</nav>
				</header>

				<main>
					<RouteHandler/>
				</main>

				<footer>
				    <section className="legal row">
				        <span className="copyright">Copyright &copy; 2015 &nbsp; <a href='http://jbrodriguez.io/'>Juan B. Rodriguez</a></span>
				        <span className="version">mediaGUI v{settings.version}</span>
				        <div className="logos">
					        <span><a href="http://jbrodriguez.io/" title="jbrodriguez.io" target="_blank"><img src="app/logo-small.png" alt="Logo for Juan B. Rodriguez" /></a></span>
				       </div>
				    </section>				
				</footer>
			</div>
		)
	}
})