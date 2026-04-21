<!-- <script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../../src/services/api'

// 1. Define our interfaces
interface Book {
  id: number;
  title: string;
  pages: number;
  published: number;
  genres: string[];
  rating: number;
  version: number;
}

// Your Go API likely wraps the single book response in an envelope like {"book": {...}}
interface BookResponse {
  book: Book;
}

// 2. Initialize the route to access the URL parameters
const route = useRoute()

// 3. Reactive state
const book = ref<Book | null>(null)
const error = ref<string | null>(null)
const isLoading = ref(true)

// 4. Fetch the book when the component loads
onMounted(async () => {
  // Grab the ID from the URL (we will define this parameter in the router)
  const bookId = route.params.id

  try {
    const response = await api.get<BookResponse>(`/books/${bookId}`)
    book.value = response.data.book
  } catch (err) {
    console.error("Failed to fetch book:", err)
    error.value = "Failed to load book details. It may have been deleted."
  } finally {
    // Stop the loading indicator whether it succeeded or failed
    isLoading.value = false
  }
})
</script>

<template>
  <div>
    <p v-if="isLoading">Loading book details...</p>
    <p v-else-if="error" style="color: red;">{{ error }}</p>

    <div v-else-if="book" class="book-details">
      <h2>{{ book.title }}</h2>
      
      <ul>
        <li><strong>Title:</strong> {{ book.title }}</li>
        <li><strong>Published:</strong> {{ book.published }}</li>
        <li><strong>Pages:</strong> {{ book.pages }}</li>
        
        <li><strong>Genres:</strong> {{ book.genres.join(", ") }}</li>
        
        <li><strong>Rating:</strong> {{ book.rating }}</li>
      </ul>
    </div>
  </div>
</template> -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../src/services/api' // Adjust path if your api service is located elsewhere

import Card from 'primevue/card'
import Tag from 'primevue/tag'
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

const book = ref<Book | null>(null)
const isLoading = ref(true)
const apiError = ref('')

onMounted(async () => {
  try {
    const response = await api.get<BookResponse>(`/books/${bookId}`)
    // Safely extract the book data
    book.value = response?.data?.book  || (response?.data as unknown as Book)
  } catch (err) {
    console.error("Failed to fetch book details:", err)
    apiError.value = "Failed to load book data."
  } finally {
    isLoading.value = false
  }
})
</script>
<template>
  <div style="display: flex; justify-content: center; padding: 2rem; width: 100%;">
    
    <Card style="width: 500px; max-width: 100%; box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);">
      
      <template #title>
        <h2 v-if="book" style="margin: 0; color: var(--p-primary-color); text-align: center;">
          {{ book.title }}
        </h2>
      </template>

      <template #subtitle>
        <div v-if="book" style="text-align: center; font-size: 1.1rem;">
          Published in {{ book.published }}
        </div>
      </template>

      <template #content>
        <div v-if="isLoading" style="text-align: center; color: #64748b; padding: 2rem 0;">
          <i class="pi pi-spin pi-spinner" style="font-size: 2rem; margin-bottom: 1rem;"></i>
          <p>Loading book details...</p>
        </div>

        <Message v-if="apiError" severity="error" :closable="false" style="margin-bottom: 1.5rem;">
          {{ apiError }}
        </Message>

        <div v-if="!isLoading && !apiError && book" style="display: flex; flex-direction: column; gap: 1.5rem; font-size: 1.1rem; margin-top: 1rem;">
          
          <div style="display: flex; justify-content: space-between; border-bottom: 1px solid #e2e8f0; padding-bottom: 0.5rem;">
            <strong style="font-weight: 600;">Pages:</strong> 
            <span>{{ book.pages }}</span>
          </div>
          
          <div style="display: flex; justify-content: space-between; border-bottom: 1px solid #e2e8f0; padding-bottom: 0.5rem;">
            <strong style="font-weight: 600;">Rating:</strong> 
            <span>
              <i class="pi pi-star-fill" style="color: #fbbf24; margin-right: 0.25rem;"></i>
              {{ book.rating }} / 5
            </span>
          </div>
          
          <div style="display: flex; flex-direction: column; gap: 0.5rem; padding-top: 0.5rem;">
            <strong style="font-weight: 600;">Genres:</strong>
            <div style="display: flex; gap: 0.5rem; flex-wrap: wrap;">
              <Tag v-for="genre in book.genres" :key="genre" :value="genre" severity="info" rounded />
            </div>
          </div>
          
        </div>
      </template>

      <template #footer>
        <div v-if="!isLoading && !apiError && book" style="display: flex; justify-content: space-between; gap: 1rem; margin-top: 1.5rem;">
          <Button label="Back" icon="pi pi-arrow-left" severity="secondary" outlined @click="router.push('/')" />
          <Button label="Edit Book" icon="pi pi-pencil" @click="router.push(`/edit/${book.id}`)" />
        </div>
      </template>

    </Card>
  </div>
</template>

