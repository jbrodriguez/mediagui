<template>
	<div class="c-mcover">
		<div class="row">
			<div class="col-xs-12">
				<div class="c-moverlay">
					<img :src="'/img/p' + movie.cover" />
					<div v-if="watched"
					     class="c-moverlaycover">
						<span>watched</span>
					</div>
				</div>
			</div>
		</div>
		<div class="row">
			<div class="col-xs-12">
				<div class="pv0 ph2">
					<p class="c-mctitle">{{movie.title}}</p>
					<div class="between-xs c-mcdetails">
						<span>{{movie.year}}</span>
						<span>{{movie.imdb_rating}}</span>
						<span>{{runtime}}</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script lang="ts">
import Vue from 'vue'
import Component from 'vue-class-component'
import { State } from 'vuex-class'
import { Prop } from 'vue-property-decorator'

import { Movie } from '@/types'
import * as constant from '@/constants'
import { hourMinute } from '@/lib/utils'

@Component
export default class Cover extends Vue {
	// @State(state => state.movies)
	// private movies!: Movie[]
	@Prop(Object)
	private movie!: Movie

	get watched() {
		return this.movie.count_watched > 0
	}

	get runtime(): string {
		return hourMinute(this.movie.runtime)
	}
}
</script>

<style lang="scss" scoped>
@import '../styles/variables.scss';

.c-mcover {
	background-color: $movie-cover-bg;
}

.c-moverlay {
	position: relative;
	overflow: hidden;
}

.c-moverlaycover {
	position: absolute;
	left: -40px;
	top: 10px;
	background-color: rgba(170, 0, 0, 0.6);
	transform: rotate(-45deg);
	box-shadow: 0 0 10px #888;

	span {
		border: 1px solid #faa;
		color: #fff;
		display: block;
		font-size: 0.8em;
		margin: 1px 0;
		padding: 5px 50px;
		text-align: center;
		text-decoration: none;
		/* shadow */
		text-shadow: 0 0 5px #444;
	}
}

.c-mctitle {
	color: $movie-title-color;
	white-space: nowrap;
	width: 100%;
	overflow: hidden;
	/* "overflow" value must be different from "visible" */
	text-overflow: ellipsis;
}

.c-mcdetails {
	display: flex;
	font-size: 0.75em;
}
</style>
