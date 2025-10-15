import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/home/HomeView.vue'
import EchoView from '../views/echo/EchoView.vue'
import NotFoundView from '../views/404/NotFoundView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: {
        title: 'Home',
      },
    },
    {
      path: '/panel',
      name: 'panel',
      component: () => import('../views/panel/PanelView.vue'),
      redirect: '/panel/status',
      children: [
        {
          path: 'status',
          name: 'panel-status',
          component: () => import('../views/panel/modules/TheStatus.vue'),
        },
        {
          path: 'setting',
          name: 'panel-setting',
          component: () => import('../views/panel/modules/TheSetting.vue'),
        },
        {
          path: 'user',
          name: 'panel-user',
          component: () => import('../views/panel/modules/TheUser.vue'),
        },
        {
          path: 'advance',
          name: 'panel-advance',
          component: () => import('../views/panel/modules/TheAdvance.vue'),
        },
      ],
    },
    {
      path: '/auth',
      name: 'auth',
      component: () => import('../views/auth/AuthView.vue'),
    },
    {
      path: '/connect',
      name: 'connect',
      component: () => import('../views/connect/ConnectView.vue'),
    },
    {
      path: '/echo/:echoId',
      name: 'echo',
      component: EchoView,
    },
    // {
    //   path: '/fediverse',
    //   name: 'fediverse',
    //   component: () => import('../views/fediverse/FediverseView.vue'),
    // },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: NotFoundView,
    },
  ],
})

export default router
