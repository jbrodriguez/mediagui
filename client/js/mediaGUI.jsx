const 	React 			= require('react'),
		cx 		 		= require('classnames'),
		Router 			= require('react-router'),
		Link			= Router.Link,
		RouteHandler 	= Router.RouteHandler,
		optionsBO 		= require('./options.js')

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
    handleFilterBy: function() {
		const filterBy = event.target.value
		optionsBO.setFilterBy(filterBy)
    },

    handleSortBy: function() {
		const sortBy = event.target.value
		optionsBO.setSortBy(sortBy)
    },

    handleSortOrder: function() {
    	const sortOrder = this.props.options.sortOrder === 'asc' ? 'desc' : 'asc'
		optionsBO.setSortOrder(sortOrder)
    },

    handleQueryTerm: function() {
		const queryTerm = event.target.value
		optionsBO.setQueryTerm(queryTerm)
    },

	render: function() {
		// console.log('somebody to love: ' + JSON.stringify(this.props, null, 4))
		const settings = this.props.settings
		const options = this.props.options
		const urlQuery = {
			query: options.query,
			filterBy: options.filterBy,
			sortBy: options.sortBy,
			sortOrder: options.sortOrder,
			limit: options.limit,
			offset: options.offset
		}

        var filterByNodes = options.filterByOptions.map(function(option){
            return <FilterWrapper key={option.id} option={option} />
        })

        var sortByNodes = options.sortByOptions.map(function(option){
            return <SortWrapper key={option.id} option={option} />
        })

        const chevron = cx({
        	'icon-chevron-down': options.sortOrder === 'desc',
        	'icon-chevron-up': options.sortOrder === 'asc',
        	// 'header__action': true
        })

		return (
			// <div className={cx("container", "body")}>
			<div className="container body">
				<header>
					<nav className="row between-xs">
						<ul className="col-xs-12 col-sm-2 center-xs">
							<li className="header__logo">
								<Link to="cover">mediaGUI</Link>
							</li>
						</ul>
						<ul className="col-xs-12 col-sm-10">
							<li className="header__menu">
								<div className="row between-xs">
									<div className="col-xs-12 col-sm-8">
										<div className="header__menu-section">
											<Link to="movies" query={urlQuery} className="spacer">MOVIES</Link>

											<select value={options.filterBy} onChange={this.handleFilterBy}>
												{filterByNodes}
											</select>
											<input type="search" placeholder="Enter search string" onChange={this.handleQueryTerm} />

											<select value={options.sortBy} onChange={this.handleSortBy} className="spacer">
												{sortByNodes}
											</select>

											<i onClick={this.handleSortOrder} className={chevron}></i>

											<span className="spacer">|</span>

											<Link to="import">{"import".toUpperCase()}</Link>
										</div>
									</div>
									<div className="col-xs-12 col-sm-4 end-xs">
										<div className="header__menu-section">
											<Link to="settings">{"settings".toUpperCase()}</Link>
											<span className="spacer">|</span>
											<Link to="duplicates">{"duplicates".toUpperCase()}</Link>
											<Link to="prune" className="spacer">{"prune".toUpperCase()}</Link>
										</div>
									</div>
								</div>
							</li>
						</ul>
					</nav>
				</header>

				<main>
					<RouteHandler { ...this.props}/>
				</main>

				<footer>
				    <section className="row legal between-xs middle-xs">
				    	<ul className="col-xs-12 col-sm-4">
				    		<div>
						        <span className="copyright spacer">Copyright &copy; 2015</span>
						        <a href='http://jbrodriguez.io/'>Juan B. Rodriguez</a>
					       	</div>
				    	</ul>
				    	<ul className="col-xs-12 col-sm-4">
				    		<div className="center-xs">
						        <span className="version">mediaGUI &nbsp;</span>
						        <span className="version">v{settings.version}</span>
					        </div>
				    	</ul>
				    	<ul className="col-xs-12 col-sm-4 end-xs middle-xs">
							<a className="end-xs middle-xs spacer" href="http://jbrodriguez.io/" title="jbrodriguez.io" target="_blank">
								<img src="app/logo-small.png" alt="Logo for Juan B. Rodriguez" />
							</a>
				    	</ul>
				    </section>				
				</footer>
			</div>
		)
	}
})