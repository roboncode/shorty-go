import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    apiKey: "",
    links: []
  },
  mutations: {
    setApiKey(state, apiKey) {
        state.apiKey = apiKey
    },
    setLinks(state, links) {
      state.links = links
    }
  },
  actions: {
    setApiKey({commit}, apiKey) {
      commit('setApiKey', apiKey)
    },
    createLink({commit}, url) {
      return axios.post('/shorten', {url}, {
        params: {
          key: state.apiKey
        }
      })
    },
    fetchLinks({commit, state}) {
      return axios.get("/links", {
        params: {
          key: state.apiKey
        }
      }).then(({data}) => {
        commit('setLinks', data)
      })
    },
    fetchLink({commit, state}, code) {
      return axios.get(`/links/${code}`, {
        params: {
          key: state.apiKey
        }
      })
    },
    deleteLink({commit}, code) {
      return axios.delete(`/links/${code}`, {
        params: {
          key: state.apiKey
        }
      })
    }
  }
})
