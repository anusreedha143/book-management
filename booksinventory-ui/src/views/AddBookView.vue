<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../src/services/api'

// --- PrimeVue Imports ---
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Button from 'primevue/button'
import Message from 'primevue/message'

// 1. Initialize Vue Router so we can redirect the user after success
const router = useRouter()

// 2. Create reactive variables for the form inputs
const title = ref('')
const pages = ref<number | null>(null)
const published = ref<number | null>(null)
const genres = ref('') // Kept as string for the input field
const rating = ref<number | null>(null)

const errors = ref<Record<string, string>>({})
const apiError = ref('') // For overall server errors

// 3. The function that runs when the form is submitted
const submitBook = async () => {
 
  errors.value = {}
  apiError.value = ''
  
  const currentYear = new Date().getFullYear()

  // --- Custom Validation Logic ---
  if (!title.value.trim()) {
    errors.value.title = "Title is required."
  }
  
  if (pages.value === null || pages.value < 1 || pages.value > 9999) {
    errors.value.pages = "Pages must be between 1 and 9,999."
  }
  
  if (published.value === null || published.value < 1000 || published.value > currentYear) {
    errors.value.published = `Published year must be between 1000 and ${currentYear}.`
  }

  if (!genres.value.trim()) {
    errors.value.genres = "At least one genre is required."
  }

  if (rating.value === null || rating.value < 0 || rating.value > 5) {
    errors.value.rating = "Rating must be between 0 and 5."
  }

  // If there are any errors in the object, stop the submission!
  if (Object.keys(errors.value).length > 0) {
    return
  }

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
    apiError.value = "Failed to create the book on the server. Please try again."
  }
}
</script>

<template>
  <div style="display: flex; justify-content: center; padding: 2rem;">
    <Card style="width: 100%; max-width: 500px; box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);">
      
      <template #title>
        <h2 style="margin: 0; color: var(--p-primary-color); text-align: center;">Create a New Book Entry</h2>
      </template>

      <template #content>
        <Message v-if="apiError" severity="error" :closable="false" style="margin-bottom: 1.5rem;">
          {{ apiError }}
        </Message>
        <!-- Purpose of adding @submit.prevent is to prevent the default form submission behavior (prevent that from bubbling down to whatever element like the form to do the default behaviour) -->
        <form @submit.prevent="submitBook" style="display: flex; flex-direction: column; gap: 1.5rem;">
          
          <div style="display: flex; flex-direction: column; gap: 0.25rem;">
            <label for="title" style="font-weight: 600;">Title</label>
            <InputText id="title" v-model="title" maxlength="150" placeholder="e.g. The Hobbit" />
            <small v-if="errors.title" style="color: #ef4444; font-weight: bold;">{{ errors.title }}</small>
          </div>

          <div style="display: flex; gap: 1rem;">
            <div style="display: flex; flex-direction: column; gap: 0.25rem; flex: 1;">
              <label for="pages" style="font-weight: 600;">Pages</label>
              <InputNumber id="pages" v-model="pages" />
              <small v-if="errors.pages" style="color: #ef4444; font-weight: bold;">{{ errors.pages }}</small>
            </div>

            <div style="display: flex; flex-direction: column; gap: 0.25rem; flex: 1;">
              <label for="published" style="font-weight: 600;">Published</label>
              <InputNumber id="published" v-model="published" :useGrouping="false" />
              <small v-if="errors.published" style="color: #ef4444; font-weight: bold;">{{ errors.published }}</small>
            </div>
          </div>

          <div style="display: flex; flex-direction: column; gap: 0.25rem;">
            <label for="genres" style="font-weight: 600;">Genres (comma separated)</label>
            <InputText id="genres" v-model="genres" placeholder="e.g. Sci-Fi, Fantasy" />
            <small v-if="errors.genres" style="color: #ef4444; font-weight: bold;">{{ errors.genres }}</small>
          </div>

          <div style="display: flex; flex-direction: column; gap: 0.25rem;">
            <label for="rating" style="font-weight: 600;">Rating</label>
            <InputNumber id="rating" v-model="rating" :step="0.1" :maxFractionDigits="1" showButtons fluid placeholder="0.0"/>
            <small v-if="errors.rating" style="color: #ef4444; font-weight: bold;">{{ errors.rating }}</small>
          </div>

          <div style="display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem;">
            <Button label="Cancel" icon="pi pi-times" severity="secondary" outlined @click="router.push('/')" type="button" />
            <Button label="Submit" icon="pi pi-check" type="submit" />
          </div>
          
        </form>
      </template>
    </Card>
  </div>
</template>