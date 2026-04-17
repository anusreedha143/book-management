<!-- <script setup lang="ts"></script>

<template>
  <h1>You did it!</h1>
  <p>
    Visit <a href="https://vuejs.org/" target="_blank" rel="noopener">vuejs.org</a> to read the
    documentation
  </p>
</template>

<style scoped></style> I want to know what is this style scoped does --> 
<script setup lang="ts">
// Import your navigation component
import NavBar from './components/NavBar.vue'
</script>

<template>
  <header>
    <h1><router-link to="/">Reading List</router-link></h1>
  </header>
  
  <NavBar />
  
  <main>
    <router-view />
  </main>
  
  <footer>
    Powered by <a href='https://vuejs.org/'>Vue.js</a> & <a href='https://golang.org/'>Go</a>
  </footer>
</template>

<style>
/* This imports your Go application's CSS so it applies to the whole Vue app! */
@import './assets/main.css';
</style>
