<template>
  <div class="home">
    <section class="hero">
      <div class="hero-overlay"></div>
      <div class="hero-content">
        <span class="tagline">Curated Collection</span>
        <h1>Wine Shop</h1>
        <p>Discover exceptional wines from the world's finest vineyards</p>
        <router-link to="/products" class="btn btn-primary">Explore Wines</router-link>
      </div>
    </section>

    <section class="intro">
      <div class="intro-content">
        <h2>The Art of Wine</h2>
        <p>Every bottle tells a story. From sun-drenched vineyards to your table, we bring you wines crafted with passion and tradition.</p>
      </div>
    </section>

    <section class="features">
      <div class="feature">
        <div class="feature-icon">🍇</div>
        <h3>Premium Selection</h3>
        <p>Hand-picked from renowned estates worldwide</p>
      </div>
      <div class="feature">
        <div class="feature-icon">🚚</div>
        <h3>Careful Delivery</h3>
        <p>Temperature-controlled shipping to your door</p>
      </div>
      <div class="feature">
        <div class="feature-icon">⭐</div>
        <h3>Expert Curation</h3>
        <p>Selected by our team of sommeliers</p>
      </div>
    </section>
  </div>
</template>

<script setup>
// Home view
</script>

<style scoped>
.home {
  min-height: calc(100vh - 80px);
}

.hero {
  position: relative;
  height: 70vh;
  min-height: 500px;
  background: url('../assets/images/hero-bg.png') center center / cover no-repeat;
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom, rgba(61, 50, 37, 0.4), rgba(61, 50, 37, 0.6));
}

.hero-content {
  position: relative;
  z-index: 1;
  text-align: center;
  color: #fff;
  padding: 20px;
}

.tagline {
  display: inline-block;
  font-size: 0.85rem;
  letter-spacing: 3px;
  text-transform: uppercase;
  color: var(--accent);
  margin-bottom: 15px;
}

.hero h1 {
  font-size: 4.5rem;
  font-weight: 400;
  letter-spacing: 8px;
  text-transform: uppercase;
  margin-bottom: 20px;
}

.hero p {
  font-size: 1.2rem;
  font-weight: 300;
  max-width: 500px;
  margin: 0 auto 35px;
  opacity: 0.9;
}

.intro {
  padding: 80px 20px;
  background: var(--bg-warm);
  text-align: center;
}

.intro-content {
  max-width: 600px;
  margin: 0 auto;
}

.intro h2 {
  font-size: 2.5rem;
  color: var(--primary);
  margin-bottom: 20px;
}

.intro p {
  font-size: 1.1rem;
  color: var(--text-muted);
  line-height: 1.8;
}

.features {
  display: flex;
  justify-content: center;
  gap: 60px;
  padding: 80px 20px;
  flex-wrap: wrap;
  max-width: 1000px;
  margin: 0 auto;
}

.feature {
  text-align: center;
  width: 250px;
}

.feature-icon {
  font-size: 2.5rem;
  margin-bottom: 20px;
}

.feature h3 {
  font-size: 1.3rem;
  color: var(--text);
  margin-bottom: 10px;
}

.feature p {
  color: var(--text-muted);
  font-size: 0.95rem;
}
</style>
