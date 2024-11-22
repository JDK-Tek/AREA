// import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'

import { createMemoryHistory, createRouter } from 'vue-router'

import HomeView from './components/Home.vue'
import LoginView from './components/Login.vue'

const routes = [
  { path: '/', component: LoginView },
  { path: '/home', component: HomeView },
]

const router = createRouter({
  history: createMemoryHistory(),
  routes,
})

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const vuetify = createVuetify({
  components,
  directives,
})

createApp(App).use(vuetify).use(router).mount('#app')
