<template>
  <v-chip :color="tag.color" class="ma-1" small label text-color="white">
    <v-icon small left>mdi-label</v-icon>
    {{ tag.title }}
  </v-chip>
</template>
<script>
export default {
  name: 'BaseTag',
  props: {
    tag: Object,
  },
};
</script>
