<template>
  <section class="row">
    <div class="col-xs-12 mb2">
      <form>
        <fieldset>
          <legend>Where are your movies stored ?</legend>
          <div class="row mb3">
            <div class="col-xs-12 addon">
              <span class="addon-item">Folder</span>
              <input class="addon-field" type="text" v-model="folder" />
              <button class="btn btn-default" @click="onAdd">Add</button>
            </div>
          </div>
          <div class="row mb3">
            <div class="col-xs-12">
              <table>
                <thead>
                  <th width="50">#</th>
                  <th>Folder</th>
                </thead>
                <tbody>
                  <tr v-for="folder in folders">
                    <td><i class="fa fa-times-circle" /></td>
                    <td>{{ folder }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </fieldset>
      </form>
    </div>
  </section>
</template>

<script>
import * as types from '../store/types'

export default {
  name: 'settings',

  data () {
    return {
      folder: ''
    }
  },

  methods: {
    onAdd () {
      // console.log(`folder-${this.folder}`)
      this.$store.dispatch(types.ADD_FOLDER, this.folder)
    }
  },

  computed: {
    folders () {
      return this.$store.state.config ? this.$store.state.config.mediaFolders : []
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
