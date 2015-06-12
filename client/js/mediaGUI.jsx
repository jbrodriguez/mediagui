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
	componentWillMount: function() {
		console.log('reminiscing')
		console.log(Array.isArray(this.props.children)); // => true
	},

	// getInitialState: function() {
	// 	return {
	// 		selectedFilter: this.props.options.filterBy,
	// 		selectedSort: this.props.options.sortBy
	// 	}
	// },

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

        const handleSortBy = function() {
			const sortBy = event.target.value

        	console.log("mediaGUI.jsx.handleSortBy:" + sortBy)

			// this.setState({selectedSort: sortBy})
			optionsBO.setSortBy(sortBy)
        }

        const handleSortOrder = function() {
        	const sortOrder = options.sortOrder === 'asc' ? 'desc' : 'asc'

        	console.log("mediaGUI.jsx.handleSortOrder:" + sortOrder)

			// this.setState({selectedSort: sortBy})
			optionsBO.setSortOrder(sortOrder)
        }

        const handleQueryTerm = function() {
			const queryTerm = event.target.value

        	console.log("mediaGUI.jsx.handleQueryTerm:" + queryTerm)

			optionsBO.setQueryTerm(queryTerm)
        }


        const chevron = cx({
        	'icon-chevron-down': options.sortOrder === 'desc',
        	'icon-chevron-up': options.sortOrder === 'asc',
        	// 'header__action': true
        })
        // const sortStyle = {marginLeft: "1em"}

								// <select value={this.state.selectedFilter}>
								// <select value={this.state.selectedSort}>
										// <a href="#" className="spacer">{"prune".toUpperCase()}</a>

		return (
			// <div className={cx("container", "body")}>
			<div className="container">
				<header>
					<nav className="row between-xs">
						<ul className="col-xs-12 col-sm-2 center-xs">
							<li className="header__logo">
								<Link to="cover">mediaGUI</Link>
							</li>
						</ul>
						<ul className="col-xs-12 col-sm-10 center-xs">
							<li className="row between-xs">
								<div className="col-xs-12 col-sm-8">
									<div className="header__menu">
										<Link to="movies" query={urlQuery} className="spacer">MOVIES</Link>

										<select value={options.filterBy}>
											{filterByNodes}
										</select>
										<input type="search" placeholder="Enter search string" onChange={handleQueryTerm} />

										<select value={options.sortBy} onChange={handleSortBy} className="spacer">
											{sortByNodes}
										</select>

										<i onClick={handleSortOrder} className={chevron}></i>

										<span className="spacer">|</span>

										<a href="#">{"import".toUpperCase()}</a>
									</div>
								</div>
								<div className="col-xs-12 col-sm-4">
									<div className="header__menu">
										<a href="#">{"settings".toUpperCase()}</a>
										<span className="spacer">|</span>
										<a href="#">{"duplicates".toUpperCase()}</a>
										<a href="#" className="spacer">{"prune".toUpperCase()}</a>
									</div>
								</div>
							</li>
						</ul>
					</nav>
				</header>

				<main>
					<RouteHandler { ...this.props}/>
				</main>

				<footer className="row">
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