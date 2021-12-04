import Vue from 'vue'
import VueRouter from 'vue-router'
import Line from '../views/Line'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Line
  }
]

const router = new VueRouter({
  routes
})

export default router
