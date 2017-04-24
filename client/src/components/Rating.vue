<template>
	<div class="rating">
		<i v-for="n in max"
		   class="fa mr1 c-star"
		   :class="[overwritten_value >= n ? 'fa-star' : 'fa-star-o']"
		   v-on:mouseover="star_over(n)"
		   v-on:mouseout="star_out"
		   v-on:click.stop.prevent="star_selected(n)"></i>
	</div>
</template>
<script>
export default {
	name: 'rating',

	data() {
		return {
			is_selected: true,
			overwritten_value: this.value,
		}
	},

	props: {
		max: Number,
		value: Number,
		disabled: Boolean,
	},

	watch: {
		overwritten_value() {
			this.selected_value = false
		},
	},

	methods: {
		star_over(index) {
			if (!this.disabled) {
				this.temp_value = this.overwritten_value
				this.overwritten_value = index
			}
		},

		star_out() {
			if (!this.disabled && !this.selected_value) {
				this.overwritten_value = this.temp_value
			}
		},

		star_selected(value) {
			if (!this.disabled) {
				this.selected_value = true
				this.overwritten_value = value
				this.$emit('rating-selected', this.overwritten_value)
			}
		},
	},
}
</script>

<style lang="scss" scoped>
.c-star {
	color: white;
	text-shadow: 0 1px 0 hsla(0, 0%, 0%, .75), 0 0 1px hsla(0, 0%, 0%, .75), 0 1px 5px hsla(0, 0%, 0%, .75);
}
</style>
