<template>
  <div class="container flex flex-column body">
    <header class="mv2">
      <header-bar></header-bar>
    </header>
    <main>
      <div class="flex flex-column">
        <router-view></router-view>
      </div>
    </main>
    <footer class="pv2 ph2 mv2 mh0">
      <footer-bar></footer-bar>
    </footer>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

import HeaderBar from "@/components/HeaderBar.vue"; // @ is an alias to /src
import FooterBar from "@/components/FooterBar.vue";

import * as constant from "@/constants";

@Component({
  components: {
    HeaderBar,
    FooterBar
  }
})
export default class App extends Vue {
  private mounted() {
    this.$store.dispatch(constant.FETCH_CONFIG);
  }
}
</script>

<style></style>
