<template>
  <v-app id="inspire">
    <v-toolbar
      :clipped-left="$vuetify.breakpoint.lgAndUp"
      color="blue darken-3"
      dark
      app
      fixed
      flat
    >
      <v-toolbar-title style="width: 300px" class="ml-0 pl-3">
        <div class="clickable" @click="$router.push('/')">URL Shortener</div>
        <v-btn
          fab
          absolute
          bottom
          right
          color="success"
          dark
          fixed
          @click="$refs.shortUrlDialog.open()"
        >
          <v-icon>add</v-icon>
        </v-btn>
      </v-toolbar-title>
      <!-- <v-text-field
        flat
        solo-inverted
        hide-details
        prepend-inner-icon="search"
        label="Search"
        class="hidden-sm-and-down"
      ></v-text-field>-->
      <v-spacer></v-spacer>
    </v-toolbar>
    <v-content>
      <v-container fluid fill-height>
        <v-layout>
          <v-flex md10 offset-md1>
            <router-view v-if="apiKey"></router-view>
          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
    <shorten-url-dialog ref="shortUrlDialog"></shorten-url-dialog>
    <api-key-dialog ref="apiKeyDialog"></api-key-dialog>
    <v-snackbar v-model="showSnackbar" :bottom="true" :left="true" :timeout="3000">{{ message }}</v-snackbar>
  </v-app>
</template>

<script>
import ApiKeyDialog from './components/ApiKeyDialog'
import ShortenUrlDialog from './components/ShortenUrlDialog'
import { mapActions, mapState } from 'vuex'
import axios from 'axios'

export default {
  components: {
    ApiKeyDialog,
    ShortenUrlDialog
  },
  props: {
    source: String
  },
  data() {
    return {
      message: "",
      showSnackbar: false
    }
  },
  computed: {
    ...mapState(['apiKey'])
  },
  methods: {
    ...mapActions(['setApiKey'])
  },
  created() {
    axios.interceptors.response.use(
      response => {
        return response
      },
      error => {
        if (error.response.status === 401) {
          this.message = "Invalid API Key"
          this.showSnackbar = true
          this.$refs.apiKeyDialog.open()
        } else {
          return Promise.reject(error)
        }
      }
    )
  },
  mounted() {
    if (!this.apiKey) {
      this.$refs.apiKeyDialog.open()
    }
  }
}
</script>

<style lang="stylus" scoped>
.clickable
  display inline-block
  cursor pointer
</style>
