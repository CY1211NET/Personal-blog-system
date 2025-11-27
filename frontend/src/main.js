import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import './style.css'
import App from './App.vue'

const app = createApp(App)
const pinia = createPinia()

// Simple Lazy Load Directive
app.directive('lazy', {
    mounted: (el) => {
        const loadImage = () => {
            if (el.dataset.src) {
                el.src = el.dataset.src;
            }
        };

        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    loadImage();
                    observer.unobserve(el);
                }
            });
        });

        observer.observe(el);
    }
});

app.use(pinia)
app.use(router)
app.mount('#app')
