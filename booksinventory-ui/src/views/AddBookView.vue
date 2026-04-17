<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../src/services/api'

// 1. Initialize Vue Router so we can redirect the user after success
const router = useRouter()

// 2. Create reactive variables for the form inputs
const title = ref('')
const pages = ref<number>(250)
const published = ref<number>(new Date().getFullYear())
const genres = ref('') // Kept as string for the input field
const rating = ref<number>(0)
const errorMessage = ref('')

// 3. The function that runs when the form is submitted
const submitBook = async () => {
  try {
    // Transform the comma-separated genres string into an array of strings
    const genresArray = genres.value
      .split(',')
      .map(g => g.trim())
      .filter(g => g !== '')

    // Construct the JSON payload exactly how your Go API expects it
    const payload = {
      title: title.value,
      pages: pages.value,
      published: published.value,
      genres: genresArray,
      rating: rating.value
    }

    // Send the POST request to your Go API
    await api.post('/books', payload)

    // On success, redirect the user back to the Home page
    router.push('/')
  } catch (error) {
    console.error("Failed to create book:", error)
    errorMessage.value = "Failed to create the book. Please check your inputs."
  }
}
</script>

<template>
  <div>
    <h2>Create a New Book Entry</h2>
    
    <p v-if="errorMessage" style="color: red;">{{ errorMessage }}</p>

    <form @submit.prevent="submitBook">
      <label>Title:</label>
      <input type="text" v-model="title" required maxlength="150"><br>

      <label>Pages:</label>
      <input type="number" v-model.number="pages" required min="1" max="9999"><br>

      <label>Published:</label>
      <input type="number" v-model.number="published" required min="1000" max="9999" step="1"><br>

      <label>Genres (comma separated):</label>
      <input type="text" v-model="genres" placeholder="e.g. Sci-Fi, Fantasy" required><br>

      <label>Rating:</label>
      <input type="number" step="0.1" v-model.number="rating" required min="0" max="5"><br>
      
      <div class="button-center">
        <button type="submit">Submit</button>
      </div>
    </form>
  </div>
</template>