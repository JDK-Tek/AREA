<style>
.main-class {
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

.TextField {
  display: flex;
  /* min-height: 30vh; */
  max-width: 20vw;
  flex-direction: column;
  justify-content: center;
  /* gap: 10px; */
  margin: 0 auto;
  margin-top: 5vh;
}
</style>


<template>
  <main class="main-class">
    <header class="App-header">
      <p>Benchmark AREA - VueJS</p>
      <p>Elise PIPET - Grégoire LANTIM - Paul PARISOT - Esteban MARQUES - John DE KETTELBUTTER</p>
      <p>Temps de chargement : {{ loadTime }} ms</p>
    </header>
    <div class="TextField">
      <TextField label="Email"></TextField>
      <TextField label="Password"></TextField>
    </div>
  </main>
</template>

<script>
import TextField from './TextField.vue';

export default {
  components: {
    TextField
  },
  data() {
    return {
      loadTime: 0
    };
  },
  mounted() {
    const performance = window.performance;
    const loadTime = performance.now();
    this.loadTime = loadTime.toFixed(2); // Arrondir à 2 décimales
  }
}
</script>