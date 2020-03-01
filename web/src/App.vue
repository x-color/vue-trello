<template>
  <v-app>
    <the-header />
    <v-content class="overflow-x-auto">
      <router-view />
    </v-content>
  </v-app>
</template>

<script>
import { mapActions } from 'vuex';
import TheHeader from '@/components/TheHeader.vue';

export default {
  name: 'App',
  components: {
    TheHeader,
  },
  created() {
    if (this.$store.state.user.user.login) {
      this.loadResources();
    }
  },
  methods: {
    ...mapActions(['loadResources']),
  },
};
</script>
