<template>
  <v-dialog v-model="show" width="800px">
    <v-card>
      <v-card-title class="grey lighten-4 py-4 title">Shorten URL</v-card-title>
      <v-container grid-list-sm class="pa-4">
        <v-layout row wrap>
          <v-flex xs12>
            <v-text-field prepend-icon="link" outline clearable label="Type URL here..." v-model="url"></v-text-field>
          </v-flex>
        </v-layout>
      </v-container>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn flat color="primary" @click="close">Cancel</v-btn>
        <v-btn flat @click="shortenUrl">Create</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  data() {
    return {
      show: false,
      url: ""
    }
  },
  methods: {
    ...mapActions(['createLink']),
    open() {
      this.url = "https://google.com"
      this.show = true
    },
    close() {
      this.show = false
    },
    shortenUrl() {
      this.createLink(this.url).then(({data}) => {
        this.$router.push(`/links/${data.code}`)
      }, ({data}) => {
        console.log('whoops', data)
      })
      this.close()
    }
  }
}
</script>
