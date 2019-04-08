<template>
  <v-dialog v-model="show" width="800px" persistent>
    <v-card>
      <v-card-title class="py-4 title">API Key Required</v-card-title>
      <v-container grid-list-sm>
        <v-layout row wrap>
          <v-flex xs12>
            <v-text-field outline clearable label="API Key" v-model="apiKey" @keypress.enter="save"></v-text-field>
          </v-flex>
        </v-layout>
      </v-container>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn depressed color="primary" @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  data() {
    return {
      show: false,
      apiKey: ''
    }
  },
  methods: {
    ...mapActions(['setApiKey']),
    open() {
      this.setApiKey('')
      this.show = true
    },
    close() {
      this.apiKey = ''
      this.show = false
    },
    save() {
      this.setApiKey(this.apiKey)
      this.close()
    }
  }
}
</script>
