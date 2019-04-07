import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    apiKey: '',
    limit: 2,
    links: [],
    link: {}
  },
  mutations: {
    setApiKey(state, apiKey) {
      state.apiKey = apiKey
    },
    setLinks(state, links) {
      state.links = links
    },
    setLink(state, link) {
      state.link = link
    }
  },
  actions: {
    setApiKey({ commit }, apiKey) {
      commit('setApiKey', apiKey)
    },
    createLink({ state }, url) {
      return axios.post(
        '/shorten',
        { url },
        {
          params: {
            key: state.apiKey
          }
        }
      )
    },
    fetchLinks({ commit, state }, page = 0) {
      return axios
        .get('/links', {
          params: {
            key: state.apiKey,
            l: state.limit,
            s: page * state.limit
          }
        })
        .then(({ data }) => {
          commit('setLinks', data)
        })
    },
    fetchLink({ commit, state }, code) {
      return axios
        .get(`/links/${code}`, {
          params: {
            key: state.apiKey
          }
        })
        .then(({ data }) => {
          commit('setLink', data)
        })
    },
    deleteLink({ state }, code) {
      return axios.delete(`/links/${code}`, {
        params: {
          key: state.apiKey
        }
      })
    }
  }
})
