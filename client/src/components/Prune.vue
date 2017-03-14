<template>
  <div class="">
    <section class="row">
      <div class="col-xs-12 mb2">
        <button class="btn btn-default" @click="onPrune">PRUNE</button>
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
import * as types from '../store/types';
import Console from './Console';

export default {
  name: 'prune',

  components: { Console },

  methods: {
    onPrune() {
      this.$store.dispatch(types.RUN_PRUNE);
    },
  },
};
</script>

<style lang="scss" scoped>
</style>
