<template>
	<div class="container body">
		<header class="mv2 mh0">
			<header-bar></header-bar>
		</header>
		<main>
			<router-view></router-view>
		</main>
		<footer class="pv2 ph0 mv2 mh0">
			<footer-bar></footer-bar>
		</footer>
	</div>
</template>

<script>
import HeaderBar from './components/HeaderBar.vue'
import FooterBar from './components/FooterBar.vue'
import * as types from './store/types'

export default {
	name: 'app',

	components: { HeaderBar, FooterBar },

	created() {
		this.$store.dispatch(types.FETCH_CONFIG)
	},
}
</script>

<style lang="scss" scoped>
@import "./styles/_settings.scss";

.body {
	display: flex;
	flex-direction: column;
	min-height: 100vh;

	box-sizing: border-box;
	max-width: 78rem;
	margin: 0 auto;
}

main {
	flex: 1;
}

footer {
	color: $footer-color;
	background-color: $footer-bg;
}

.c-hlogo {
	background-color: $headerLogoBackground;

	a {
		color: $headerLogoAnchorColor;
	}
}

.c-hmenu {
	background-color: $headerMenuBackground;
}
</style>
