<template>
	<div class="c-console">
		<div class="row">
			<div class="col-xs-12">
				<p v-for="(line, index) in lines" :key="index" class="console__line">{{line}}</p>
			</div>
		</div>
	</div>
</template>

<script lang="ts">
import Vue from 'vue'
import Component from 'vue-class-component'
import { State } from 'vuex-class'

import * as constant from '@/constants'

@Component
export default class Console extends Vue {
	@State(state => state.lines)
	private lines!: string[]

	public runImport() {
		this.$store.dispatch(constant.RUN_IMPORT)
	}
}
</script>

<style lang="scss" scoped>
@import '../styles/variables.scss';

.c-console {
	display: flex;
	overflow: auto; // width: 100%;
	height: 750px; // height: 100%;
	// border: 1px solid $console-border;
	padding: 0.75em;
	color: $console-color;
	background-color: $console-bg; // flex: 1;
	.console__line {
		font-size: 0.8em; // line-height: 0.5em;
	}

	p {
		margin-bottom: 0;
		white-space: nowrap;
	}
}
</style>
