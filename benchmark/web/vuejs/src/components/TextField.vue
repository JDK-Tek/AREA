<!-- components/TextField.vue -->
<template>
  <v-text-field :label="label" variant="outlined"></v-text-field>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true
    }
  }
}
</script>