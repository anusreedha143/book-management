<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../src/services/api'

// 1. Interfaces
interface Book {
  id: number;
  title: string;
  pages: number;
  published: number;
  genres: string[];
  rating: number;
  // Note: We don't need the version here for the update payload based on our earlier Go API fix
}

interface BookResponse {
  book: Book;
}

// 2. Vue Router hooks
const route = useRoute()
const router = useRouter()
const bookId = route.params.id

// 3. Reactive State for the form
const title = ref('')
const pages = ref<number | null>(null)
const published = ref<number | null>(null)
const genres = ref('') // Kept as string for the input field
const rating = ref<number | null>(null)

const isLoading = ref(true)
const errorMessage = ref('')

// 4. Fetch the existing book data when the page loads
onMounted(async () => {
  try {
    const response = await api.get<BookResponse>(`/books/${bookId}`)
    const bookData = response.data.book
    
    // Pre-populate the form fields
    title.value = bookData.title
    pages.value = bookData.pages
    published.value = bookData.published
    // Join the array into a comma-separated string for the input field
    genres.value = bookData.genres.join(', ')
    rating.value = bookData.rating
  } catch (err) {
    console.error("Failed to fetch book for editing:", err)
    errorMessage.value = "Failed to load book data."
  } finally {
    isLoading.value = false
  }
})

// 5. Submit the updated data
const updateBook = async () => {
  try {
    // Transform the comma-separated string back into an array
    const genresArray = genres.value
      .split(',')
      .map(g => g.trim())
      .filter(g => g !== '')

    // Construct the payload (Remember: No 'version' field required for your API!)
    const payload = {
      title: title.value,
      pages: pages.value,
      published: published.value,
      genres: genresArray,
      rating: rating.value
    }

    // Send a PUT request to update the book
    await api.put(`/books/${bookId}`, payload)

    // Redirect back to the home page on success
    router.push('/')
  } catch (error) {
    console.error("Failed to update book:", error)
    errorMessage.value = "Failed to save changes. Please verify your inputs."
  }
}
</script>

<template>
  <div>
    <h1>Edit Book</h1>
    
    <p v-if="isLoading">Loading book data...</p>
    <p v-if="errorMessage" style="color: red;">{{ errorMessage }}</p>

    <div v-if="!isLoading && !errorMessage">
      <form @submit.prevent="updateBook">
        <div>
          <label for="title">Title</label>
          <input type="text" v-model="title" id="title" required>
        </div>

        <div>
          <label for="pages">Pages</label>
          <input type="number" v-model.number="pages" id="pages" required min="1" max="9999">
        </div>

        <div>
          <label for="published">Published</label>
          <input type="number" v-model.number="published" id="published" required min="1000" max="9999" step="1">
        </div>

        <div>
          <label for="genres">Genres</label>
          <input type="text" v-model="genres" id="genres" required>
        </div>

        <div>
          <label for="rating">Rating</label>
          <input type="number" step="0.1" v-model.number="rating" id="rating" min="0" max="5" required>
        </div>

        <button type="submit">Save Changes</button>
      </form>
    </div>
  </div>
</template>