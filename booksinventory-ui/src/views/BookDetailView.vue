<script setup lang="ts">
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
</template>