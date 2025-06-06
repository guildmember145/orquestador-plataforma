import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '@/stores';
import { LoginView, RegisterView, WorkflowsDashboard } from '@/views';

// Añadimos el tipo explícito al array de rutas para ayudar a TypeScript
const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'Home',
        redirect: () => {
            const authStore = useAuthStore();
            return authStore.isAuthenticated ? '/dashboard/workflows' : '/login';
        },
    },
    {
        path: '/login',
        name: 'Login',
        component: LoginView,
    },
    {
        path: '/register',
        name: 'Register',
        component: RegisterView,
    },
    {
        path: '/dashboard',
        redirect: '/dashboard/workflows',
    },
    {
        path: '/dashboard/workflows',
        name: 'WorkflowsDashboard',
        component: WorkflowsDashboard,
        meta: { requiresAuth: true },
    },
    {
        path: '/dashboard/workflows/new',
        name: 'WorkflowCreate',
        component: () => import('@/views/WorkflowCreateView.vue'),
        meta: { requiresAuth: true },
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore();
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

    if (requiresAuth && !authStore.isAuthenticated) {
        next('/login');
    } else if ((to.path === '/login' || to.path === '/register') && authStore.isAuthenticated) {
        next('/dashboard/workflows');
    } else {
        next();
    }
});

export default router;