import axios from "axios";

const api = axios.create({
  // Vite injects the correct URL based on your .env file
  //baseURL: import.meta.env.VITE_API_BASE_URL,
  baseURL: (import.meta as unknown as { env: { VITE_API_BASE_URL: string } }).env.VITE_API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

export default api;
