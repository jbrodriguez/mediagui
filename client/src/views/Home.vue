<template>
  <Covers />
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Covers from '@/components/Covers.vue' // @ is an alias to /src

@Component({
	components: {
		Covers,
	},
})
export default class Home extends Vue {}
</script>