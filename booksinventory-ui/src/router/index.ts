import { createRouter, createWebHistory } from 'vue-router'

// 1. Import all of your View components
import HomeView from '../views/HomeView.vue'
import AddBookView from '../views/AddBookView.vue'
import ViewBookView from '../views/BookDetailView.vue'
import EditBookView from '../views/EditBookView.vue'

// 2. Define the routes mapping URLs to Components
const router = createRouter({
  // createWebHistory uses the browser's native History API to create clean URLs without the '#' symbol
  // history: createWebHistory(import.meta.env.BASE_URL),
  history: createWebHistory((import.meta as unknown as { env: { BASE_URL: string } }).env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/book/create',
      name: 'add-book',
      component: AddBookView
    },
    {
      // The colon ':' creates a dynamic route parameter. 
      // Vue will capture whatever number is here and pass it to route.params.id
      path: '/book/:id',
      name: 'view-book',
      component: ViewBookView
    },
    {
      // Matches the router-link we put in the HomeView table
      path: '/edit/:id', 
      name: 'edit-book',
      component: EditBookView
    }
  ]
})

export default router