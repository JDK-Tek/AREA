<template>
    <div class="TextField">
        <TextField label="Email"></TextField>
        <TextField label="Password"></TextField>
        <v-btn variant="elevated" color="blue" @click="goToHome">
        Confirm
        </v-btn>
    </div>
</template>

<style>

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

<script>
import TextField from './TextField.vue';


export default {
    components: {
        TextField
    },
    name: 'Login',
    methods: {
        goToHome() {
            this.$router.push('/home')
        },
    }
};

</script>