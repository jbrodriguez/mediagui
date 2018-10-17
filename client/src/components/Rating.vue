<template>
	<div class="rating">
		<font-awesome-icon v-for="(n, index) in max" class="mr1 c-star" 
			:key="index" 
			:icon="star(n)" 
			v-on:mouseover="starOver(n)"
			v-on:mouseout="starOut"
			v-on:click.stop.prevent="starSelected(n)" />
	</div>
</template>

<script lang="ts">
import Vue from 'vue'
import Component from 'vue-class-component'
import { State } from 'vuex-class'

import { OptionsState, DomainState, RootState } from '@/types'
import * as constant from '@/constants'
import { Prop, Watch } from 'vue-property-decorator'

@Component
export default class Rating extends Vue {
	@Prop()
	private max!: number

	@Prop()
	private value!: number

	@Prop()
	private disabled!: boolean

	private isSelected: boolean = true
	private overwrittenValue: number = this.value
	private selectedValue: boolean = false
	private tempValue: number = this.overwrittenValue

	@State((state: RootState) => state.options)
	private options!: OptionsState

	@State((state: RootState) => state.domain)
	private domain!: DomainState

	@Watch('overwrittenValue')
	private onOverwrittenValueChanged() {
		this.selectedValue = false
	}

	// private mounted() {
	// 	this.overwrittenValue = this.value
	// }

	private star(n: number) {
		// console.log(`value(${n})-over(${this.overwrittenValue})-value(${this.value})`)
		return this.overwrittenValue >= n ? 'star' : ['far', 'star']
		// return this.options.sortOrder === 'asc' ? 'chevron-circle-up' : 'chevron-circle-down'
	}

	private starOver(index: number) {
		if (!this.disabled) {
			this.tempValue = this.overwrittenValue
			this.overwrittenValue = index
		}
	}

	private starOut() {
		if (!this.disabled && !this.selectedValue) {
			this.overwrittenValue = this.tempValue
		}
	}

	private starSelected(value: number) {
		if (!this.disabled) {
			this.selectedValue = true
			this.overwrittenValue = value
			this.$emit('rating-selected', this.overwrittenValue)
		}
	}
}

// export default {
// 	name: 'rating',

// 	data() {
// 		return {
// 			is_selected: true,
// 			overwritten_value: this.value,
// 		}
// 	},

// 	props: {
// 		max: Number,
// 		value: Number,
// 		disabled: Boolean,
// 	},

// 	watch: {
// 		overwritten_value() {
// 			this.selected_value = false
// 		},
// 	},

// 	methods: {
// 		star_over(index) {
// 			if (!this.disabled) {
// 				this.temp_value = this.overwritten_value
// 				this.overwritten_value = index
// 			}
// 		},

// 		star_out() {
// 			if (!this.disabled && !this.selected_value) {
// 				this.overwritten_value = this.temp_value
// 			}
// 		},

// 		star_selected(value) {
// 			if (!this.disabled) {
// 				this.selected_value = true
// 				this.overwritten_value = value
// 				this.$emit('rating-selected', this.overwritten_value)
// 			}
// 		},
// 	},
// }
</script>

<style lang="scss" scoped>
.c-star {
	color: white;
	text-shadow: 0 3px 0 hsla(0, 0%, 0%, 0.75), 0 0 3px hsla(0, 0%, 0%, 0.75), 0 3px 5px hsla(0, 0%, 0%, 0.75);
	width: 14px;
	height: 14px;
}
</style>
