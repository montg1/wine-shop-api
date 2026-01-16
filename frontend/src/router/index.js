import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

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
    },
    // Admin Routes
    {
        path: '/admin',
        component: () => import('../views/admin/AdminLayout.vue'),
        meta: { requiresAuth: true, requiresAdmin: true },
        children: [
            {
                path: '',
                name: 'AdminDashboard',
                component: () => import('../views/admin/AdminDashboard.vue')
            },
            {
                path: 'products',
                name: 'AdminProducts',
                component: () => import('../views/admin/AdminProducts.vue')
            },
            {
                path: 'products/new',
                name: 'AdminProductNew',
                component: () => import('../views/admin/ProductForm.vue')
            },
            {
                path: 'products/:id/edit',
                name: 'AdminProductEdit',
                component: () => import('../views/admin/ProductForm.vue')
            }
        ]
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Navigation Guard
router.beforeEach(async (to, from, next) => {
    const token = localStorage.getItem('token')

    // Requires authentication
    if (to.meta.requiresAuth && !token) {
        return next('/login')
    }

    // Requires admin role
    if (to.meta.requiresAdmin) {
        const authStore = useAuthStore()
        // Wait for user to be loaded if not already
        if (!authStore.user && token) {
            await authStore.fetchUser()
        }
        if (!authStore.isAdmin) {
            return next('/products')
        }
    }

    next()
})

export default router
