<template>
	<nav class="row between-xs">
		<ul class="col-xs-12 col-sm-2 center-xs">
			<li class="pv2 ph0 c-hlogo">
				<router-link to="/">mediaGUI</router-link>
			</li>
		</ul>
		<ul class="col-xs-12 col-sm-10">
			<li class="pv2 ph0 c-hmenu">
				<div class="row between-xs">
					<div class="col-xs-12 col-sm-8">
						<div class="c-hmenusection">
							<router-link to="/movies" class="mv0 mh2">MOVIES</router-link>

							<select :value="filterBy" @change="changeFilter">
								<option v-for="option in options.filterByOptions" :value="option.value" :key="option.value">{{option.label}}</option>
							</select>

							<input type="search" placeholder="Enter search string" :value="options.query" @input="updateQuery">

							<select :value="sortBy" @change="changeSort">
								<option v-for="sort in options.sortByOptions" :value="sort.value" :key="sort.value">{{sort.label}}</option>
							</select>

							<font-awesome-icon class="fa mv0 ml2 c-fai" :icon="chevron" @click="changeOrder" />

							<span class="mv0 mh2">|</span>

							<router-link to="/import" class="mv0">IMPORT</router-link>

							<span class="mv0 ml2"></span>

							<router-link to="/add" class="mv0">ADD</router-link>
						</div>
					</div>
					<div class="col-xs-12 col-sm-4 end-xs">
						<div class="c-hmenusection">
							<router-link to="/settings">SETTINGS</router-link>

							<span class="mv0 mh2">|</span>

							<router-link to="/duplicates">DUPLICATES</router-link>

							<router-link to="/prune" class="mv0 mh2">PRUNE</router-link>
						</div>
					</div>
				</div>
			</li>
		</ul>
	</nav>
</template>

<script lang="ts">
import debounce from 'lodash.debounce'
import { Component, Prop, Vue } from 'vue-property-decorator'
import { State } from 'vuex-class'
// import format from 'date-fns/format'

import * as constant from '@/constants'
import { OptionsState } from '@/types'
// import { Statement, Category } from '@/types'

@Component
export default class HeaderBar extends Vue {
	@State(state => state.options.filterBy)
	private filterBy!: string

	@State(state => state.options.sortBy)
	private sortBy!: string

	@State(state => state.options)
	private options!: OptionsState

	get chevron() {
		return this.options.sortOrder === 'asc' ? 'chevron-circle-up' : 'chevron-circle-down'
	}

	// get selected() {
	// 	return this.page
	// }

	// set selected(value: number) {
	// 	console.log(`selected(${selected})`)
	// }

	private updateQuery = debounce((e: Event) => {
		this.$store.commit(constant.SET_QUERY, (e.target as HTMLInputElement).value)
		this.$store.dispatch(constant.FETCH_MOVIES)
	}, 750)

	private changeFilter(e: Event) {
		this.$store.commit(constant.SET_FILTER, (e.target as HTMLInputElement).value)
		this.$store.dispatch(constant.FETCH_MOVIES)
	}

	private changeSort(e: Event) {
		this.$store.commit(constant.SET_SORT, (e.target as HTMLInputElement).value)
		this.$store.dispatch(constant.FETCH_MOVIES)
	}

	private changeOrder() {
		this.$store.commit(constant.FLIP_ORDER)
		this.$store.dispatch(constant.FETCH_MOVIES)
	}
}
</script>

<style lang="scss" scoped>
@import '../styles/variables.scss';
// @import '../styles/custom.scss';

.c-hlogo {
	background-color: $logo-background;

	a {
		color: $logo-anchor;
	}
}

.c-hmenu {
	background-color: $menu-background;

	a {
		font-size: 0.9em;
		color: $menu-anchor;
	}
}

.c-hmenusection {
	// background-color: $headerMenuBackground;
	a {
		font-size: 0.9em; // color: $headerMenuAnchorColor;
		&.router-link-exact-active {
			padding-bottom: 3px;
			border-bottom: 2px solid $menu-anchor;
		}
	}
}
</style>
