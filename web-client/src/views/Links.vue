<template>
  <v-layout column>
    <v-list class="list" two-line>
      <v-list-tile
        class="list-item"
        v-for="item in links"
        :key="item.id"
        avatar
        :to="{ name:'link', params:{ code:item.code } }"
      >
        <v-list-tile-avatar>
          <v-icon>link</v-icon>
        </v-list-tile-avatar>

        <v-list-tile-content>
          <v-list-tile-title>{{ item.shortUrl }}</v-list-tile-title>
          <v-list-tile-sub-title>{{ item.longUrl }}</v-list-tile-sub-title>
        </v-list-tile-content>

        <v-list-tile-action>
          <v-btn icon ripple dark color="primary" @click.stop.prevent="openLink(item)">
            <v-icon>search</v-icon>
          </v-btn>
        </v-list-tile-action>
      </v-list-tile>
    </v-list>
    <v-footer app fixed color="primary">
      <v-layout justify-center>
        <v-btn small flat dark :disabled="!hasPrev" @click="prev">
          <v-icon>keyboard_arrow_left</v-icon>Prev
        </v-btn>
        <v-btn small flat dark :disabled="!hasNext" @click="next">
          Next
          <v-icon>keyboard_arrow_right</v-icon>
        </v-btn>
      </v-layout>
    </v-footer>
  </v-layout>
</template>

<script>
import { mapActions, mapState } from 'vuex'

export default {
  computed: {
    ...mapState(['links', 'limit']),
    hasPrev() {
      return Boolean(this.$route.query.p)
    },
    hasNext() {
      return this.links.length === this.limit
    }
  },
  methods: {
    ...mapActions(['setApiKey', 'fetchLinks']),
    openLink(link) {
      window.open(link.longUrl, '_blank')
    },
    prev() {
      this.$router.go(-1)
    },
    next() {
      let page = Number(this.$route.query.p || 0)
      this.$router.push({
        name: 'links',
        query: {
          p: page + 1
        }
      })
    }
  },
  watch: {
    '$route.query.p'(val) {
      this.fetchLinks(val)
    }
  },
  created() {
    this.fetchLinks(this.$route.query.p)
  }
}
</script>

<style lang="stylus" scoped>
.list
  background transparent !important

  .list-item
    border-bottom 1px solid rgba(0, 0, 0, 0.1)

  .list-item:last-child
    border-bottom none
</style>
