import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

// 1. Import PrimeVue, the Aura theme, and PrimeIcons
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
// import 'primeicons/primeicons.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 2. Initialize PrimeVue with the Aura preset
app.use(PrimeVue, {
    theme: {
        preset: Aura
    }
})

app.mount('#app')