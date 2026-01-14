<template>
  <div class="home">
    <section class="hero">
      <div class="hero-content">
        <h1>🍷 Wine Shop</h1>
        <p>Discover exceptional wines from around the world</p>
        <router-link to="/products" class="btn btn-primary">Browse Wines</router-link>
      </div>
    </section>

    <section class="features">
      <div class="feature">
        <span class="icon">🍇</span>
        <h3>Premium Selection</h3>
        <p>Hand-picked wines from renowned vineyards</p>
      </div>
      <div class="feature">
        <span class="icon">🚚</span>
        <h3>Fast Delivery</h3>
        <p>Get your wines delivered to your door</p>
      </div>
      <div class="feature">
        <span class="icon">⭐</span>
        <h3>Expert Ratings</h3>
        <p>Curated by sommelier experts</p>
      </div>
    </section>
  </div>
</template>

<script setup>
// Home view - no logic needed
</script>

<style scoped>
.home {
  min-height: calc(100vh - 80px);
}

.hero {
  background: linear-gradient(135deg, #2d1f3d 0%, #4a2c5a 100%);
  padding: 100px 20px;
  text-align: center;
  border-radius: 0 0 50px 50px;
}

.hero h1 {
  font-size: 3.5rem;
  margin-bottom: 10px;
  color: #fff;
}

.hero p {
  font-size: 1.3rem;
  color: #d4b8e0;
  margin-bottom: 30px;
}

.features {
  display: flex;
  justify-content: center;
  gap: 40px;
  padding: 60px 20px;
  flex-wrap: wrap;
}

.feature {
  background: var(--card-bg);
  padding: 30px;
  border-radius: 16px;
  text-align: center;
  width: 250px;
  transition: transform 0.3s;
}

.feature:hover {
  transform: translateY(-5px);
}

.feature .icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 15px;
}

.feature h3 {
  margin-bottom: 10px;
  color: var(--primary);
}
</style>
