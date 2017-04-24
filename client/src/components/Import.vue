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

<script>
import * as types from '../store/types'
import Console from './Console'

export default {
	name: 'import',

	components: { Console },

	methods: {
		runImport() {
			this.$store.dispatch(types.RUN_IMPORT)
		},
	},
}
</script>

<style lang="scss" scoped>

</style>
