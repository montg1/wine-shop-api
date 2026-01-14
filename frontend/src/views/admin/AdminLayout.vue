<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <div class="sidebar-header">
        <h2>🍷 Admin</h2>
      </div>
      <nav class="sidebar-nav">
        <router-link to="/admin" class="nav-item" exact>
          <span class="icon">📊</span> Dashboard
        </router-link>
        <router-link to="/admin/products" class="nav-item">
          <span class="icon">🍾</span> Products
        </router-link>
        <router-link to="/" class="nav-item">
          <span class="icon">🏠</span> Back to Shop
        </router-link>
      </nav>
    </aside>
    
    <main class="admin-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
// Admin layout component
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
}

.admin-sidebar {
  width: 250px;
  background: #1a1a2e;
  border-right: 1px solid rgba(139, 69, 102, 0.3);
  padding: 20px 0;
}

.sidebar-header {
  padding: 0 20px 20px;
  border-bottom: 1px solid rgba(139, 69, 102, 0.2);
  margin-bottom: 20px;
}

.sidebar-header h2 {
  color: var(--primary);
  margin: 0;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  color: #aaa;
  text-decoration: none;
  transition: all 0.3s;
}

.nav-item:hover {
  background: rgba(139, 69, 102, 0.15);
  color: #fff;
}

.nav-item.router-link-active {
  background: rgba(139, 69, 102, 0.25);
  color: var(--primary);
  border-left: 3px solid var(--primary);
}

.nav-item .icon {
  font-size: 1.2rem;
}

.admin-content {
  flex: 1;
  padding: 30px;
  background: var(--bg);
  overflow-y: auto;
}
</style>
