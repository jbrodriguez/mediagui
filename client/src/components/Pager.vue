<template>
  <ul :class="containerClass">
    <li :class="prevClass">
      <a @click="prevPage()" :class="prevLinkClass" tabindex="0">{{ prevText }}</a>
    </li>
    <li v-for="page in pages" :class="[{ active: page.selected, disabled: page.disabled }, pageClass]">
      <a v-if="page.disabled" :class="pageLinkClass" tabindex="0">{{ page.content }}</a>
      <a v-else @click="handlePageSelected(page.index)" :class="pageLinkClass" tabindex="0">{{ page.content }}</a>
    </li>
    <li :class="nextClass">
      <a @click="nextPage()" :class="nextLinkClass" tabindex="0">{{ nextText }}</a>
    </li>
  </ul>
</template>

<script>
export default {
  name: 'Pager',

  props: {
    pageCount: {
      type: Number,
      required: true,
    },
    initialPage: {
      type: Number,
      default: 1,
    },
    forcePage: {
      type: Number,
      default: 1,
    },
    clickHandler: {
      type: Function,
      default: () => {},
    },
    pageRange: {
      type: Number,
      default: 3,
    },
    marginPages: {
      type: Number,
      default: 1,
    },
    prevText: {
      type: String,
      default: 'Prev',
    },
    nextText: {
      type: String,
      default: 'Next',
    },
    containerClass: {
      type: String,
    },
    pageClass: {
      type: String,
    },
    pageLinkClass: {
      type: String,
    },
    prevClass: {
      type: String,
    },
    prevLinkClass: {
      type: String,
    },
    nextClass: {
      type: String,
    },
    nextLinkClass: {
      type: String,
    },
  },

  data() {
    return {
      selected: this.initialPage - 1,
    };
  },

  beforeUpdate() {
    if (this.forcePage !== this.selected) {
      this.selected = this.forcePage;
    }
  },

  computed: {
    pages: function () { // eslint-disable-line
      const items = {};
      if (this.pageCount <= this.pageRange) {
        for (let index = 0; index < this.pageCount; index++) { // eslint-disable-line
          const page = {
            index,
            content: index + 1,
            selected: index === this.selected,
          };
          items[index] = page;
        }
      } else {
        let leftPart = this.pageRange / 2;
        let rightPart = this.pageRange - leftPart;

        if (this.selected < leftPart) {
          leftPart = this.selected;
          rightPart = this.pageRange - leftPart;
        } else if (this.selected > this.pageCount - this.pageRange / 2) { // eslint-disable-line
          rightPart = this.pageCount - this.selected;
          leftPart = this.pageRange - rightPart;
        }

        for (let index = 0; index < this.pageCount; index++) { // eslint-disable-line
          const page = {
            index,
            content: index + 1,
            selected: index === this.selected,
          };

          if (index <= this.marginPages - 1 || index >= this.pageCount - this.marginPages) {
            items[index] = page;
            continue; // eslint-disable-line
          }

          const breakView = {
            content: '...',
            disabled: true,
          };

          if ((this.selected - leftPart) > this.marginPages && items[this.marginPages] !== breakView) { // eslint-disable-line
            items[this.marginPages] = breakView;
          }

          if ((this.selected + rightPart) < (this.pageCount - this.marginPages - 1) && items[this.pageCount - this.marginPages - 1] !== breakView) { // eslint-disable-line
            items[this.pageCount - this.marginPages - 1] = breakView;
          }

          let overCount = this.selected + rightPart - this.pageCount + 1; // eslint-disable-line

          if (overCount > 0 && index === this.selected - leftPart - overCount) {
            items[index] = page;
          }

          if ((index >= this.selected - leftPart) && (index <= this.selected + rightPart)) {
            items[index] = page;
            continue; // eslint-disable-line
          }
        }
      }
      return items;
    },
  },

  methods: {
    handlePageSelected(selected) {
      if (this.selected === selected) return;

      this.selected = selected;

      this.clickHandler(this.selected + 1);
    },
    prevPage() {
      if (this.selected <= 0) return;

      this.selected--; // eslint-disable-line

      this.clickHandler(this.selected + 1);
    },
    nextPage() {
      if (this.selected >= this.pageCount - 1) return;

      this.selected++; // eslint-disable-line

      this.clickHandler(this.selected + 1);
    },
  },
};
</script>

<style lang="css" scoped>
a {
  cursor: pointer;
}
</style>
