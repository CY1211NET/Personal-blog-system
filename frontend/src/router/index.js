import { createRouter, createWebHistory } from 'vue-router';

const routes = [
    {
        path: '/',
        name: 'Home',
        component: () => import('../views/ArticleList.vue'),
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/Login.vue'),
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('../views/Register.vue'),
    },
    {
        path: '/articles/new',
        name: 'CreateArticle',
        component: () => import('../views/ArticleEditor.vue'),
    },
    {
        path: '/articles/:id',
        name: 'ArticleDetail',
        component: () => import('../views/ArticleDetail.vue'),
    },
    {
        path: '/articles/:id/edit',
        name: 'EditArticle',
        component: () => import('../views/ArticleEditor.vue'),
    },
    {
        path: '/profile',
        component: () => import('../views/Profile.vue'),
    },
    {
        path: '/timeline',
        name: 'Timeline',
        component: () => import('../views/Timeline.vue'),
    },
    {
        path: '/about',
        name: 'About',
        component: () => import('../views/About.vue'),
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
