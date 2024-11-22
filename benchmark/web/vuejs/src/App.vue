<template>
  <div class="App">
    <header class="App-header">
      <p>Benchmark AREA - VueJS</p>
      <p>Elise PIPET - Grégoire LANTIM - Paul PARISOT - Esteban MARQUES - John DE KETTELBUTTER</p>
      <p>Temps de chargement : {{ loadTime }} ms</p>
    </header>
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  data() {
    return {
      loadTime: 0
    };
  },
  mounted() {
    const performance = window.performance;
    const loadTime = performance.now();
    this.loadTime = loadTime.toFixed(2);
  }
}
</script>

<style>
.App {
  text-align: center;
  background-color: #f1f1f1;
  min-height: 100vh;
}

.App-header {
  display: flex;
  flex-direction: column;
  color: rgb(92, 92, 92);
  font-size: 1.3rem;
  min-width: 100vw;
  gap: 0.4rem;
}

</style>