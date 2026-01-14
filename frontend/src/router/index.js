import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: () => import('../views/HomeView.vue')
    },
    {
        path: '/products',
        name: 'Products',
        component: () => import('../views/ProductsView.vue')
    },
    {
        path: '/products/:id',
        name: 'ProductDetail',
        component: () => import('../views/ProductDetailView.vue')
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/LoginView.vue')
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('../views/RegisterView.vue')
    },
    {
        path: '/cart',
        name: 'Cart',
        component: () => import('../views/CartView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/orders',
        name: 'Orders',
        component: () => import('../views/OrdersView.vue'),
        meta: { requiresAuth: true }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Navigation Guard
router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    if (to.meta.requiresAuth && !token) {
        next('/login')
    } else {
        next()
    }
})

export default router
