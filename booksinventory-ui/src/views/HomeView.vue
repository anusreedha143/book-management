<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../../src/services/api'
import { useRouter } from 'vue-router'

// PrimeVue Components
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
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

const books = ref<Book[]>([])
const error = ref<string | null>(null)
const router = useRouter()

// Fetch all books
onMounted(async () => {
  try {
    const response = await api.get('/books')
    // Bulletproof assignment to ensure it never crashes
    books.value = response?.data?.books || []
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
    books.value = books.value.filter(book => book.id !== id)
  } catch (err) {
    alert("Failed to delete the book. Please try again.")
    console.error(err)
  }
}

// Handle row click
const onRowClick = (event: { data: Book }) => {
  console.log("Row clicked:", event.data);
  router.push(`/book/${event.data?.id}`)
}
</script>

<template>
  <article style="width: 100%; display: block; text-align: left; margin: 0; padding: 1rem 0;">
    <Message v-if="error" severity="error" :closable="false">{{ error }}</Message>
    
    <DataTable v-if="books.length > 0" :value="books" stripedRows showGridlines rowHover responsiveLayout="scroll" @row-click="onRowClick">
      <template #empty> No books found. </template>

      <Column field="title" header="Title">
        <template #body="slotProps">
          <span style="color: var(--p-primary-color); font-weight: bold;">
            {{ slotProps.data.title }}
          </span>
        </template>
      </Column>
      <Column field="pages" header="Pages"></Column>
      <Column field="published" header="Published"></Column>
      <Column field="rating" header="Rating"></Column>

      <Column header="Actions" :exportable="false" style="min-width: 12rem">
        <template #body="slotProps">
          <div @click.stop>
            <Button icon="pi pi-trash" label="Delete" severity="danger" outlined size="small" style="margin-right: 0.5rem" @click="deleteBook(slotProps.data.id, slotProps.data.title)" />
            <router-link :to="`/edit/${slotProps.data.id}`" style="text-decoration: none;">
              <Button icon="pi pi-pencil" label="Edit" severity="secondary" size="small" />
           </router-link>
          </div>
        </template>
      </Column>
    </DataTable>

    <p v-else-if="!error">There's nothing to see here yet!</p>
  </article>
</template>
<style scoped>
/* Target the table rows directly */
:deep(.p-datatable-tbody > tr:hover) {
  background-color: #f1f5f9 !important; /* A nice light slate color */
  cursor: pointer;
}
</style>