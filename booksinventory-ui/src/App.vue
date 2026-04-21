<script setup lang="ts">
import { useRouter } from 'vue-router'
import Menubar from 'primevue/menubar'
import Button from 'primevue/button'

const router = useRouter()
</script>

<template>
  <div class="app-container">
    <Menubar class="custom-menubar">
      <template #start>
        <div class="logo-container" @click="router.push('/')">
          <i class="pi pi-book logo-icon"></i>
          <span class="logo-text">Reading List</span>
        </div>
      </template>

      <template #end>
        <Button label="Add Book" icon="pi pi-plus" @click="router.push('/book/create')" outlined severity="contrast" size="small" />
      </template>
    </Menubar>
    
    <main>
      <router-view />
    </main>
    
    <footer>
      Powered by <a href='https://vuejs.org/' target="_blank">Vue.js</a> & <a href='https://golang.org/' target="_blank">Go</a>
    </footer>
  </div>
</template>

<style>
/* This imports your existing global CSS */
@import './assets/main.css';
/* 1. Ensure the body takes up the full screen */
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  background-color: #f8fafc; /* Optional: gives a nice subtle gray background to make your white cards pop */
}

/* 2. Force the main Vite div to stretch across the screen */
#app {
  width: 100% !important;
  display: block !important; 
}

/* 3. Update your app-container to ensure it centers */
.app-container {
  width: 100%; /* Tell it to stretch until it hits the max-width */
  max-width: 1200px;
  margin: 0 auto !important; /* The 'auto' pushes equally on the left and right, centering it */
  padding: 1rem;
  min-height: 100vh; 
  display: flex;
  flex-direction: column;
}

.custom-menubar {
  margin-bottom: 2rem;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  padding: 0.5rem 1rem;
}

.logo-container {
  display: flex;
  align-items: center;
  margin-right: 2rem;
  color: var(--p-primary-color);
  cursor: pointer;
}

.logo-icon {
  font-size: 1.5rem;
  margin-right: 0.5rem;
}

.logo-text {
  font-size: 1.5rem;
  font-weight: bold;
}

footer {
  text-align: center;
  padding: 2rem 0;
  margin-top: 2rem;
  color: #64748b;
  border-top: 1px solid #e2e8f0;
}

footer a {
  color: var(--p-primary-color);
  text-decoration: none;
  font-weight: bold;
}

footer a:hover {
  text-decoration: underline;
}

/* Forcefully override Vite's default centering */
/* Forcefully strip Vite's flex layout from the entire page */
html, body, #app {
  display: block !important;
  place-items: start !important;
  justify-content: flex-start !important;
  align-items: flex-start !important;
  min-height: auto !important;
  height: auto !important;
}
</style>
