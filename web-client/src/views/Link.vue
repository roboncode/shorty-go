<template>
  <v-container fill-height align-center column fluid>
    <v-flex xs12>
      <v-layout v-if="busy" column align-center>
        <v-progress-circular indeterminate size="48" color="grey lighten-2"></v-progress-circular>
      </v-layout>
      <v-layout v-else column align-center justify-center wrap>
        <a :href="link.shortUrl" target="_blank" class="shortUrl">{{link.shortUrl}}</a>
        <a :href="link.longUrl" target="_blank" class="longUrl">{{link.longUrl}}</a>
      </v-layout>
    </v-flex>
    <v-footer v-if="!busy" app fixed>
      <v-layout justify-end>
        <v-btn small flat color="primary" :to="{name: 'links'}">
          <v-icon left>arrow_back</v-icon>
          <span>Back to links</span>
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn small depressed color="error" @click="remove(link)">Delete</v-btn>
      </v-layout>
    </v-footer>
  </v-container>
</template>

<script>
import { mapActions, mapState } from 'vuex'
export default {
  data() {
    return {
      busy: true
    }
  },
  computed: {
    ...mapState(['link'])
  },
  methods: {
    ...mapActions(['fetchLink', 'deleteLink']),
    remove(link) {
      this.busy = true
      this.deleteLink(link.code).then(() => {
        this.$router.push({ name: 'links' })
      })
    },
    getLinks() {
      this.busy = true
      this.fetchLink(this.$route.params.code).then(() => {
        this.busy = false
      })
    }
  },
  watch: {
    $route() {
      this.getLinks()
    }
  },
  created() {
    this.getLinks()
  }
}
</script>

<style lang="stylus" scoped>
.shortUrl
  padding 10px 20px
  background #1365c0
  color white
  font-size 24px
  line-height 32px
  border-radius 8px
  text-decoration none
  cursor pointer
  white-space nowrap
  overflow hidden
  text-overflow ellipsis

  @media only screen and (orientation: portrait) and (min-width: 360px)
    font-size 18px
    max-width 310px

  &:hover
    background #267ee2

.longUrl
  padding 16px
  font-size 18px
  text-align center
  /* These are technically the same, but use both */
  overflow-wrap break-word
  word-wrap break-word
  -ms-word-break break-all
  /* This is the dangerous one in WebKit, as it breaks things wherever */
  word-break break-all
  /* Instead use this non-standard one: */
  word-break break-word
  /* Adds a hyphen where the word breaks, if supported (No Blink) */
  -ms-hyphens auto
  -moz-hyphens auto
  -webkit-hyphens auto
  hyphens auto

  @media only screen and (orientation: portrait) and (min-width: 360px)
    font-size 16px
    text-align left
</style>
