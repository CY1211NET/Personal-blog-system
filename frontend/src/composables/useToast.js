import { ref } from 'vue';

const toasts = ref([]);
let idCounter = 0;

export function useToast() {
    const addToast = (message, type = 'info', duration = 3000) => {
        const id = idCounter++;
        const toast = { id, message, type };
        toasts.value.push(toast);

        setTimeout(() => {
            removeToast(id);
        }, duration);
    };

    const removeToast = (id) => {
        toasts.value = toasts.value.filter(t => t.id !== id);
    };

    return {
        toasts,
        addToast,
        success: (msg) => addToast(msg, 'success'),
        error: (msg) => addToast(msg, 'error'),
        info: (msg) => addToast(msg, 'info')
    };
}
