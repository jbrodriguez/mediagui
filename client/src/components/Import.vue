<template>
	<div class="">
		<section class="row">
			<div class="col-xs-12 mb2">
				<button class="btn btn-default"
				        @click="runImport">Import</button>
			</div>
		</section>
		<section class="row">
			<div class="col-xs-12">
				<console></console>
			</div>
		</section>
	</div>
</template>

<script lang="ts">
import Vue from 'vue'
import Component from 'vue-class-component'

import * as constant from '@/constants'

import Console from './Console.vue'

@Component({
	components: {
		Console,
	},
})
export default class Import extends Vue {
	public runImport() {
		this.$store.dispatch(constant.RUN_IMPORT)
	}
}
</script>

<style lang="scss" scoped>
</style>
