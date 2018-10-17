<template>
	<section class="row c-flegal between-xs middle-xs">
		<ul class="col-xs-12 col-sm-4">
			<div>
				<span class="mv0 mr1">Copyright &copy;</span>
				<a href="https://jbrio.net/">Juan B. Rodriguez</a>
			</div>
		</ul>
		<ul class="col-xs-12 col-sm-4">
			<div class="center-xs">
				<span>mediaGUI &nbsp;</span>
				<span>v{{version}}</span>
			</div>
		</ul>
		<ul class="col-xs-12 col-sm-4 end-xs middle-xs">
			<li class="row end-xs middle-xs">
				<a class="end-xs middle-xs mv0 mr3"
					href="https://www.themoviedb.org/"
					title="themoviedb.org"
					height="40"
					target="_blank">
						<img src="@/assets/tmdb.png"
							alt="Logo for tmdb" />
				</a>

				<a class="end-xs middle-xs mv0 mr3"
					href="https://jbrio.net/"
					title="jbrio.net"
					target="_blank">
						<img src="@/assets/logo-small.png"
							alt="Logo for Juan B. Rodriguez" />
				</a>
			</li>
		</ul>
	</section>
</template>

<script lang="ts">
import Vue from 'vue'
import Component from 'vue-class-component'
import { State } from 'vuex-class'

@Component
export default class FooterBar extends Vue {
	@State(state => state.config.version)
	private version!: string
}
</script>

<style lang="scss" scoped>
@import '../styles/variables.scss';

.c-flegal {
	a,
	img,
	span,
	div {
		display: flex;
	}

	a {
		color: $menu-anchor;
	}
}
</style>
