import { defineStore } from 'pinia';
import api from '../api/axios';

export const useAuthStore = defineStore('auth', {
    state: () => ({
        token: localStorage.getItem('token') || '',
        user: null,
    }),
    getters: {
        isAuthenticated: (state) => !!state.token,
    },
    actions: {
        async login(username, password) {
            try {
                const response = await api.post('/auth/login', { username, password });
                this.token = response.data.token;
                localStorage.setItem('token', this.token);
                // Fetch user profile if needed
                // await this.fetchUser();
                return true;
            } catch (error) {
                console.error('Login failed:', error);
                throw error;
            }
        },
        logout() {
            this.token = '';
            this.user = null;
            localStorage.removeItem('token');
        },
    },
});
