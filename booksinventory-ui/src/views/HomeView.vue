<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../../src/services/api'

interface Book {
  id: number;
  title: string;
  pages: number;
  published: number;
  genres: string[];
  rating: number;
}

interface BooksResponse {
  books: Book[];
}

const books = ref<Book[]>([])
const error = ref<string | null>(null)

// Fetch all books
onMounted(async () => {
  try {
    const response = await api.get<BooksResponse>('/books')
    books.value = response.data.books 
  } catch (err: unknown) {
    error.value = "Failed to load books from the API."
    console.error("API Error:", err)
  }
})

// Delete a book
const deleteBook = async (id: number, title: string) => {
  if (!confirm(`Are you sure you want to delete "${title}"?`)) return
  
  try {
    await api.delete(`/books/${id}`)
    // Instantly remove the deleted book from the UI table without refreshing
    books.value = books.value.filter(book => book.id !== id)
  } catch (err) {
    alert("Failed to delete the book. Please try again.")
    console.error(err)
  }
}
</script>

<template>
  <article>
    <p v-if="error" style="color: red;">{{ error }}</p>
    
    <table v-if="books.length > 0">
      <thead>
        <tr>
          <th>Title</th>
          <th>Pages</th>
          <th>Published</th>
          <th>Rating</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="book in books" :key="book.id">
          <td><router-link :to="`/book/${book.id}`">{{ book.title }}</router-link></td>
          <td>{{ book.pages }}</td>
          <td>{{ book.published }}</td>
          <td>{{ book.rating }}</td>
          <td>
            <button @click="deleteBook(book.id, book.title)" class="btn-delete" style="margin-right: 10px;">
              Delete
            </button>
            <router-link :to="`/edit/${book.id}`">Edit</router-link>
          </td>
        </tr>
      </tbody>
    </table>
    
    <p v-else-if="!error">There's nothing to see here yet!</p>
  </article>
</template>