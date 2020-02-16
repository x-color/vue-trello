<template>
  <v-app>
    <the-header />
    <v-content class="overflow-x-auto">
      <router-view />
    </v-content>
  </v-app>
</template>

<script>
import TheHeader from '@/components/TheHeader.vue';

export default {
  name: 'App',
  components: {
    TheHeader,
  },
};
</script>
