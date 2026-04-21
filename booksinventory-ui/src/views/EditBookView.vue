<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../src/services/api'

import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Button from 'primevue/button'
import Message from 'primevue/message'

interface Book {
  id: number;
  title: string;
  pages: number;
  published: number;
  genres: string[];
  rating: number;
}

interface BookResponse {
  book: Book;
}

const route = useRoute()
const router = useRouter()
const bookId = route.params.id

// Form State
const title = ref('')
const pages = ref<number | null>(null)
const published = ref<number | null>(null)
const genres = ref('') 
const rating = ref<number | null>(null)

// Validation and Loading State
const isLoading = ref(true)
const errors = ref<Record<string, string>>({})
const apiError = ref('')

onMounted(async () => {
  try {
    const response = await api.get<BookResponse>(`/books/${bookId}`)
    const bookData = response?.data?.book
    
    title.value = bookData.title
    pages.value = bookData.pages
    published.value = bookData.published
    genres.value = bookData.genres.join(', ')
    rating.value = bookData.rating
  } catch (err) {
    console.error("Failed to fetch book for editing:", err)
    apiError.value = "Failed to load book data."
  } finally {
    isLoading.value = false
  }
})

const updateBook = async () => {
  // Clear previous errors
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

  // --- API Submission ---
  try {
    const genresArray = genres.value
      .split(',')
      .map(g => g.trim())
      .filter(g => g !== '')

    const payload = {
      title: title.value,
      pages: pages.value,
      published: published.value,
      genres: genresArray,
      rating: rating.value
    }

    await api.put(`/books/${bookId}`, payload)
    router.push('/')
  } catch (error) {
    console.error("Failed to update book:", error)
    apiError.value = "Failed to save changes. Please try again."
  }
}
</script>

<template>
  <div style="display: flex; justify-content: center; padding: 2rem;">
    <Card style="width: 100%; max-width: 500px; box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);">
      
      <template #title>
        <h2 style="margin: 0; color: var(--p-primary-color); text-align: center;">Edit Book</h2>
      </template>

      <template #content>
        <div v-if="isLoading" style="text-align: center; color: #64748b; padding: 2rem 0;">
          <i class="pi pi-spin pi-spinner" style="font-size: 2rem; margin-bottom: 1rem;"></i>
          <p>Loading book data...</p>
        </div>

        <Message v-if="apiError" severity="error" :closable="false" style="margin-bottom: 1.5rem;">
          {{ apiError }}
        </Message>

        <form v-if="!isLoading && !apiError" @submit.prevent="updateBook" style="display: flex; flex-direction: column; gap: 1.5rem;">
          
          <div style="display: flex; flex-direction: column; gap: 0.25rem;">
            <label for="title" style="font-weight: 600;">Title</label>
            <InputText id="title" v-model="title" maxlength="150" />
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
            <InputText id="genres" v-model="genres" />
            <small v-if="errors.genres" style="color: #ef4444; font-weight: bold;">{{ errors.genres }}</small>
          </div>

          <div style="display: flex; flex-direction: column; gap: 0.25rem;">
            <label for="rating" style="font-weight: 600;">Rating</label>
            <InputNumber id="rating" v-model="rating" :step="0.1" :maxFractionDigits="1" showButtons fluid />
            <small v-if="errors.rating" style="color: #ef4444; font-weight: bold;">{{ errors.rating }}</small>
          </div>

          <div style="display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem;">
            <Button label="Cancel" icon="pi pi-times" severity="secondary" outlined @click="router.push('/')" type="button" />
            <Button label="Save Changes" icon="pi pi-check" type="submit" />
          </div>
          
        </form>
      </template>
    </Card>
  </div>
</template>