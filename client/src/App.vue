<template>
  <v-app id="inspire">
    <v-toolbar :clipped-left="$vuetify.breakpoint.lgAndUp"
      color="blue darken-3" dark app fixed flat>
      <v-toolbar-title style="width: 300px" class="ml-0 pl-3">
        <div class="clickable" @click="$router.push('/')">URL Shortener</div>
        <v-btn fab absolute bottom right color="success" dark fixed @click="$refs.shortUrlDialog.open()">
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
      ></v-text-field> -->
      <v-spacer></v-spacer>
    </v-toolbar>
    <v-content>
      <v-container fluid fill-height>
        <v-layout>
          <v-flex md10 offset-md1>
            <router-view></router-view>
          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
    <shorten-url-dialog ref="shortUrlDialog"></shorten-url-dialog>
  </v-app>
</template>

<script>
import ShortenUrlDialog from "@/components/ShortenUrlDialog";
import { mapActions } from "vuex";

export default {
  components: {
    ShortenUrlDialog
  },
  props: {
    source: String
  },
  methods: {
    ...mapActions(["setApiKey"])
  },
  created() {
    this.setApiKey("123456789");
  }
};
</script>

<style lang="stylus" scoped>
.clickable
  display inline-block
  cursor pointer
</style>
