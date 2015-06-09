const 	React 			= require('react'),
		Router 			= require('react-router'),
		Link			= Router.Link,
		RouteHandler 	= Router.RouteHandler

var FilterWrapper = React.createClass({
  render: function() {
    return <option value={this.props.option.value}>{this.props.option.label}</option>
  }
})

var SortWrapper = React.createClass({
  render: function() {
    return <option value={this.props.option.value}>{this.props.option.label}</option>
  }
})

module.exports = React.createClass({
	componentWillMount: function() {
		console.log('reminiscing')
		console.log(Array.isArray(this.props.children)); // => true
	},

	getInitialState: function() {
		return {
			selectedFilter: this.props.options.filterBy,
			selectedSort: this.props.options.sortBy
		}
	},

	render: function() {
		// console.log('somebody to love: ' + JSON.stringify(this.props, null, 4))
		const settings = this.props.settings
		const options = this.props.options

        var filterByNodes = options.filterByOptions.map(function(option){
            return <FilterWrapper key={option.id} option={option} />
        })

        var sortByNodes = options.sortByOptions.map(function(option){
            return <SortWrapper key={option.id} option={option} />
        })

		return (
			<div className="body">
				<header>
					<nav className="row">
						<ul>
							<Link to="cover">Home</Link>
						</ul>
						<ul>
							<Link to="movies">Movies</Link>
							<select value={this.state.selectedFilter}>
								{filterByNodes}
							</select>
							<input type="search" placeholder="Enter search string" data-ng-model="home.options.searchTerm" ng-model-options="{ debounce: 750 }" />
							<select value={this.state.selectedSort}>
								{sortByNodes}
							</select>
							<a href="#" data-ng-click="home.sortOrder()"><i className="fa" data-ng-class="home.options.sortOrder === 'desc' ? 'fa-chevron-circle-down' : 'fa-chevron-circle-up'"></i></a>
						</ul>
					</nav>
				</header>

				<main>
					<RouteHandler { ...this.props}/>
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