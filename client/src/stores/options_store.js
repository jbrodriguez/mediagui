import Bacon from 'baconjs';
import R from 'ramda';
import ffux from 'ffux';

const Options = ffux.createStore({
	actions: ["setFilterBy", "setSortBy", "setSortOrder", "setOffset", "setQueryTerm"],

	state: (initialOptions, {setFilterBy, setSortBy, setSortOrder, setOffset, setQueryTerm}, {storage}) => {
		const dQueryTerm = setQueryTerm.debounce(750)

		return Bacon.update(
			initialOptions,
			setFilterBy, _setFilterBy,
			setSortBy, _setSortBy,
			setSortOrder, _setSortOrder,
			setOffset, _setOffset,
			dQueryTerm, _setQueryTerm
		)

		function _setFilterBy(options, filterBy) {
			storage.set('filterBy', filterBy)
			return R.merge(options, {filterBy: filterBy})
		}

		function _setSortBy(options, sortBy) {
			storage.set('sortBy', sortBy)
			return R.merge(options, {sortBy: sortBy})
		}

		function _setSortOrder(options, sortOrder) {
			storage.set('sortOrder', sortOrder)
			return R.merge(options, {sortOrder: sortOrder})
		}

		function _setOffset(options, offset) {
			return R.merge(options, {offset: offset})
		}

		function _setQueryTerm(options, term) {
			return R.merge(options, {query: term, offset: 0})
		}
	}
})

module.exports = Options